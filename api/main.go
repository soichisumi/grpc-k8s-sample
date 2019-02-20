package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/soichisumi/grpc-auth-sample/api-pb"
	"github.com/soichisumi/grpc-auth-sample/auth"
	"github.com/soichisumi/grpc-auth-sample/meta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
)

const (
	port = ":3000"
	rsaPrivateKeyPath = "./privKey.pem"
	rsaPublicKeyPath = "./privKey.pem.pub.pkcs8"
)

var (
	PrivKey *rsa.PrivateKey       // to generate token
	PubKey  *rsa.PublicKey        // to validate token
)

func validateToken(token string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			log.Printf("Unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("invalid token")
		}
		return PubKey, nil
	})
	if err == nil && jwtToken.Valid {
		return jwtToken, nil
	}
	return nil, err
}

func authenticationFunc() auth.AuthenticationFunc {
	return func(ctx context.Context) (context.Context, error) {
		authorization, err := meta.Authorization(ctx)
		if err != nil {
			return nil, err
		}

		_, err = validateToken(authorization)
		if err != nil {
			return nil, err
		}

		return ctx, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %+v\n", err)
	}

	privKey, err := ioutil.ReadFile(rsaPrivateKeyPath)
	if err != nil {
		log.Fatalf("Error reading the jwt private key: %s", err)
	}
	parsedPrivKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKey)
	if err != nil {
		log.Fatalf("Error parsing the jwt private key: %s", err)
	}
	PrivKey = parsedPrivKey

	pubKey, err := ioutil.ReadFile(rsaPublicKeyPath)
	if err != nil {
		log.Fatalf("Error reading the jwt public key: %s", err)
	}
	parsedPubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		log.Fatalf("Error parsing the jwt public key: %s", err)
	}
	PubKey = parsedPubKey

	server, err := NewServer(parsedPrivKey)
	if err != nil{
		log.Fatalf("failed to create server: %+v\n", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthenticationInterceptor(authenticationFunc())),
	)

	apipb.RegisterUserServiceServer(s, server)
	reflection.Register(s)
	fmt.Printf("grpc server is running on port:%s...\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}

}
