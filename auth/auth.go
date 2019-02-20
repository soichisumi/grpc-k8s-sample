package auth

import (
	"context"
	"google.golang.org/grpc"
)

// AuthenticationFunc はすべてのサービスに共通する認証処理を行う関数を表す
type AuthenticationFunc func(ctx context.Context) (context.Context, error)

func isAuthorizedMethod(methodName string) bool {
	switch methodName {
	case "/apipb.UserService/Login", "/apipb.UserService/AddUser":
		return false
	}
	return true
}

// AuthenticationInterceptor はリクエストごとの認証処理を行う
func AuthenticationInterceptor(authFunc AuthenticationFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var err error
		if isAuthorizedMethod(info.FullMethod) {
			ctx, err = authFunc(ctx)
			if err != nil {
				return nil, err
			}
		}
		return handler(ctx, req)
	}
}
