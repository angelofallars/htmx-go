// Just some basic example usage of the library
package main

import (
	"fmt"
	"net/http"

	"github.com/angelofallars/htmx-go"
)

func main() {
	r := htmx.NewResponse().
		AddTriggerAfterSwap(
			htmx.Trigger("myEvent"),
			htmx.Trigger("cat"),
			htmx.TriggerValue("myEvent2", "weee"),
		).
		AddTriggerAfterSettle(
			htmx.TriggerValue("myEvent2", "weee"),
		).
		AddTrigger(
			htmx.TriggerKeyValue("myEvent4", map[string]string{
				"level":   "info",
				"message": "Here is a Message",
			}),
			htmx.TriggerKeyValue("myEvent3", map[string]string{
				"level":   "info",
				"message": "Here is a Message",
			}),
		).
		AddTrigger(htmx.Trigger("myEvent4"))

	fmt.Println(r)
	fmt.Println(r.Headers())
}

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
