package user

import (
	"github.com/go-kit/kit/log"
	auth "github.com/inhumanLightBackend/auth/logic"
	"google.golang.org/grpc"
)

var (
	AuthserverPort = ":1011"
)

func AuthGRPCService(instance string, logger log.Logger) (auth.Service, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(instance, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	
	return auth.NewGRPCClient(conn, logger), conn, nil
}