package api

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/soichisumi/protobuf-trial/pbtrial"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
)

type UserServiceServer struct {
	Users map[string]pbtrial.User //todo: goroutine unsafe
	PrivKey *rsa.PrivateKey // to generate token
	PubKey *rsa.PublicKey // to validate token
}

func validateToken(token string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			log.Printf("Unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("invalid token")
		}
		return publicKey, nil
	})
	if err == nil && jwtToken.Valid {
		return jwtToken, nil
	}
	return nil, err
}


func NewServer(rsaPrivateKeyPath string) (*UserServiceServer, error) {
	key, err := ioutil.ReadFile(rsaPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading the jwt private key: %s", err)
	}
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		return nil, fmt.Errorf("Error parsing the jwt private key: %s", err)
	}
	userDB := make(map[string]pbtrial.User)
	return &UserServiceServer{
		Users: userDB,
		PrivKey: parsedKey,
	}, nil
}

func (s *UserServiceServer) AddUser(ctx context.Context, req *pbtrial.AddUserRequest) (*pbtrial.AddUserResponse, error) {
	user := req.User
	if user.Name == "" {
		fmt.Printf("username is empty. user: %+v\n", user)
		return &pbtrial.AddUserResponse{}, status.Error(codes.InvalidArgument, "")
	}
	s.Users[user.Name] = *user

	return &pbtrial.AddUserResponse{}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pbtrial.GetUserRequest) (*pbtrial.GetUserResponse, error) {
	username := req.Username
	if username == "" {
		fmt.Printf("username is empty. username: %s\n", username)
		return &pbtrial.GetUserResponse{}, status.Error(codes.InvalidArgument, "")
	}
	user, ok := s.Users[username]
	if !ok {
		fmt.Println("given user is not found")
		return &pbtrial.GetUserResponse{}, status.Error(codes.NotFound, "")
	}
	return &pbtrial.GetUserResponse{
		User: &user,
	}, nil
}

