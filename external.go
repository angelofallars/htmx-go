package htmx

import (
	"context"
	"io"
)

// Interface to integrate with Templ components
type templComponent interface {
	// Render the template.
	Render(ctx context.Context, w io.Writer) error
}
