package user

import "go.uber.org/fx"

// User module providers
var Module = fx.Options(
	fx.Provide(
		NewUserService,
		NewUserController,
	),
)
