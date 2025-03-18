package plugin

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Plugin represents a Goblin plugin that can be loaded into the application
type Plugin interface {
	// Name returns the unique name of the plugin
	Name() string

	// Version returns the version of the plugin
	Version() string

	// Description returns a short description of what the plugin does
	Description() string

	// Dependencies returns a list of plugin names that this plugin depends on
	Dependencies() []string

	// OnRegister is called when the plugin is registered with the application
	OnRegister(ctx context.Context) error

	// OnStart is called when the application starts
	OnStart(ctx context.Context) error

	// OnStop is called when the application stops
	OnStop(ctx context.Context) error

	// RegisterRoutes registers the plugin's routes with the Gin router
	RegisterRoutes(router *gin.Engine) error

	// RegisterDependencies registers the plugin's dependencies with the fx container
	RegisterDependencies(app *fx.App) error
}

// PluginManager manages the lifecycle of plugins
type PluginManager struct {
	plugins map[string]Plugin
	order   []string // Plugin load order based on dependencies
}

// NewPluginManager creates a new plugin manager
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
		order:   make([]string, 0),
	}
}

// RegisterPlugin registers a new plugin
func (pm *PluginManager) RegisterPlugin(plugin Plugin) error {
	name := plugin.Name()
	if _, exists := pm.plugins[name]; exists {
		return fmt.Errorf("plugin %s is already registered", name)
	}

	// Check dependencies
	for _, dep := range plugin.Dependencies() {
		if _, exists := pm.plugins[dep]; !exists {
			return fmt.Errorf("plugin %s depends on %s which is not registered", name, dep)
		}
	}

	pm.plugins[name] = plugin
	return nil
}

// GetPlugin returns a registered plugin by name
func (pm *PluginManager) GetPlugin(name string) (Plugin, bool) {
	plugin, exists := pm.plugins[name]
	return plugin, exists
}

// GetAllPlugins returns all registered plugins
func (pm *PluginManager) GetAllPlugins() []Plugin {
	plugins := make([]Plugin, 0, len(pm.plugins))
	for _, plugin := range pm.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// SortPlugins sorts plugins based on their dependencies
func (pm *PluginManager) SortPlugins() error {
	visited := make(map[string]bool)
	temp := make(map[string]bool)
	order := make([]string, 0)

	var visit func(name string) error
	visit = func(name string) error {
		if temp[name] {
			return fmt.Errorf("circular dependency detected: %s", name)
		}
		if visited[name] {
			return nil
		}

		temp[name] = true
		plugin := pm.plugins[name]

		for _, dep := range plugin.Dependencies() {
			if err := visit(dep); err != nil {
				return err
			}
		}

		temp[name] = false
		visited[name] = true
		order = append(order, name)
		return nil
	}

	for name := range pm.plugins {
		if !visited[name] {
			if err := visit(name); err != nil {
				return err
			}
		}
	}

	pm.order = order
	return nil
}

// GetLoadOrder returns the sorted plugin load order
func (pm *PluginManager) GetLoadOrder() []string {
	return pm.order
}

// PluginConfig represents the configuration for a plugin
type PluginConfig struct {
	Enabled bool                   `json:"enabled"`
	Options map[string]interface{} `json:"options"`
}

// PluginContext provides context information to plugins
type PluginContext struct {
	Config     *PluginConfig
	App        *fx.App
	Router     *gin.Engine
	Logger     *log.Logger
	Context    context.Context
	CancelFunc context.CancelFunc
}

// NewPluginContext creates a new plugin context
func NewPluginContext(
	config *PluginConfig,
	app *fx.App,
	router *gin.Engine,
	logger *log.Logger,
) *PluginContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &PluginContext{
		Config:     config,
		App:        app,
		Router:     router,
		Logger:     logger,
		Context:    ctx,
		CancelFunc: cancel,
	}
}

// Stop stops the plugin context
func (pc *PluginContext) Stop() {
	pc.CancelFunc()
}
