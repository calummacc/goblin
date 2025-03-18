# Goblin Framework - Module System Example

This example demonstrates how to use the Module System in the Goblin Framework, which is inspired by NestJS's module system.

## Overview

The Module System in Goblin Framework provides a way to organize your application into modular components, similar to NestJS. Each module can:

- Import other modules
- Provide services and controllers
- Export providers for use in other modules
- Handle lifecycle events (initialization and cleanup)

## Key Components

### 1. Module Interface

```go
type Module interface {
    GetMetadata() ModuleMetadata
    OnModuleInit(ctx context.Context) error
    OnModuleDestroy(ctx context.Context) error
}
```

### 2. Module Metadata

```go
type ModuleMetadata struct {
    Imports     []Module
    Exports     []interface{}
    Providers   []interface{}
    Controllers []interface{}
}
```

### 3. Base Module

The `BaseModule` struct provides a default implementation of the `Module` interface:

```go
type BaseModule struct {
    metadata ModuleMetadata
}
```

## Usage Example

### 1. Creating a Module

```go
type UserModule struct {
    *core.BaseModule
}

func NewUserModule() *UserModule {
    module := &UserModule{}
    module.BaseModule = core.NewBaseModule(core.ModuleMetadata{
        Providers: []interface{}{
            NewDatabase,
            NewUserService,
        },
        Controllers: []interface{}{
            NewUserController,
        },
    })
    return module
}
```

### 2. Implementing Lifecycle Hooks

```go
func (m *UserModule) OnModuleInit(ctx context.Context) error {
    // Initialize module resources
    return nil
}

func (m *UserModule) OnModuleDestroy(ctx context.Context) error {
    // Clean up module resources
    return nil
}
```

### 3. Importing Other Modules

```go
type AppModule struct {
    *core.BaseModule
}

func NewAppModule() *AppModule {
    module := &AppModule{}
    module.BaseModule = core.NewBaseModule(core.ModuleMetadata{
        Imports: []core.Module{
            NewUserModule(),
            NewAuthModule(),
        },
    })
    return module
}
```

### 4. Using the Module Manager

```go
// Create the module manager
moduleManager := core.NewModuleManager()

// Register modules
appModule := NewAppModule()
if err := moduleManager.RegisterModule(appModule); err != nil {
    log.Fatalf("Failed to register app module: %v", err)
}

// Initialize modules
ctx := context.Background()
if err := moduleManager.InitializeModules(ctx); err != nil {
    log.Fatalf("Failed to initialize modules: %v", err)
}

// Get module dependencies
deps := moduleManager.GetModuleDependencies(appModule)

// Get module providers
providers := moduleManager.GetModuleProviders(appModule)

// Get module controllers
controllers := moduleManager.GetModuleControllers(appModule)
```

## Best Practices

1. **Module Organization**:
   - Keep related functionality together in a module
   - Use clear, descriptive module names
   - Keep modules focused and single-purpose

2. **Dependency Management**:
   - Explicitly declare module dependencies using `Imports`
   - Export only what's necessary using `Exports`
   - Avoid circular dependencies between modules

3. **Lifecycle Management**:
   - Implement `OnModuleInit` for resource initialization
   - Implement `OnModuleDestroy` for cleanup
   - Handle errors appropriately in lifecycle hooks

4. **Provider Registration**:
   - Register providers in the appropriate module
   - Use constructor functions for dependency injection
   - Keep providers focused and single-purpose

## Running the Example

1. Navigate to the example directory:
   ```bash
   cd examples/module_example
   ```

2. Run the example:
   ```bash
   go run main.go
   ```

3. The example will:
   - Create and register modules
   - Initialize the application
   - Make a test request to `/users`
   - Clean up and shut down

## Expected Output

```
Response: {"users":["user1","user2","user3"]}
```

## Next Steps

1. Explore the example code in detail
2. Try creating your own modules
3. Experiment with different module configurations
4. Add more complex functionality to the modules 