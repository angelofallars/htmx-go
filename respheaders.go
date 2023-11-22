package htmx

import (
	"encoding/json"
	"net/http"
)

const trueString = "true"
const falseString = "false"

type (
	// Interface to define valid 'hx-swap' values.
	swapper interface {
		swap() string
	}
	// Concrete 'hx-swap' values.
	swapValue string
)

func (s swapValue) swap() string {
	return string(s)
}

const (
	// Replace the inner html of the target element.
	SwapInnerHTML swapValue = "innerHTML"

	// Replace the entire target element with the response.
	SwapOuterHTML swapValue = "outerHTML"

	// Insert the response before the target element.
	SwapBeforeBegin swapValue = "beforebegin"

	// Insert the response before the first child of the target element.
	SwapAfterBegin swapValue = "afterbegin"

	// Insert the response after the last child of the target element.
	SwapBeforeEnd swapValue = "beforeend"

	// Insert the response after the target element.
	SwapAfterEnd swapValue = "afterend"

	// Deletes the target element regardless of the response.
	SwapDelete swapValue = "delete"

	// Does not append content from response (out of band items will still be processed).
	SwapNone swapValue = "none"
)

// Location allows you to do a client-side redirect that does not do a full page reload.
//
// If you want to redirect to a specific target on the page rather than the default of document.body,
// you can use response.LocationWithContext().
//
// Sets the 'HX-Location' header.
func (r Response) Location(path string) Response {
	r.headers[HeaderLocation] = path
	return r
}

// A context object that is used by htmx.response.LocationWithContext()
// to finely determine the behavior of client-side redirection.
//
// In the browser, these are the parameters that will be used by 'htmx.ajax()'.
type LocationContext struct {
	// The source element of the request.
	Source string
	// An event that “triggered” the request.
	Event string
	// A JavaScript callback that will handle the response HTML.
	Handler string
	// The target to swap the response into.
	Target string
	// How the response will be swapped in relative to the target.
	Swap swapper
	// Values to submit with the request.
	Values []string
	// Headers to submit with the request
	Headers map[string]string
	// Allows you to select the content you want swapped from a response.
	Select string
}

