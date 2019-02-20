package meta

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	// AuthorizationKey は認証トークンに対応するキーを表す
	AuthorizationKey = "authorization"
)

// Authorization は gRPC メタデータからユーザー認証トークンを取得する
func Authorization(ctx context.Context) (string, error) {
	return fromMeta(ctx, AuthorizationKey)
}

func fromMeta(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("not found metadata")
	}
	vs := md[key]
	if len(vs) == 0 {
		return "", fmt.Errorf("not found %s in metadata", key)
	}
	return vs[0], nil
}
