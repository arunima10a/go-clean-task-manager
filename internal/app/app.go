package app

import (
	"errors"

	"net"

	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/arunima10a/task-manager/config"
	amqpV1 "github.com/arunima10a/task-manager/internal/controller/amqp/v1"
	grpcV1 "github.com/arunima10a/task-manager/internal/controller/grpc/v1"
	v1 "github.com/arunima10a/task-manager/internal/controller/http/v1"
	"github.com/arunima10a/task-manager/internal/usecase"
	"github.com/arunima10a/task-manager/internal/usecase/repo"
	"github.com/arunima10a/task-manager/pkg/httpserver"
	"github.com/arunima10a/task-manager/pkg/logger"
	"github.com/arunima10a/task-manager/pkg/postgres"
	"github.com/arunima10a/task-manager/pkg/rabbitmq"
	"github.com/arunima10a/task-manager/pkg/redis"
	"github.com/arunima10a/task-manager/pkg/webapi"
	pb "github.com/arunima10a/task-manager/proto/v1"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {

	l := logger.New(cfg.Log.Level)
	motivationURL := "https://zenquotes.io/api/random"
	motivationAPI := webapi.NewMotivationAPI(motivationURL)

	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Error(err, "app - Run - postgres.New")
		return
	}
	defer pg.Close()

	m, err := migrate.New(
		"file://migrations",
		cfg.PG.URL,
	)
	if err != nil {
		l.Error(err, "app - Run - migrate.New")
		return
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		l.Error(err, "app - Run - m.Up")
		return
	}
	l.Info("Migrations applied successfully")

	rmq, err := rabbitmq.New("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		l.Error(err, "app - Run - rabbitmq.New")
	}
	defer func() { _ = rmq.Close() }()

	taskMQ := repo.NewTaskRMQ(rmq)

	l.Info("Background Worker started")

	re, err := redis.New(cfg.Redis.URL)
	if err != nil {
		l.Error(err, "app - Run - redis.New")
	}
	defer func() { _ = re.Close() }()

	taskCache := repo.NewTaskCache(re)

	var wg sync.WaitGroup

	taskRepo := repo.New(pg, l)
	TaskUseCase := usecase.New(taskRepo, motivationAPI, l, &wg, taskMQ, taskCache)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcV1.AuthInterceptor("your-secret-key")),
	)

	catRepo := repo.NewcategoryRepo(pg)
	catUseCase := usecase.NewCategoryInteractor(catRepo)

	taskConsumer := amqpV1.NewTaskConsumer(TaskUseCase, l)
	taskConsumer.Start(rmq.Conn)

	TaskHandler := grpcV1.NewTaskHandler(TaskUseCase, catUseCase)

	pb.RegisterTaskServiceServer(grpcServer, TaskHandler)

	l_grpc, err := net.Listen("tcp", ":50051")
	if err != nil {
		l.Error(err, "app - Run - net.Listen (gRPC)")
	}

	go func() {
		l.Info("gRPC server started on port 50051")
		if err := grpcServer.Serve(l_grpc); err != nil {
			l.Error(err, "app - Run - grpcServer.Serve")

		}
	}()

	handler := gin.New()
	userRepo := repo.NewUserRepo(pg)
	authUseCase := usecase.NewAuth(userRepo, taskRepo, pg, "your-secret-key")
	v1.NewRouter(handler, TaskUseCase, authUseCase, catUseCase, pg, re)

	httpServer := httpserver.New(
		handler,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(time.Second*10),
	)

	l.Info("Starting %s app, version: %s on port %s", cfg.App.Name, Version, cfg.HTTP.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		l.Info("%s", "app - Run - signal: "+s.String())
	case err := <-httpServer.Notify():
		l.Error(err, "app - Run - httpSever.Notify")
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(err, "app - Run - httpServer.Shutdown")
	}

	l.Info("Shutting down gRPC server...")

	grpcServer.GracefulStop()

	l.Info("Waiting for background workers to finish...")
	wg.Wait()

	l.Info("Closing database connections...")
	pg.Close()

	l.Info("App excited cleanly")
}
