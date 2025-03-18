package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"goblin/core"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserController(t *testing.T) {
	// Tắt Gin debug mode
	gin.SetMode(gin.TestMode)

	// Tạo Gin engine
	engine := gin.New()

	// Tạo controller manager
	cm := core.NewControllerManager(engine)

	// Đăng ký controller
	controller := NewUserController()
	err := cm.RegisterController(controller)
	assert.NoError(t, err)

	// Test GET /users
	t.Run("GetUsers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response []User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
	})

	// Test GET /users/:id
	t.Run("GetUser", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/1", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.ID)
	})

	// Test POST /users
	t.Run("CreateUser", func(t *testing.T) {
		user := User{Name: "Test User"}
		body, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		engine.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)

		var response User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, response.Name)
	})

	// Test PUT /users/:id
	t.Run("UpdateUser", func(t *testing.T) {
		user := User{Name: "Updated User"}
		body, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		engine.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, response.Name)
	})

	// Test DELETE /users/:id
	t.Run("DeleteUser", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "User 1 deleted")
	})

	// Test Guard
	t.Run("AuthGuard", func(t *testing.T) {
		// Test without token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, 403, w.Code)

		// Test with token
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/users", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		engine.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})
}
