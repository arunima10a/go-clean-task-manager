package v1

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/metadata"


)

type ctxKey string

func AuthInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{} , error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is missing")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is required")
		}

		tokenString := strings.Replace(values[0], "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
			return []byte(secret), nil
		})

		if err != nil  || token.Valid{
			return nil, status.Error(codes.Unauthenticated, "invalid token")


		}

		claims, _ := token.Claims.(jwt.MapClaims)
		userID := int(claims["user_id"].(float64))

		ctx = context.WithValue(ctx, ctxKey("user_id"), userID)

		return handler(ctx, req)
	}
}
