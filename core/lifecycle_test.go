package core

import (
	"context"
	"sync"
	"testing"
	"time"
)

// MockModule là một module giả sử dụng cho kiểm thử
type MockModule struct {
	moduleInitCalled    bool
	appBootstrapCalled  bool
	appShutdownCalled   bool
	moduleDestroyCalled bool
	initError           error
	bootstrapError      error
	shutdownError       error
	destroyError        error
	moduleInitOrder     int
	appBootstrapOrder   int
	appShutdownOrder    int
	moduleDestroyOrder  int
	mu                  sync.Mutex
}

// GetMetadata trả về metadata của module
func (m *MockModule) GetMetadata() ModuleMetadata {
	return ModuleMetadata{
		Imports:     make([]Module, 0),
		Exports:     make([]interface{}, 0),
		Providers:   make([]interface{}, 0),
		Controllers: make([]interface{}, 0),
	}
}

// OnModuleInit được gọi khi module khởi tạo
func (m *MockModule) OnModuleInit(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.moduleInitCalled = true
	return m.initError
}

// OnApplicationBootstrap được gọi khi ứng dụng khởi động
func (m *MockModule) OnApplicationBootstrap(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.appBootstrapCalled = true
	return m.bootstrapError
}

// OnApplicationShutdown được gọi khi ứng dụng shutdown
func (m *MockModule) OnApplicationShutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.appShutdownCalled = true
	return m.shutdownError
}

// OnModuleDestroy được gọi khi module bị hủy
func (m *MockModule) OnModuleDestroy(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.moduleDestroyCalled = true
	return m.destroyError
}

// MockProvider là một provider giả sử dụng cho kiểm thử
type MockProvider struct {
	moduleInitCalled    bool
	appBootstrapCalled  bool
	appShutdownCalled   bool
	moduleDestroyCalled bool
	initError           error
	bootstrapError      error
	shutdownError       error
	destroyError        error
	mu                  sync.Mutex
}

// OnModuleInit được gọi khi module khởi tạo
func (p *MockProvider) OnModuleInit(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.moduleInitCalled = true
	return p.initError
}

// OnApplicationBootstrap được gọi khi ứng dụng khởi động
func (p *MockProvider) OnApplicationBootstrap(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.appBootstrapCalled = true
	return p.bootstrapError
}

// OnApplicationShutdown được gọi khi ứng dụng shutdown
func (p *MockProvider) OnApplicationShutdown(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.appShutdownCalled = true
	return p.shutdownError
}

// OnModuleDestroy được gọi khi module bị hủy
func (p *MockProvider) OnModuleDestroy(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.moduleDestroyCalled = true
	return p.destroyError
}

func TestLifecycleManager_RunModuleInit(t *testing.T) {
	manager := NewLifecycleManager()
	ctx := context.Background()

	mockModule := &MockModule{}
	mockProvider := &MockProvider{}

	manager.RegisterModules([]Module{mockModule})
	manager.RegisterProviders([]interface{}{mockProvider})

	if err := manager.RunModuleInit(ctx); err != nil {
		t.Fatalf("RunModuleInit should not return error, got: %v", err)
	}

	if !mockModule.moduleInitCalled {
		t.Error("OnModuleInit was not called on module")
	}

	if !mockProvider.moduleInitCalled {
		t.Error("OnModuleInit was not called on provider")
	}

	if manager.GetState() != StateModuleInit {
		t.Errorf("Expected state to be StateModuleInit, got: %v", manager.GetState())
	}
}

func TestLifecycleManager_RunAppBootstrap(t *testing.T) {
	manager := NewLifecycleManager()
	ctx := context.Background()

	mockModule := &MockModule{}
	mockProvider := &MockProvider{}

	manager.RegisterModules([]Module{mockModule})
	manager.RegisterProviders([]interface{}{mockProvider})

	if err := manager.RunAppBootstrap(ctx); err != nil {
		t.Fatalf("RunAppBootstrap should not return error, got: %v", err)
	}

	if !mockModule.appBootstrapCalled {
		t.Error("OnApplicationBootstrap was not called on module")
	}

	if !mockProvider.appBootstrapCalled {
		t.Error("OnApplicationBootstrap was not called on provider")
	}

	if manager.GetState() != StateRunning {
		t.Errorf("Expected state to be StateRunning, got: %v", manager.GetState())
	}
}

func TestLifecycleManager_RunAppShutdown(t *testing.T) {
	manager := NewLifecycleManager()
	ctx := context.Background()

	mockModule := &MockModule{}
	mockProvider := &MockProvider{}

	shutdownHookCalled := false
	shutdownHook := func(ctx context.Context) error {
		shutdownHookCalled = true
		return nil
	}

	manager.RegisterModules([]Module{mockModule})
	manager.RegisterProviders([]interface{}{mockProvider})
	manager.RegisterShutdownHook(shutdownHook)

	if err := manager.RunAppShutdown(ctx); err != nil {
		t.Fatalf("RunAppShutdown should not return error, got: %v", err)
	}

	if !mockModule.appShutdownCalled {
		t.Error("OnApplicationShutdown was not called on module")
	}

	if !mockProvider.appShutdownCalled {
		t.Error("OnApplicationShutdown was not called on provider")
	}

	if !shutdownHookCalled {
		t.Error("Shutdown hook was not called")
	}

	if manager.GetState() != StateAppShutdown {
		t.Errorf("Expected state to be StateAppShutdown, got: %v", manager.GetState())
	}
}

