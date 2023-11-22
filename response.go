package htmx

import (
	"context"
	"errors"
	"net/http"
)

// Response contains HTMX headers to write to a response.
type Response struct {
	// The HTMX headers that will be written to a response.
	headers map[string]string

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
	if r.locationWithContextErr != nil {
		return errors.Join(r.locationWithContextErr...)
	}

	header := w.Header()
	for k, v := range r.headers {
		header.Set(k, v)
	}

	return nil
}

// Headers returns a copy of the headers. Any modifications to the
// returned headers will not affect the headers in this struct.
func (r Response) Headers() map[string]string {
	m := make(map[string]string)
	for k, v := range r.headers {
		m[k] = v
	}
	return m
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
