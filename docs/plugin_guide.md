# Plugin Development Guide

This guide explains how to create plugins for the Goblin Framework.

## Overview

Plugins in Goblin Framework are modular components that can be loaded into your application to extend its functionality. They follow a similar pattern to modules but are designed to be more lightweight and focused on specific features.

## Plugin Interface

Every plugin must implement the `Plugin` interface:

```go
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
    RegisterDependencies(container *fx.Container) error
}
```

## Creating a Plugin

Here's a step-by-step guide to creating a plugin:

1. Create a new package for your plugin:

```go
package myplugin

import (
    "context"
    "goblin/plugin"
    "log"
    "github.com/gin-gonic/gin"
    "go.uber.org/fx"
)
```

2. Create your plugin struct:

```go
type MyPlugin struct {
    config *plugin.PluginConfig
    logger *log.Logger
}
```

3. Create a constructor:

```go
func NewMyPlugin(config *plugin.PluginConfig, logger *log.Logger) *MyPlugin {
    return &MyPlugin{
        config: config,
        logger: logger,
    }
}
```

4. Implement the Plugin interface:

```go
func (p *MyPlugin) Name() string {
    return "myplugin"
}

func (p *MyPlugin) Version() string {
    return "1.0.0"
}

func (p *MyPlugin) Description() string {
    return "My custom plugin for Goblin Framework"
}

func (p *MyPlugin) Dependencies() []string {
    return []string{} // List any dependencies here
}

func (p *MyPlugin) OnRegister(ctx context.Context) error {
    p.logger.Printf("Registering my plugin...")
    return nil
}

func (p *MyPlugin) OnStart(ctx context.Context) error {
    p.logger.Printf("Starting my plugin...")
    return nil
}

func (p *MyPlugin) OnStop(ctx context.Context) error {
    p.logger.Printf("Stopping my plugin...")
    return nil
}

func (p *MyPlugin) RegisterRoutes(router *gin.Engine) error {
    my := router.Group("/my")
    {
        my.GET("/", p.handleGet)
        my.POST("/", p.handlePost)
    }
    return nil
}

func (p *MyPlugin) RegisterDependencies(container *fx.Container) error {
    // Register your services and dependencies here
    container.Provide(NewMyService)
    return nil
}
```

5. Implement your handlers and services:

```go
type MyService struct {
    // Your service fields
}

func NewMyService() *MyService {
    return &MyService{}
}

func (p *MyPlugin) handleGet(c *gin.Context) {
    c.JSON(200, gin.H{"message": "Hello from my plugin!"})
}

func (p *MyPlugin) handlePost(c *gin.Context) {
    // Handle POST request
}
```

## Plugin Configuration

Plugins can be configured using the `PluginConfig` struct:

```go
type PluginConfig struct {
    Enabled bool                   `json:"enabled"`
    Options map[string]interface{} `json:"options"`
}
```

Example configuration:

```go
config := &plugin.PluginConfig{
    Enabled: true,
    Options: map[string]interface{}{
        "option1": "value1",
        "option2": 42,
    },
}
```

## Using Plugins in Your Application

1. Create a plugin manager:

```go
pluginManager := plugin.NewPluginManager()
```

2. Create and register your plugins:

```go
myPlugin := myplugin.NewMyPlugin(config, logger)
if err := pluginManager.RegisterPlugin(myPlugin); err != nil {
    log.Fatal(err)
}
```

3. Sort plugins based on dependencies:

```go
if err := pluginManager.SortPlugins(); err != nil {
    log.Fatal(err)
}
```

4. Create a plugin context:

```go
pluginCtx := plugin.NewPluginContext(
    config,
    container,
    router,
    logger,
)
```

5. Register dependencies and routes:

```go
for _, name := range pluginManager.GetLoadOrder() {
    p, _ := pluginManager.GetPlugin(name)
    p.RegisterDependencies(container)
    p.RegisterRoutes(router)
}
```

6. Start plugins:

```go
for _, name := range pluginManager.GetLoadOrder() {
    p, _ := pluginManager.GetPlugin(name)
    p.OnStart(pluginCtx.Context)
}
```

## Best Practices

1. **Dependency Management**:
   - Clearly declare plugin dependencies
   - Handle dependency injection properly
   - Use the fx container for dependency management

2. **Error Handling**:
   - Return meaningful errors
   - Log errors appropriately
   - Clean up resources in OnStop

3. **Configuration**:
   - Use the PluginConfig for configuration
   - Validate configuration options
   - Provide sensible defaults

4. **Resource Management**:
   - Clean up resources in OnStop
   - Use context for cancellation
   - Handle goroutines properly

5. **Testing**:
   - Write unit tests for your plugin
   - Test plugin lifecycle methods
   - Mock dependencies for testing

## Example Plugins

Check out the example plugins in the `examples/plugins` directory:

- `auth`: Authentication plugin with JWT support
- `cache`: In-memory caching plugin

These examples demonstrate different aspects of plugin development and can serve as templates for your own plugins.

## Publishing Plugins

To make your plugin available to others:

1. Create a GitHub repository for your plugin
2. Use semantic versioning for releases
3. Document your plugin's features and configuration
4. Add examples and tests
5. Submit a pull request to add your plugin to the Goblin plugins list

## Support

For questions and support:

1. Check the documentation
2. Look at example plugins
3. Open an issue on GitHub
4. Join the Goblin community 