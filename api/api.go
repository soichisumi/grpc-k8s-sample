package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/soichisumi/grpc-auth-sample/api-pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type UserServiceServer struct {
	Users   map[string]apipb.User // Todo: goroutine unsafe
	PrivKey *rsa.PrivateKey
}

func NewServer(rsaPrivateKey *rsa.PrivateKey) (*UserServiceServer, error) {
	userDB := make(map[string]apipb.User)
	return &UserServiceServer{
		Users:   userDB,
		PrivKey: rsaPrivateKey,
	}, nil
}

func (s *UserServiceServer) AddUser(ctx context.Context, req *apipb.AddUserRequest) (*apipb.AddUserResponse, error) {
	user := req.User
	fmt.Println("given user: ")
	spew.Println(user)
	if user.Name == "" {
		fmt.Printf("username is empty. user: %+v\n", user)
		return &apipb.AddUserResponse{}, status.Error(codes.InvalidArgument, "")
	}
	fmt.Println("s.Users: ")
	spew.Println(s.Users)
	s.Users[user.Name] = *user

	return &apipb.AddUserResponse{}, nil
}

// authorization required
func (s *UserServiceServer) GetUser(ctx context.Context, req *apipb.GetUserRequest) (*apipb.GetUserResponse, error) {

	username := req.Username
	if username == "" {
		fmt.Printf("username is empty. username: %s\n", username)
		return &apipb.GetUserResponse{}, status.Error(codes.InvalidArgument, "")
	}
	user, ok := s.Users[username]
	if !ok {
		fmt.Println("given user is not found")
		return &apipb.GetUserResponse{}, status.Error(codes.NotFound, "")
	}
	return &apipb.GetUserResponse{
		User: &user,
	}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *apipb.LoginRequest) (*apipb.LoginResponse, error) {
	if req.User.Name == "" || req.User.Password == "" {
		return &apipb.LoginResponse{}, status.Error(codes.InvalidArgument, "invalid request.")
	}
	dbUser, ok := s.Users[req.User.Name]
	if !ok {
		return &apipb.LoginResponse{}, status.Error(codes.InvalidArgument, "user does not exists.")
	}
	if req.User.Name != dbUser.Name || req.User.Password != dbUser.Password {
		return &apipb.LoginResponse{}, status.Error(codes.InvalidArgument, "invalid userid or password.")
	}

	// create token
	token := jwt.New(jwt.SigningMethodRS256)

	// set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = req.User.Name
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	spew.Printf("dump. claim: %+v\n", claims)

	tokenString, err := token.SignedString(s.PrivKey)
	if err != nil {
		fmt.Println(err)
	}
	return &apipb.LoginResponse{Token: tokenString}, nil
}
