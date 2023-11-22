// ext.go provides this library with interfaces
// for integrations with popular Go libraries.
package htmx

import (
	"context"
	"io"
)

// Interface to integrate with Templ components
type TemplComponent interface {
	// Render the template.
	Render(ctx context.Context, w io.Writer) error
}
