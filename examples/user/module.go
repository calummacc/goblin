package user

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("UserModule",
		fx.Provide(
			NewUserController,
			NewUserService,
		),
	)
}
