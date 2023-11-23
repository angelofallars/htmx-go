<img src="https://github.com/angelofallars/htmx-go/assets/39676098/c1a14954-27fd-4276-8948-0800e5372b14" width="400px">

# htmx-go (WIP)

A **type-safe** library for working with [HTMX](https://htmx.org/) in Go. Less time fiddling with HTTP
headers, more time developing awesome Hypermedia-driven applications.

Easily check if requests are from HTMX, and utilize a type-safe, declarative syntax for HTMX response headers to control HTMX behavior from the server.

Has some integration with [templ](https://templ.guide/) components.
Uses standard `net/http` types.

```go
import (
  "net/http"
  "github.com/angelofallars/htmx-go"
)

func(w http.ResponseWriter, r *http.Request) {
	if htmx.IsHTMXRequest(r) {
		htmx.NewResponse().
			Reswap(htmx.SwapBeforeEnd).
			Retarget("#errors").
			ReplaceURL("/errors").
			Write(w)
	}
}
```

## Installation

Use go get.

```sh
go get -u github.com/angelofallars/htmx-go
```

Then import htmx-go:

```go
import "github.com/angelofallars/htmx-go"
```

## HTMX requests

### Check request origin

You can determine if a request is from HTMX. With this, you can add custom handling for non-HTMX requests.

You can also use this for checking if this is a GET request for the initial (very first) page loads on your website, as initial page load requests don't use HTMX.

```go
func(w http.ResponseWriter, r *http.Request) {
	if htmx.IsHTMXRequest(r) {
		// logic for handling HTMX requests
	} else {
		// logic for handling non-HTMX requests (e.g. render a full page for first-time visitors)
	}
}
```

### Check if request is Boosted (`hx-boost`)

```go
func(w http.ResponseWriter, r *http.Request) {
	if htmx.IsBoosted(r) {
		// logic for handling boosted requests
	} else {
		// logic for handling non-boosted requests
	}
}
```

## HTMX responses

htmx-go takes inspiration from [Lip Gloss](https://github.com/charmbracelet/lipgloss) for a
declarative way of specifying HTMX response headers.

### Basic usage

Make a response writer with `htmx.NewResponse()`, and add a header to it to make the page refresh:

``` go
func(w http.ResponseWriter, r *http.Request) {
 	writer := htmx.NewResponse().Refresh(true)
 	writer.Write(w)
}
```

### Retarget response to a different element
```go
func(w http.ResponseWriter, r *http.Request) {
	htmx.NewResponse().
		// Override 'hx-target' to specify which target to load into
		Retarget("#errors").
		// Also override the 'hx-swap' value of the request
		Reswap(htmx.SwapBeforeEnd).
		Write(w)
}
```

### Triggers

Add triggers without dealing with serializing the trigger headers.

Define event triggers:

- `htmx.Trigger(eventName string)` - A trigger with no details.
- `htmx.TriggerValue(eventName string, detail string)` - A trigger with one detail.
- `htmx.TriggerKeyValue(eventName string, details map[string]string)` - A trigger with key/value
details.

Set trigger headers using the triggers above:

- `Response.AddTrigger(trigger ..EventTrigger)` - appends to the `HX-Trigger` header
- `Response.AddTriggerAfterSettle(trigger ..EventTrigger)` - appends to the `HX-Trigger-After-Settle` header
- `Response.AddTriggerAfterSwap(trigger ..EventTrigger)` - appends to the `HX-Trigger-After-Swap` header

```go

htmx.NewResponse().
	AddTrigger(htmx.Trigger("myEvent"))
// HX-Trigger: myEvent

htmx.NewResponse().
	AddTrigger(htmx.TriggerValue("showMessage", "Here Is A Message"))
// HX-Trigger: {"showMessage":"Here Is A Message"}

htmx.NewResponse().
	AddTrigger(
			htmx.TriggerValue("hello", "world"),
			htmx.TriggerKeyValue("myEvent", map[string]string{
				"level":   "info",
				"message": "Here Is A Message",
			}),
  )
// HX-Trigger: {"hello":"world","myEvent":{"level":"info","message":"Here is a Message"}}
```

### Code organization

HTMX response writers can be declared outside of functions with `var` so you can reuse them in several
places. 

> [!CAUTION]
> If you're adding additional headers to a global response writer, always use the `.Clone()` method
> to avoid accidentally modifying the global response writer.

```go
var deleter = htmx.NewResponse().
    Reswap(htmx.SwapDelete)

func(w http.ResponseWriter, r *http.Request) {
	deleter.Clone().
		Reselect("#messages").
		Write(w)
}
```

### Templ integration

HTMX pairs well with [Templ](https://templ.guide), and this library is no exception. You can render
both the necessary HTMX response headers and Templ components in one step with the
`.RenderTempl()` method.

```go
// hello.templ
templ Hello() {
    <div>Hello { name }!</div>
}

// main.go
func(w http.ResponseWriter, r *http.Request) {
	htmx.NewResponse().
		Retarget("#hello").
		RenderTempl(r.Context(), w, Hello())
}
```

## Header names

If you need to work with HTMX headers directly, htmx-go provides constant values for all 
HTTP header field names of HTMX so you don't have to write them yourself. This mitigates the risk of writing
header names with typos.

```go
// Request headers
const (
	HeaderBoosted               = "HX-Boosted"
	HeaderCurrentURL            = "HX-Current-URL"
	HeaderHistoryRestoreRequest = "HX-History-Restore-Request"
	HeaderPrompt                = "HX-Prompt"
	HeaderRequest               = "Hx-Request"
	HeaderTarget                = "HX-Target"
	HeaderTriggerName           = "Hx-Trigger-Name"
)

// Common headers
const (
	HeaderTrigger = "HX-Trigger"
)

// Response headers
const (
	HeaderLocation           = "HX-Location"
	HeaderPushURL            = "HX-Push-Url"
	HeaderRedirect           = "HX-Redirect"
	HeaderRefresh            = "HX-Refresh"
	HeaderReplaceUrl         = "HX-Replace-Url"
	HeaderReswap             = "HX-Reswap"
	HeaderRetarget           = "HX-Retarget"
	HeaderReselect           = "HX-Reselect"
	HeaderTriggerAfterSettle = "HX-Trigger-After-Settle"
	HeaderTriggerAfterSwap   = "HX-Trigger-After-Swap"
)
```

## Compatibility

This library is compatible with the standard `net/http` library, as well as other routers like Chi
and Gorilla Mux that use the standard `http.HandlerFunc` handler type.

With the Echo web framework, try passing in `context.Request()` and
`context.Response().Writer` for requests and responses, respectively.

With the Gin web framework on the other hand, try using `context.Request` and
`context.Writer`.

## Additional resources

- [HTMX - HTTP Header Reference](https://htmx.org/reference/#headers)

## License

[MIT](./LICENSE)
