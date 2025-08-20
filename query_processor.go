package fj

import (
	"sync"
)

// PathHandler defines the interface for handling different types of path components
// in the Chain of Responsibility pattern
type PathHandler interface {
	// CanHandle determines if this handler can process the given path segment
	CanHandle(segment string) bool
	
	// Handle processes the path segment and returns the result
	Handle(ctx Context, segment string, parser *parser) (Context, bool)
	
	// SetNext sets the next handler in the chain
	SetNext(handler PathHandler)
	
	// GetNext returns the next handler in the chain
	GetNext() PathHandler
}

// BasePathHandler provides common functionality for all path handlers
type BasePathHandler struct {
	next PathHandler
	mu   sync.RWMutex
}

// SetNext sets the next handler in the chain (thread-safe)
func (h *BasePathHandler) SetNext(handler PathHandler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.next = handler
}

// GetNext returns the next handler in the chain (thread-safe)
func (h *BasePathHandler) GetNext() PathHandler {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.next
}

// HandleNext passes processing to the next handler in the chain
func (h *BasePathHandler) HandleNext(ctx Context, segment string, parser *parser) (Context, bool) {
	h.mu.RLock()
	next := h.next
	h.mu.RUnlock()
	
	if next != nil {
		return next.Handle(ctx, segment, parser)
	}
	return Context{}, false
}

// ObjectPropertyHandler handles object property access like "user.name"
type ObjectPropertyHandler struct {
	BasePathHandler
}

// CanHandle checks if the segment represents object property access
func (h *ObjectPropertyHandler) CanHandle(segment string) bool {
	if isEmptyString(segment) {
		return false
	}
	// Simple property access - no special characters except dots and alphanumeric
	for _, r := range segment {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
			 (r >= '0' && r <= '9') || r == '_' || r == '.' || r == '-') {
			return false
		}
	}
	return true
}

// Handle processes object property access
func (h *ObjectPropertyHandler) Handle(ctx Context, segment string, parser *parser) (Context, bool) {
	if !h.CanHandle(segment) {
		return h.HandleNext(ctx, segment, parser)
	}
	
	if ctx.kind != JSON {
		return Context{}, false
	}
	
	// Use existing Get functionality for compatibility
	result := ctx.Get(segment)
	return result, result.Exists()
}

// ArrayElementHandler handles array element access like "users[0]" or "items[*]"  
type ArrayElementHandler struct {
	BasePathHandler
}

// CanHandle checks if the segment represents array access
func (h *ArrayElementHandler) CanHandle(segment string) bool {
	if isEmptyString(segment) {
		return false
	}
	// Look for array access patterns [index], [*], or array queries
	return len(segment) > 2 && (segment[len(segment)-1] == ']' || 
		containsRune(segment, '['))
}

// Handle processes array element access
func (h *ArrayElementHandler) Handle(ctx Context, segment string, parser *parser) (Context, bool) {
	if !h.CanHandle(segment) {
		return h.HandleNext(ctx, segment, parser)
	}
	
	// Use existing Get functionality for array access
	result := ctx.Get(segment)
	return result, result.Exists()
}

// WildcardHandler handles wildcard queries like "users.*.name"
type WildcardHandler struct {
	BasePathHandler
}

// CanHandle checks if the segment contains wildcard characters
func (h *WildcardHandler) CanHandle(segment string) bool {
	if isEmptyString(segment) {
		return false
	}
	return containsRune(segment, '*') || containsRune(segment, '?')
}

// Handle processes wildcard queries
func (h *WildcardHandler) Handle(ctx Context, segment string, parser *parser) (Context, bool) {
	if !h.CanHandle(segment) {
		return h.HandleNext(ctx, segment, parser)
	}
	
	// Use existing Get functionality for wildcard processing
	result := ctx.Get(segment)
	return result, result.Exists()
}

// ConditionalHandler handles conditional queries like "users[age>25].name"
type ConditionalHandler struct {
	BasePathHandler
}

// CanHandle checks if the segment contains conditional query syntax
func (h *ConditionalHandler) CanHandle(segment string) bool {
	if isEmptyString(segment) {
		return false
	}
	// Look for conditional operators in brackets
	if !containsRune(segment, '[') || !containsRune(segment, ']') {
		return false
	}
	// Check for conditional operators
	return containsAny(segment, []string{"=", "!=", "<", ">", "<=", ">=", "%"})
}

// Handle processes conditional queries
func (h *ConditionalHandler) Handle(ctx Context, segment string, parser *parser) (Context, bool) {
	if !h.CanHandle(segment) {
		return h.HandleNext(ctx, segment, parser)
	}
	
	// Use existing Get functionality for conditional processing  
	result := ctx.Get(segment)
	return result, result.Exists()
}

// HandlerChain manages the chain of path handlers
type HandlerChain struct {
	head PathHandler
	mu   sync.RWMutex
}

// NewHandlerChain creates a new handler chain with default handlers
func NewHandlerChain() *HandlerChain {
	// Create handlers
	objectHandler := &ObjectPropertyHandler{}
	arrayHandler := &ArrayElementHandler{}
	wildcardHandler := &WildcardHandler{}
	conditionalHandler := &ConditionalHandler{}
	
	// Chain them together - order matters for precedence
	conditionalHandler.SetNext(arrayHandler)
	arrayHandler.SetNext(wildcardHandler)
	wildcardHandler.SetNext(objectHandler)
	
	return &HandlerChain{
		head: conditionalHandler,
	}
}

// Process processes a path segment through the handler chain
func (c *HandlerChain) Process(ctx Context, segment string, parser *parser) (Context, bool) {
	c.mu.RLock()
	head := c.head
	c.mu.RUnlock()
	
	if head != nil {
		return head.Handle(ctx, segment, parser)
	}
	return Context{}, false
}

// SetHead sets the head of the handler chain (thread-safe)
func (c *HandlerChain) SetHead(handler PathHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.head = handler
}