package helpers

import "context"

type copyContext struct {
	// This context contains the cancel/timeout/todo/background channel
	context.Context
	// This context contains the values
	values context.Context
}

// NewCopyContext is function to duplicate context value
func NewCopyContext(baseContext, previousContext context.Context) context.Context {
	return &copyContext{
		Context: baseContext,
		values:  previousContext,
	}
}

// Value is function to insert context value
func (c *copyContext) Value(key any) any {
	return c.values.Value(key)
}
