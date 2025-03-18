# Exception Filters in Goblin Framework

Exception Filters provide a centralized way to handle exceptions in a Goblin Framework application. This example demonstrates how to use and implement exception filters.

## Concepts

### Exception Filters

Exception Filters are responsible for processing any unhandled exceptions in your application. They allow you to:

- Implement custom exception handling logic
- Control the exact flow of control and server responses after an exception occurs
- Map exceptions to HTTP responses in a consistent way
- Add logging or other cross-cutting concerns

### Filter Types

The framework provides several filter types:

- **HttpException**: A built-in exception type with status code and message
- **ExceptionFilter**: The base interface for all exception filters
- **BaseExceptionFilter**: A simple implementation for extension
- **LifecycleExceptionFilter**: Extends filters with lifecycle hooks

## Implementation

An Exception Filter must implement the `ExceptionFilter` interface:

```go
type ExceptionFilter interface {
    // Catch handles an exception
    Catch(exception error, ctx *ExceptionContext)
    // CanHandle checks if this filter can handle the given exception
    CanHandle(exception error) bool
}
```

### Example Usage

This example demonstrates:

1. Creating custom exception types (`BusinessException`, `UserNotFoundException`)
2. Implementing a custom exception filter (`HttpExceptionFilter`)
3. Registering the filter globally
4. Throwing and handling different types of exceptions

## Running the Example

To run this example:

```bash
go run main.go
```

Then try these routes:

- `GET /` - Shows available test routes
- `GET /users/1` - Returns a normal user
- `GET /users/404` - Demonstrates `UserNotFoundException`
- `GET /users/400` - Demonstrates `BusinessException`
- `GET /users/500` - Demonstrates generic error handling
- `GET /users/403` - Demonstrates `HttpException`

## Integration with Goblin Framework

Exception Filters work seamlessly with other Goblin Framework features:

1. **Middleware Integration**: Filters are applied through middleware
2. **Controller-Specific Filters**: Apply filters to specific controllers or methods
3. **Global Filters**: Apply filters application-wide
4. **Prioritization**: Controller-specific filters take precedence over global ones

## Decorators

```go
// Apply filters to a specific controller or method
@UseFilters(new MyExceptionFilter())
func (c *MyController) MyMethod() { /* ... */ }

// Apply filters globally to a module
@GlobalFilters(new MyExceptionFilter())
func MyModule() { /* ... */ }
```

## Benefits

- **Separation of Concerns**: Keep exception handling logic separate from business logic
- **Reusability**: Create filters that can be reused across the application
- **Consistency**: Handle exceptions consistently across the application
- **Testability**: Easily test exception handling logic in isolation 