// Internal version of locationContext that just contains the "path" field
// and JSON tags to serialize it to JSON.
type locationContext struct {
	// Path to redirect to.
	Path    string            `json:"path"`
	Source  string            `json:"source,omitempty"`
	Event   string            `json:"event,omitempty"`
	Handler string            `json:"handler,omitempty"`
	Target  string            `json:"target,omitempty"`
	Swap    swapper           `json:"swap,omitempty"`
	Values  []string          `json:"values,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Select  string            `json:"select,omitempty"`
}

// LocationWithContext allows you to do a client-side redirect that does not do a full page reload,
// redirecting to a specific target on the page with the given context.
//
// For simple redirects, you can just use response.Location().
//
// Sets the 'HX-Location' header.
func (r Response) LocationWithContext(path string, ctx LocationContext) Response {
	// Replace the error at the start because the last errors shouldn't really matter
	r.locationWithContextErr = make([]error, 0)

	c := locationContext{
		Path:    path,
		Source:  ctx.Source,
		Event:   ctx.Event,
		Handler: ctx.Handler,
		Target:  ctx.Target,
		Swap:    ctx.Swap,
		Values:  ctx.Values,
		Headers: ctx.Headers,
		Select:  ctx.Select,
	}

	bytes, err := json.Marshal(c)
	if err != nil {
		r.locationWithContextErr = append(r.locationWithContextErr, err)
		return r
	}

	r.headers[HeaderLocation] = string(bytes)

	return r
}

// PushURL pushes a new URL into the browser location history.
//
// Sets the 'HX-Push-Url' header.
func (r Response) PushURL(url string) Response {
	r.headers[HeaderPushURL] = url
	return r
}

// PreventPushURL prevents the browser’s history from being updated.
//
// Sets the 'HX-Push-Url' header.
func (r Response) PreventPushURL() Response {
	r.headers[HeaderPushURL] = falseString
	return r
}

// Redirect does a client-side redirect to a new location.
//
// Sets the 'HX-Redirect' header.
func (r Response) Redirect(path string) Response {
	r.headers[HeaderRedirect] = path
	return r
}

// If set to true, Refresh makes the client-side do a full refresh of the page.
//
// Sets the 'HX-Refresh' header.
func (r Response) Refresh(shouldRefresh bool) Response {
	if shouldRefresh {
		r.headers[HeaderRefresh] = trueString
	} else {
		r.headers[HeaderRefresh] = falseString
	}
	return r
}

// ReplaceURL replaces the current URL in the browser location history.
//
// Sets the 'HX-Replace-Url' header.
func (r Response) ReplaceURL(url string) Response {
	r.headers[HeaderReplaceUrl] = url
	return r
}

// PreventReplaceURL prevents the browser’s current URL from being updated.
//
// Sets the 'HX-Replace-Url' header to 'false'.
func (r Response) PreventReplaceURL() Response {
	r.headers[HeaderReplaceUrl] = falseString
	return r
}

// Reswap allows you to specify how the response will be swapped. Accepts 'htmx.Swap*' values from this library.
//
// Sets the 'HX-Reswap' header.
func (r Response) Reswap(s swapper) Response {
	r.headers[HeaderReswap] = s.swap()
	return r
}

// Retarget accepts a CSS selector that updates the target of the content update to a different element on the page.
//
// Sets the 'HX-Retarget' header.
func (r Response) Retarget(selector string) Response {
	r.headers[HeaderRetarget] = selector
	return r
}

// Reselect accepts a CSS selector that allows you to choose which part of the response is used to be swapped in.
// Overrides an existing hx-select on the triggering element.
//
// Sets the 'HX-Reselect' header.
func (r Response) Reselect(selector string) Response {
	r.headers[HeaderReselect] = selector
	return r
}

type (
	// Just a bare interface for restricting the trigger types
	trigger interface {
		htmxTrigger()
	}

	// Example output: HX-Trigger: myEvent
	Trigger string

	// Example output: HX-Trigger: {"showMessage":"Here Is A Message"}
	//
	// Unexported with a public constructor function for type safety reasons
	triggerWithValue struct {
		key   string
		value string
	}
	// Example output: HX-Trigger: {"showMessage":{"level" : "info", "message" : "Here Is A Message"}}
	//
	// Unexported with a public constructor function for type safety reasons
	triggerKeyValue struct {
		key   string
		value map[string]string
	}
)

// Trigger satisfies htmx.trigger
func (t Trigger) htmxTrigger() {}

// TriggerWithValue satisfies htmx.trigger
func (t triggerWithValue) htmxTrigger() {}

// TriggerKeyValue satisfies htmx.trigger
func (t triggerKeyValue) htmxTrigger() {}

// Example output: HX-Trigger: {"showMessage":"Here Is A Message"}
func TriggerWithValue(key string, value string) triggerWithValue {
	return triggerWithValue{
		key:   key,
		value: value,
	}
}

// Example output: HX-Trigger: {"showMessage":{"level" : "info", "message" : "Here Is A Message"}}
func TriggerKeyValue(key string, value map[string]string) triggerKeyValue {
	return triggerKeyValue{
		key:   key,
		value: value,
	}
}

// AddTrigger adds a trigger for events as soon as the response is received.
func (r Response) AddTrigger(trigger trigger) Response {
	// TODO: AddTrigger
	return r
}

// AddTriggerAfterSettle adds a trigger for events after the settling step.
func (r Response) AddTriggerAfterSettle(trigger trigger) Response {
	// TODO: AddTriggerAfterSettle
	return r
}

// AddTriggerAfterSwap adds a trigger for events after the swap step.
func (r Response) AddTriggerAfterSwap(trigger trigger) Response {
	// TODO: AddTriggerAfterSwap
	return r
}

func testResponses(w http.ResponseWriter, r *http.Request) {
	NewResponse().
		Refresh(true).
		ReplaceURL("/items").
		Write(w)

	NewResponse().AddTrigger(TriggerWithValue("event", "eee"))

	NewResponse().Retarget("#errors")

	NewResponse().Reswap(SwapBeforeEnd)
}
