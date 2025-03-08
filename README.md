# Goblin Framework

[![GoDoc](https://godoc.org/github.com/calummacc/goblin?status.svg)](https://godoc.org/github.com/calummacc/goblin)  
[![Go Report Card](https://goreportcard.com/badge/github.com/calummacc/goblin)](https://goreportcard.com/report/github.com/calummacc/goblin)

**A modular Go framework**, leveraging **Uber's `fx` and `dig`** for dependency injection, lifecycle management, and modularity. Built for building scalable, maintainable web applications with a modern architecture.

---

## Key Features

- **Modular Design**: Structure applications using reusable modules (e.g., `user`, `auth`).
- **Dependency Injection (DI)**: Powered by **`fx`** (Uber's framework) and **`dig`**, enabling automatic wiring of dependencies.
- **Router & Middleware**: Built on **Gin** for high-performance HTTP routing and middleware (logging, authentication, error handling).
- **Lifecycle Hooks**: Manage app startup/shutdown with `fx`-based lifecycle events.
- **CLI Tool**: Generate projects, modules, controllers, and services with `cobra`.
- **Extensible Patterns**: Support for Guards, Interceptors, Exception Filters, and Validation (Pipes).
- **Database Integration**: Optional ORM support via **GORM**.

---

## Getting Started

### Installation

```bash
go get github.com/calummacc/goblin
```

### Create a New Project

```bash
goblin new my-goblin-app
cd my-goblin-app
```

### Run the Application

```bash
go run main.go
```

Your app will start on `http://localhost:8485`.

---

## Usage Example

### 1. Define a Module

Create a `user` module with a service and controller:

```go
// internal/modules/user/user_service.go
package user

type UserService struct {}

func (s *UserService) GetUser() string {
    return "Hello from UserService!"
}
```

```go
// internal/modules/user/user_controller.go
package user

import "github.com/gin-gonic/gin"

type UserController struct {}

func (c UserController) GetUserHandler(ctx *gin.Context) {
    ctx.String(200, "User data")
}
```

### 2. Register the Module

```go
// internal/modules/user/user_module.go
package user

import "go.uber.org/fx"

var Module = fx.Options(
    fx.Provide(
        NewUserService,
        NewUserController,
    ),
)
```

### 3. Define Routes

```go
// internal/core/router.go
package core

import "github.com/gin-gonic/gin"

func NewGinEngine() *gin.Engine {
    r := gin.Default()
    r.GET("/user", func(c *gin.Context) {
        // Inject UserController here (TODO: Implement DI)
        c.String(200, "User endpoint")
    })
    return r
}
```

---

## CLI Commands

### Generate a New Project

```bash
goblin new [project-name]
```

### Generate a Module

```bash
goblin generate module [module-name]
```

### Generate a Controller/Service

```bash
goblin generate controller [module-name].user
goblin generate service [module-name].user
```

---

## Architecture

### Project Structure

```
goblin/
├── cmd/               # CLI tool (built with Cobra)
├── internal/
│   ├── core/          # Core framework logic (DI, router, lifecycle)
│   └── modules/       # Business modules (e.g., user, auth)
├── examples/          # Example projects
├── go.mod             # Dependency management
└── main.go            # Entry point
```

### Dependency Injection Flow

1. **Providers** define objects (services, repositories) via `fx.Provide()`.
2. **Modules** group providers (e.g., `user.Module`).
3. **App Module** aggregates all modules in `core/di.go`.
4. **Lifecycle Hooks** start/stop the app via `fx.Hook`.

---

## Roadmap

- **Phase 1 (Completed)**: Core skeleton, CLI, DI, router.
- **Phase 2**: Guards, Interceptors, Exception Filters, GORM integration.
- **Phase 3**: Performance optimizations (e.g., `fasthttp`), CI/CD, docs.

---

## Contributing

1. Fork the repo and create your branch.
2. Implement features or fixes.
3. Add tests for new functionality.
4. Submit a pull request.

---

## License

MIT License  
Copyright (c) 2023 Your Name

---

**Join the Goblin Framework community** on [GitHub Discussions](https://github.com/calummacc/goblin/discussions) or [Slack](https://goblin-framework.slack.com).

---

Let me know if you need adjustments or additional sections!
