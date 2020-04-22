package cocain

import (
	"github.com/go-kit/kit/log"
	auth "github.com/inhumanLightBackend/auth/logic"
	notifications "github.com/inhumanLightBackend/notifications/logic"
	support "github.com/inhumanLightBackend/support/logic"
	user "github.com/inhumanLightBackend/user/logic"
	"google.golang.org/grpc"
)

var (
	AuthPort = ":1011"
	SupportPort = ":5051"
	NotificationsPort = ":6061"
	UserPort = ":7071"
)

func AuthGRPCService(instance string, logger log.Logger) (auth.Service, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(instance, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return auth.NewGRPCClient(conn, logger), conn, nil
}

func UserGRPCService(instance string, logger log.Logger) (user.Service, *grpc.ClientConn, error){
	conn, err := grpc.Dial(instance, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return user.NewGRPCClient(conn, logger), conn, nil
}

func NotificationsGRPCService(instance string, logger log.Logger) (notifications.Service, *grpc.ClientConn, error){
	conn, err := grpc.Dial(instance, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return notifications.NewGRPCClient(conn, logger), conn, nil
}

func SupportGRPCService(instance string, logger log.Logger) (support.Service, *grpc.ClientConn, error){
	conn, err := grpc.Dial(instance, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return support.NewGRPCClient(conn, logger), conn, nil
}