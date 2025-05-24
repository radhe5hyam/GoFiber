package test

import (
	"testing"

	"github.com/radhe5hyam/GoFiber/http"
)

type Context struct {
	Params map[string]string
}

//
func Test_AddandFindSegment(t *testing.T) {
	root := http.NewNode("/")

	dummy := func(ctx *http.Context) {}

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

	root.AddSegment("GET", "/users/:id", func(ctx *http.Context) {
		if ctx.Params["id"] != "123" {
			t.Errorf("Expected id=123, got %s", ctx.Params["id"])
		}
	})

	handler, params, _, _ := root.FindSegment("GET", "/users/123")
	if handler == nil {
		t.Fatal("Expected handler match")
	}
	handler(&http.Context{Params: params})

	getCalled := false
	postCalled := false

	root.AddSegment("GET", "/users/:id", func(ctx *http.Context) {
		getCalled = true
		if ctx.Params["id"] != "123" {
			t.Errorf("Expected id=123, got %s", ctx.Params["id"])
		}
	})
	handler, params, found, allowed := root.FindSegment("GET", "/users/123")
	if !found || !allowed || handler == nil {
		t.Fatal("Expected GET handler for /users/123")
	}
	handler(&http.Context{Params: params})
	if !getCalled {
		t.Error("GET handler was not called")
	}

	root.AddSegment("POST", "/users/:id", func(ctx *http.Context) {
		postCalled = true
	})
	handler, params, found, allowed = root.FindSegment("POST", "/users/123")
	if !found || !allowed || handler == nil {
		t.Fatal("Expected POST handler")
	}
	handler(&http.Context{Params: params})
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