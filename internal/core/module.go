package goblin

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module interface (defines the structure of a module)
type Module interface {
	Name() string               // Module name
	Provide() fx.Option         // Fx options for this module
	RegisterRoutes(*gin.Engine) // Registers the module's routes
}

// ModuleRegistry to manage module registration (dynamic loading)
type ModuleRegistry struct {
	modules map[string]Module
}

// NewModuleRegistry creates a new ModuleRegistry
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make(map[string]Module),
	}
}

// Register registers a module with the registry
func (r *ModuleRegistry) Register(module Module) {
	r.modules[module.Name()] = module
}

// LoadModules loads modules dynamically.  This requires a mechanism to read module definitions (e.g., from a configuration file).  This is a placeholder.
func (r *ModuleRegistry) LoadModules() error {
	//In a real-world scenario, you would load modules dynamically here.  For example, read module definitions from config file
	//This is a placeholder function that simulates loading modules.
	return nil
}

// RegisterRoutesForModules registers routes for all registered modules
func (r *ModuleRegistry) RegisterRoutesForModules(engine *gin.Engine) {
	for moduleName, module := range r.modules {
		log.Printf("Registering routes for module: %s", moduleName)
		module.RegisterRoutes(engine)
	}
}
