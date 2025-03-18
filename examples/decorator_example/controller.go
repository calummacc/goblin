package decorator_example

import (
	"fmt"
	"goblin/core"
	"goblin/decorators"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Product represents a product in our system
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

// AdminGuard is a guard that only allows admins
type AdminGuard struct {
	core.BaseGuard
}

// CanActivate checks if the user is an admin
func (g *AdminGuard) CanActivate(ctx *core.Context) (bool, error) {
	// Get the Authorization header
	token := ctx.GinContext.GetHeader("Authorization")

	// In a real app, you'd verify the token
	// This is just a simple example
	return token == "admin-token", nil
}

// LoggingMiddleware is a middleware that logs requests
func LoggingMiddleware(c *gin.Context) {
	fmt.Printf("[%s] %s\n", c.Request.Method, c.Request.URL.Path)
	c.Next()
}

// ProductController is a controller for managing products
type ProductController struct {
	core.BaseController
	products map[string]Product
}

// NewProductController creates a new product controller
func NewProductController() *ProductController {
	// Create some sample products
	products := map[string]Product{
		"1": {
			ID:          "1",
			Name:        "Laptop",
			Description: "A powerful laptop for developers",
			Price:       1299.99,
			Category:    "Electronics",
		},
		"2": {
			ID:          "2",
			Name:        "Smartphone",
			Description: "A high-end smartphone",
			Price:       899.99,
			Category:    "Electronics",
		},
		"3": {
			ID:          "3",
			Name:        "Headphones",
			Description: "Noise-cancelling headphones",
			Price:       299.99,
			Category:    "Accessories",
		},
	}

	controller := &ProductController{
		products: products,
	}

	// Initialize controller
	applyControllerDecorators(controller)

	return controller
}

// Apply decorators to controller
func applyControllerDecorators(controller *ProductController) {
	// Controller base path
	decorators.Controller("/products")(controller)

	// Route decorators
	decorators.Get("/")(controller, "GetProducts")
	decoratorsDesc, ok := decorators.Description("Get all products").(struct {
		ControllerDecorator func(interface{})
		MethodDecorator     func(interface{}, string)
	})
	if ok {
		decoratorsDesc.MethodDecorator(controller, "GetProducts")
	}

	decoratorsTags, ok := decorators.Tags("products").(struct {
		ControllerDecorator func(interface{})
		MethodDecorator     func(interface{}, string)
	})
	if ok {
		decoratorsTags.MethodDecorator(controller, "GetProducts")
	}

	decorators.Public()(controller, "GetProducts")
	decorators.UseMiddleware(LoggingMiddleware)(controller, "GetProducts")

	// GetProduct decorators
	decorators.Get("/:id")(controller, "GetProduct")
	if ok {
		decoratorsDesc.MethodDecorator(controller, "GetProduct")
		decoratorsTags.MethodDecorator(controller, "GetProduct")
	}
	decorators.Public()(controller, "GetProduct")

	// CreateProduct decorators
	decorators.Post("/")(controller, "CreateProduct")
	decoratorsCreateDesc, ok := decorators.Description("Create a new product").(struct {
		ControllerDecorator func(interface{})
		MethodDecorator     func(interface{}, string)
	})
	if ok {
		decoratorsCreateDesc.MethodDecorator(controller, "CreateProduct")
	}

	decoratorsAdminTags, ok := decorators.Tags("products", "admin").(struct {
		ControllerDecorator func(interface{})
		MethodDecorator     func(interface{}, string)
	})
	if ok {
		decoratorsAdminTags.MethodDecorator(controller, "CreateProduct")
	}

	decorators.UseGuards(&AdminGuard{})(controller, "CreateProduct")

	// UpdateProduct decorators
	decorators.Put("/:id")(controller, "UpdateProduct")
	decoratorsUpdateDesc, ok := decorators.Description("Update an existing product").(struct {
		ControllerDecorator func(interface{})
		MethodDecorator     func(interface{}, string)
	})
	if ok {
		decoratorsUpdateDesc.MethodDecorator(controller, "UpdateProduct")
	}

	if ok {
		decoratorsAdminTags.MethodDecorator(controller, "UpdateProduct")
	}

	decorators.UseGuards(&AdminGuard{})(controller, "UpdateProduct")

	// DeleteProduct decorators
	decorators.Delete("/:id")(controller, "DeleteProduct")
	decoratorsDeleteDesc, ok := decorators.Description("Delete a product").(struct {
		ControllerDecorator func(interface{})
		MethodDecorator     func(interface{}, string)
	})
	if ok {
		decoratorsDeleteDesc.MethodDecorator(controller, "DeleteProduct")
	}

	if ok {
		decoratorsAdminTags.MethodDecorator(controller, "DeleteProduct")
	}

	decorators.UseGuards(&AdminGuard{})(controller, "DeleteProduct")
	decorators.Deprecated()(controller, "DeleteProduct")

	// SearchProducts decorators
	decorators.Get("/search")(controller, "SearchProducts")
	decoratorsSearchDesc, ok := decorators.Description("Search for products").(struct {
		ControllerDecorator func(interface{})
		MethodDecorator     func(interface{}, string)
	})
	if ok {
		decoratorsSearchDesc.MethodDecorator(controller, "SearchProducts")
	}

	if ok {
		decoratorsTags.MethodDecorator(controller, "SearchProducts")
	}

	decorators.Public()(controller, "SearchProducts")
}

// GetProducts returns all products
// @Get("/")
// @Description("Get all products")
// @Tags("products")
// @Public()
func (c *ProductController) GetProducts(ctx *gin.Context) {
	// Convert the map to a slice
	productList := make([]Product, 0, len(c.products))
	for _, product := range c.products {
		productList = append(productList, product)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": productList,
	})
}

