// Just some basic example usage of the library
package main

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	if !htmx.IsHTMXRequest(r) {
		w.Write([]byte("only HTMX requests allowed"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	writer := htmx.NewResponse().
		Reswap(htmx.SwapBeforeBegin).
		Redirect("/cats").
		LocationWithContext("/hello", htmx.LocationContext{
			Target: "#testdiv",
			Source: "HELLO",
		}).
		Refresh(false)

	writer.Write(w)
}