func TestLifecycleManager_RunModuleDestroy(t *testing.T) {
	manager := NewLifecycleManager()
	ctx := context.Background()

	mockModule := &MockModule{}
	mockProvider := &MockProvider{}

	manager.RegisterModules([]Module{mockModule})
	manager.RegisterProviders([]interface{}{mockProvider})

	if err := manager.RunModuleDestroy(ctx); err != nil {
		t.Fatalf("RunModuleDestroy should not return error, got: %v", err)
	}

	if !mockModule.moduleDestroyCalled {
		t.Error("OnModuleDestroy was not called on module")
	}

	if !mockProvider.moduleDestroyCalled {
		t.Error("OnModuleDestroy was not called on provider")
	}

	if manager.GetState() != StateModuleDestroy {
		t.Errorf("Expected state to be StateModuleDestroy, got: %v", manager.GetState())
	}
}

func TestLifecycleManager_ErrorHandling(t *testing.T) {
	manager := NewLifecycleManager()
	ctx := context.Background()

	// Test error propagation in ModuleInit
	mockModuleWithError := &MockModule{initError: context.Canceled}
	manager.RegisterModules([]Module{mockModuleWithError})

	if err := manager.RunModuleInit(ctx); err != context.Canceled {
		t.Errorf("RunModuleInit should propagate module init error, got: %v", err)
	}

	// Test error propagation in AppBootstrap
	manager = NewLifecycleManager()
	mockModuleWithError = &MockModule{bootstrapError: context.Canceled}
	manager.RegisterModules([]Module{mockModuleWithError})

	if err := manager.RunAppBootstrap(ctx); err != context.Canceled {
		t.Errorf("RunAppBootstrap should propagate bootstrap error, got: %v", err)
	}
}

func TestLifecycleManager_FullLifecycle(t *testing.T) {
	manager := NewLifecycleManager()
	ctx := context.Background()

	mockModule1 := &MockModule{}
	mockModule2 := &MockModule{}
	mockProvider1 := &MockProvider{}
	mockProvider2 := &MockProvider{}

	var shutdownHookCalled bool
	shutdownHook := func(ctx context.Context) error {
		shutdownHookCalled = true
		return nil
	}

	manager.RegisterModules([]Module{mockModule1, mockModule2})
	manager.RegisterProviders([]interface{}{mockProvider1, mockProvider2})
	manager.RegisterShutdownHook(shutdownHook)

	// Execute full lifecycle
	if err := manager.RunModuleInit(ctx); err != nil {
		t.Fatalf("RunModuleInit failed: %v", err)
	}

	if err := manager.RunAppBootstrap(ctx); err != nil {
		t.Fatalf("RunAppBootstrap failed: %v", err)
	}

	// Simulate running application for a short time
	time.Sleep(100 * time.Millisecond)

	if err := manager.RunAppShutdown(ctx); err != nil {
		t.Fatalf("RunAppShutdown failed: %v", err)
	}

	if err := manager.RunModuleDestroy(ctx); err != nil {
		t.Fatalf("RunModuleDestroy failed: %v", err)
	}

	// Verify all lifecycle methods were called
	if !mockModule1.moduleInitCalled || !mockModule2.moduleInitCalled ||
		!mockProvider1.moduleInitCalled || !mockProvider2.moduleInitCalled {
		t.Error("OnModuleInit was not called on all modules/providers")
	}

	if !mockModule1.appBootstrapCalled || !mockModule2.appBootstrapCalled ||
		!mockProvider1.appBootstrapCalled || !mockProvider2.appBootstrapCalled {
		t.Error("OnApplicationBootstrap was not called on all modules/providers")
	}

	if !mockModule1.appShutdownCalled || !mockModule2.appShutdownCalled ||
		!mockProvider1.appShutdownCalled || !mockProvider2.appShutdownCalled {
		t.Error("OnApplicationShutdown was not called on all modules/providers")
	}

	if !mockModule1.moduleDestroyCalled || !mockModule2.moduleDestroyCalled ||
		!mockProvider1.moduleDestroyCalled || !mockProvider2.moduleDestroyCalled {
		t.Error("OnModuleDestroy was not called on all modules/providers")
	}

	if !shutdownHookCalled {
		t.Error("Shutdown hook was not called")
	}

	if manager.GetState() != StateModuleDestroy {
		t.Errorf("Expected final state to be StateModuleDestroy, got: %v", manager.GetState())
	}
}
