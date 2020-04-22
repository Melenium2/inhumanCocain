package support

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	auth "github.com/inhumanLightBackend/auth/logic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func authenticate() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			error, ok := ctx.Value(auth.CtxErrorKey).(bool)
			if ok && error {
				return nil, status.Error(codes.PermissionDenied, "Error not authenticated")
			}
			return next(ctx, request)
		}
	}
}

func translateMetadataToContext() grpc.ServerRequestFunc {
	return func(ctx context.Context, m metadata.MD) context.Context {
		authHeader, ok := m["authorization"]
		if !ok {
			return context.WithValue(ctx, auth.CtxErrorKey, true)
		}
		token := authHeader[0]
		return context.WithValue(ctx, auth.CtxUserKey, token)
	}
} 

func translateJwtToMetadata() grpc.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		t := ctx.Value(auth.CtxUserKey).(string)
		key, val := encodeKeyValue("authorization", t)
		(*md)[key] = append((*md)[key], val)
		return ctx
	}
}

func encodeKeyValue(key, val string) (string, string) {
	key = strings.ToLower(key)
	if strings.HasSuffix(key, "-bin") {
		val = base64.StdEncoding.EncodeToString([]byte(val))
	}
	return key, val
}