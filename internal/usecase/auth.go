package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/arunima10a/task-manager/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthInteractor struct {
	userRepo  UserRepoInterface
	taskRepo      TaskRepo
	tx        TransactionManager
	jwtSecret string
}

func NewAuth(u UserRepoInterface, t TaskRepo, tx TransactionManager, secret string) *AuthInteractor {
	return &AuthInteractor{
		userRepo:  u,
		taskRepo:  t,
		tx:        tx,
		jwtSecret: secret,
	}
}

func (a *AuthInteractor) SignUp(ctx context.Context, email, password string) error {
	return a.tx.RunInTx(ctx, func(txCtx context.Context) error {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := entity.User{Email: email, Password: string(hashed)}

		userID, err := a.userRepo.Create(txCtx, user)
		if err != nil {
			return err
		}
		welcomeTask := entity.Task{
			Title: "Welcome",
			Description: "Account created succesfully.",
			Status: "active",
		}
		

		return a.taskRepo.Store(txCtx, welcomeTask, userID)

	})

}

func (a *AuthInteractor) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(a.jwtSecret))
}
