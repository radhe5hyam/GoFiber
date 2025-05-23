package test

import (
	"testing"

	"github.com/radhe5hyam/GoFiber/http/router"
)

func Test_AddandFindSegment(t *testing.T) {
	root := router.NewNode("/")

	dummy := func(params map[string]string) {}

	root.AddSegment("GET", "/users", dummy)
	root.AddSegment("GET", "/posts", dummy)
	root.AddSegment("GET", "/users/42/posts", dummy)

	handler, _, _, _ := root.FindSegment("GET", "/users")
	if handler == nil {
		t.Error("Expected match for /users")
	}
	handler, _, _, _ = root.FindSegment("GET", "/users/42")
	if handler != nil {
		t.Error("Did not expect match for /users/42")
	}
	handler, _, _, _ = root.FindSegment("GET", "/users/42/posts")
	if handler == nil {
		t.Error("Expected match for /users/42/posts")
	}

	root.AddSegment("GET", "/users/:id", func(params map[string]string) {
		if params["id"] != "123" {
			t.Errorf("Expected id=123, got %s", params["id"])
		}
	})

	handler, params, _, _ := root.FindSegment("GET", "/users/123")
	if handler == nil {
		t.Fatal("Expected handler match")
	}
	handler(params)


	getCalled := false
	postCalled := false

	root.AddSegment("GET", "/users/:id", func(params map[string]string) {
		getCalled = true
		if params["id"] != "123" {
			t.Errorf("Expected id=123, got %s", params["id"])
		}
	})
	handler, params, found, allowed := root.FindSegment("GET", "/users/123")
	if !found || !allowed || handler == nil {
		t.Fatal("Expected GET handler for /users/123")
	}
	handler(params)
	if !getCalled {
		t.Error("GET handler was not called")
	}


	root.AddSegment("POST", "/users/:id", func(params map[string]string) {
		postCalled = true
	})
	handler, _, found, allowed = root.FindSegment("POST", "/users/123")
	if !found || !allowed || handler == nil {
		t.Fatal("Expected POST handler")
	}
	handler(nil)
	if !postCalled {
		t.Error("POST handler was not called")
	}

	
	handler, _, found, allowed = root.FindSegment("PUT", "/users/123")
	if !found {
		t.Error("Expected path to exist")
	}
	if allowed {
		t.Error("PUT should not be allowed")
	}
	if handler != nil {
		t.Error("Handler should be nil for unsupported method")
	}
}