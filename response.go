package htmx

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

// Response contains HTMX headers to write to a response.
type Response struct {
	// The HTMX headers that will be written to a response.
	headers map[string]string

	// The HTTP status code to use
	statusCode int

	// Triggers for 'HX-Trigger'
	triggers []EventTrigger

	// Triggers for 'HX-Trigger-After-Settle'
	triggersAfterSettle []EventTrigger

	// Triggers for 'HX-Trigger-After-Swap'
	triggersAfterSwap []EventTrigger

	// JSON marshalling might fail, so we need to keep track of this error
	// to return when `Write` is called
	locationWithContextErr []error
}

// NewResponse returns a new HTMX response header writer.
//
// Any subsequent method calls that write to the same header
// will overwrite the last set header value.
func NewResponse() Response {
	return Response{
		headers: make(map[string]string),
	}
}

// Clone returns a clone of this HTMX response writer, preventing any mutation
// on the original response.
func (r Response) Clone() Response {
	n := NewResponse()

	for k, v := range r.headers {
		n.headers[k] = v
	}

	return n
}

// Write applies the defined HTMX headers to a given response writer.
func (r Response) Write(w http.ResponseWriter) error {
	if len(r.locationWithContextErr) > 0 {
		return errors.Join(r.locationWithContextErr...)
	}

	headers, err := r.Headers()
	if err != nil {
		return err
	}

	headerWriter := w.Header()
	for k, v := range headers {
		headerWriter.Set(k, v)
	}

	// Status code needs to be written after the other headers
	// so the other headers can be written
	if r.statusCode != 0 {
		w.WriteHeader(r.statusCode)
	}

	return nil
}

// RenderHTML renders an HTML document fragment along with the defined HTMX headers.
func (r Response) RenderHTML(w http.ResponseWriter, html template.HTML) (int, error) {
	err := r.Write(w)
	if err != nil {
		return 0, err
	}

	return w.Write([]byte(html))
}

// RenderTempl renders a Templ component along with the defined HTMX headers.
func (r Response) RenderTempl(ctx context.Context, w http.ResponseWriter, c templComponent) error {
	err := r.Write(w)
	if err != nil {
		return err
	}

	err = c.Render(ctx, w)
	if err != nil {
		return err
	}

	return nil
}

// MustWrite applies the defined HTMX headers to a given response writer, otherwise it panics.
//
// Under the hood this uses [Response.Write].
func (r Response) MustWrite(w http.ResponseWriter) {
	err := r.Write(w)
	if err != nil {
		panic(err)
	}
}

// MustRenderHTML renders an HTML document fragment along with the defined HTMX headers, otherwise it panics.
//
// Under the hood this uses [Response.RenderHTML].
func (r Response) MustRenderHTML(w http.ResponseWriter, html template.HTML) {
	_, err := r.RenderHTML(w, html)
	if err != nil {
		panic(err)
	}
}

// MustRenderTempl renders a Templ component along with the defined HTMX headers, otherwise it panics.
//
// Under the hood this uses [Response.RenderTempl].
func (r Response) MustRenderTempl(ctx context.Context, w http.ResponseWriter, c templComponent) {
	err := r.RenderTempl(ctx, w, c)
	if err != nil {
		panic(err)
	}
}

// Headers returns a copied map of the headers. Any modifications to the
// returned headers will not affect the headers in this struct.
func (r Response) Headers() (map[string]string, error) {
	m := make(map[string]string)

	for k, v := range r.headers {
		m[k] = v
	}

	if r.triggers != nil {
		triggers, err := triggersToString(r.triggers)
		if err != nil {
			return nil, fmt.Errorf("marshalling triggers failed: %w", err)
		}
		m[HeaderTrigger] = triggers
	}

	if r.triggersAfterSettle != nil {
		triggers, err := triggersToString(r.triggersAfterSettle)
		if err != nil {
			return nil, fmt.Errorf("marshalling triggers after settle failed: %w", err)
		}
		m[HeaderTriggerAfterSettle] = triggers
	}

	if r.triggersAfterSwap != nil {
		triggers, err := triggersToString(r.triggersAfterSwap)
		if err != nil {
			return nil, fmt.Errorf("marshalling triggers after swap failed: %w", err)
		}
		m[HeaderTriggerAfterSwap] = triggers
	}

	return m, nil
}
