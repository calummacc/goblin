package core

import (
	"go.uber.org/fx"

	// Import your modules or packages
	"github.com/calummacc/goblin/internal/modules/user"
)

// AppModule collects all modules in a single fx.Options
var AppModule = fx.Options(
	RouterModule, // router.go
	user.Module,  // internal/modules/user
)
