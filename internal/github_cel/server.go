package github_cel

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net"
	"rigsecurity/internal/github_cel/config"
	"rigsecurity/internal/github_cel/interfaces"
)

func NewListener(config *config.Config) (net.Listener, error) {
	return net.Listen("tcp", ":50051")

}
func NewServer(lc fx.Lifecycle, lis net.Listener, s interfaces.GitAccessInformationV1Server) *grpc.Server {
	server := grpc.NewServer()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				interfaces.RegisterGitAccessInformationV1Server(server, s)

				if err := server.Serve(lis); err != nil {
					fmt.Errorf("")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Stop()
			return nil
		},
	})
	return server
}
