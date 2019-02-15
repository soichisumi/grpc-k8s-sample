.PHONY: pb
pb:
	protoc -I=. -I=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. ./api-pb/types.proto

.PHONY: pbgw
pbgw:
	protoeasy --go --grpc --grpc-gateway --go-import-path github.com/soichisumi/grpc-auth-sample --out api-pb api-pb


.PHONY: gen-priv
gen-rsa:
	ssh-keygen -t rsa -b 4096 -f privKey.pem

.PHONY: gen-pub
gen-pub:
	ssh-keygen -f privkey.pem.pub -e -m pkcs8 > privkey.pem.pub.pkcs8