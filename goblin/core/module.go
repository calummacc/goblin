// goblin/core/module.go
package core

import (
	"go.uber.org/fx"
)

// GoblinModule represents a module in the Goblin framework
type GoblinModule struct {
	Name    string
	Options fx.Option
}

// ModuleOptions configures a module
type ModuleOptions struct {
	Imports     []GoblinModule
	Controllers []interface{}
	Providers   []interface{}
	Exports     []interface{}
}

// NewModule creates a new Goblin module (similar to @Module() in NestJS)
func NewModule(name string, options ModuleOptions) GoblinModule {
	var fxOptions []fx.Option

	// Add providers
	if len(options.Providers) > 0 {
		fxOptions = append(fxOptions, fx.Provide(options.Providers...))
	}

	// Add controllers as providers
	if len(options.Controllers) > 0 {
		fxOptions = append(fxOptions, fx.Provide(options.Controllers...))
	}

	// Import other modules
	for _, module := range options.Imports {
		fxOptions = append(fxOptions, module.Options)
	}

	// Create the module
	return GoblinModule{
		Name:    name,
		Options: fx.Options(fxOptions...),
	}
}

// DependencyInjection manages the dependency injection container
type DependencyInjection struct {
	app *fx.App
}

// NewDependencyInjection creates a new DI container
func NewDependencyInjection(options ...fx.Option) *DependencyInjection {
	return &DependencyInjection{
		app: fx.New(options...),
	}
}
