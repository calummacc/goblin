# Goblin Framework

Goblin is a lightweight, modular Go framework built on top of Gin and Uber-FX, designed to help developers create scalable and maintainable web applications using clean architecture principles.

## Features

- 🚀 **Modular Architecture**: Easy-to-use module system for organizing your application 
- 💉 **Dependency Injection**: Built-in DI container powered by Uber-FX
- 🛠️ **Clean Architecture**: Follows clean architecture principles with Repository, Service, and Controller patterns
- 🔌 **Middleware Support**: Pre-built middleware for logging, recovery, and error handling
- ⚙️ **Configuration Management**: Flexible configuration system with JSON support
- 🔒 **Thread-Safe**: Concurrent-safe operations with proper mutex implementation
- 📝 **Logging**: Built-in request logging and error tracking
- 🎯 **Request ID Tracking**: Unique ID generation for request tracing

## Installation

```bash
go get github.com/onepiecehung/goblin
```

## Quick Start

### 1. Create a new project

```bash
mkdir myapp
cd myapp
go mod init myapp
```

### 2. Create basic module structure

```go
// main.go
package main

import (
    "context"
    "github.com/onepiecehung/goblin/internal/core"
    "log"
)

func main() {
    app := core.NewApplication()
    appModule := NewAppModule()
    app.AddModule(appModule)
    app.Configure()
    
    if err := app.Run(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

### 3. Create a module

```go
// user/user.module.go
package user

type UserModule struct {
    core.BaseModule
    controller *Controller
    service    Service
    repository Repository
}

func (m *UserModule) RegisterRoutes(router *gin.RouterGroup) {
    users := router.Group("/users")
    {
        users.GET("", m.controller.GetUsers)
        users.POST("", m.controller.CreateUser)
    }
}
```

## Architecture

```
myapp/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── core/
│   │   ├── application.go
│   │   ├── module.go
│   │   └── container.go
│   └── middleware/
│       ├── logger.go
│       └── recovery.go
├── pkg/
│   └── config/
│       └── config.go
└── modules/
    └── user/
        ├── user.module.go
        ├── user.controller.go
        ├── user.service.go
        ├── user.repository.go
        └── user.model.go
```

## Middleware

Built-in middleware:

- Logger: Request logging
- Recovery: Panic recovery
- RequestID: Request tracking
- ErrorHandler: Centralized error handling

## Example

Creating a simple user module:

```go
// user/controller.go
func (c *Controller) GetUsers(ctx *gin.Context) {
    users, err := c.service.GetAllUsers()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, users)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

- DS112 (@ds112)

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Uber-FX](https://github.com/uber-go/fx)