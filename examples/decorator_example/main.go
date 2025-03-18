package decorator_example

import (
	"fmt"
	"goblin/decorators"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Create a new product controller
	controller := NewProductController()

	// Parse route metadata from controller
	routesInfo := decorators.ParseRouteMetadata(controller)

	// Print route information
	fmt.Println("Routes:")
	for methodName, routeInfo := range routesInfo {
		fmt.Printf("- %s %s -> %s.%s\n",
			routeInfo.Method,
			routeInfo.Path,
			"ProductController",
			methodName,
		)

		// Print additional metadata
		if routeInfo.Description != "" {
			fmt.Printf("  Description: %s\n", routeInfo.Description)
		}

		if len(routeInfo.Tags) > 0 {
			fmt.Printf("  Tags: %v\n", routeInfo.Tags)
		}

		if routeInfo.IsPublic {
			fmt.Printf("  Public: true\n")
		}

		if routeInfo.Deprecated {
			fmt.Printf("  Deprecated: true\n")
		}

		if len(routeInfo.Guards) > 0 {
			fmt.Printf("  Guards: %d\n", len(routeInfo.Guards))
		}

		if len(routeInfo.Middleware) > 0 {
			fmt.Printf("  Middleware: %d\n", len(routeInfo.Middleware))
		}

		fmt.Println()
	}

	// Convert route info to Gin routes
	routes := decorators.ConvertToRoutes("/products", routesInfo)

	// Register routes with Gin
	for _, route := range routes {
		router.Handle(route.Method, route.Path, append(route.Middlewares, route.Handler)...)
	}

	// Start the server
	log.Println("Starting server on http://localhost:8080")
	router.Run(":8080")
}
