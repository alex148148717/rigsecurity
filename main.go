package main

import (
	"context"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"rigsecurity/internal/github_cel"
	"rigsecurity/internal/github_cel/config"
	"rigsecurity/internal/github_cel/domain"
	"rigsecurity/internal/github_cel/infrastructure"
	"rigsecurity/internal/github_cel/interfaces"
)

func main() {

	fx.New(
		fx.Provide(
			config.LoadConfig,
			domain.NewAccessInformationService,
			interfaces.NewGitAccessInformationServer,
			github_cel.NewListener,
			github_cel.NewServer,
			infrastructure.NewCelBuilder,
			infrastructure.NetGitBuilder,
		),
		fx.Invoke(func(lc fx.Lifecycle, config *config.Config) {

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {

					return nil
				},
			})

		}),
		fx.Invoke(
			func(*grpc.Server) {},
		),
	).Run()
}
