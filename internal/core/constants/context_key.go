package constants

// ContextKey is a custom type to avoid key collisions in context values
type ContextKey string

// Context key constants for get values from context
const (
	// XRequestIDKey is a key for x request id on context with standard pattern
	XRequestIDKey ContextKey = "X-Request-Id"
)

// HTTP Header constants
const (
	// XRequestIDHeader define request ID HTTP with standard header
	XRequestIDHeader = "X-Request-Id"
)