// GetProduct returns a single product by ID
// @Get("/:id")
// @Description("Get a product by ID")
// @Tags("products")
// @Public()
func (c *ProductController) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	product, exists := c.products[id]
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Product with ID %s not found", id),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// CreateProduct creates a new product
// @Post("/")
// @Description("Create a new product")
// @Tags("products", "admin")
// @UseGuards(AdminGuard)
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var product Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Generate a new ID (simplistic approach for the example)
	product.ID = fmt.Sprintf("%d", len(c.products)+1)

	// Save the product
	c.products[product.ID] = product

	ctx.JSON(http.StatusCreated, gin.H{
		"product": product,
	})
}

// UpdateProduct updates an existing product
// @Put("/:id")
// @Description("Update an existing product")
// @Tags("products", "admin")
// @UseGuards(AdminGuard)
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	_, exists := c.products[id]
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Product with ID %s not found", id),
		})
		return
	}

	var product Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Ensure the ID from the path is used
	product.ID = id

	// Update the product
	c.products[id] = product

	ctx.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// DeleteProduct deletes a product
// @Delete("/:id")
// @Description("Delete a product")
// @Tags("products", "admin")
// @UseGuards(AdminGuard)
// @Deprecated()
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	_, exists := c.products[id]
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Product with ID %s not found", id),
		})
		return
	}

	// Delete the product
	delete(c.products, id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Product with ID %s was deleted", id),
	})
}

// SearchProducts searches for products
// @Get("/search")
// @Description("Search for products")
// @Tags("products")
// @Public()
func (c *ProductController) SearchProducts(ctx *gin.Context) {
	query := ctx.Query("q")
	category := ctx.Query("category")

	// Create a filtered list
	var filtered []Product
	for _, product := range c.products {
		// Apply filters
		matchesQuery := query == "" || strings.Contains(
			strings.ToLower(product.Name+" "+product.Description),
			strings.ToLower(query),
		)
		matchesCategory := category == "" || product.Category == category

		if matchesQuery && matchesCategory {
			filtered = append(filtered, product)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": filtered,
	})
}
