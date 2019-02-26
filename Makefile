pb:
	protoc -I=. -I=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. ./api-pb/types.proto

pbgw:
	protoeasy --go --grpc --grpc-gateway --go-import-path github.com/soichisumi/grpc-auth-sample --out api-pb api-pb

gen-rsa:
	ssh-keygen -t rsa -b 4096 -f privKey.pem

gen-pub:
	ssh-keygen -f privkey.pem.pub -e -m pkcs8 > privkey.pem.pub.pkcs8

# tag for azure container registry: hogehoge.azurecr.io/grpc-xx:latest
build-containers:
	docker build --file Dockerfile.api --tag testxk8s.azurecr.io/grpc-api .
	docker build --file Dockerfile.gw --tag testxk8s.azurecr.io/grpc-gw .
	docker push testxk8s.azurecr.io/grpc-api
	docker push testxk8s.azurecr.io/grpc-gw
