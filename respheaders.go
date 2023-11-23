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
	//
	// Valid value for Response.Reswap().
	SwapInnerHTML swapValue = "innerHTML"

	// Replace the entire target element with the response.
	//
	// Valid value for Response.Reswap().
	SwapOuterHTML swapValue = "outerHTML"

	// Insert the response before the target element.
	//
	// Valid value for Response.Reswap().
	SwapBeforeBegin swapValue = "beforebegin"

	// Insert the response before the first child of the target element.
	//
	// Valid value for Response.Reswap().
	SwapAfterBegin swapValue = "afterbegin"

	// Insert the response after the last child of the target element.
	//
	// Valid value for Response.Reswap().
	SwapBeforeEnd swapValue = "beforeend"

	// Insert the response after the target element.
	//
	// Valid value for Response.Reswap().
	SwapAfterEnd swapValue = "afterend"

	// Deletes the target element regardless of the response.
	//
	// Valid value for Response.Reswap().
	SwapDelete swapValue = "delete"

	// Does not append content from response (out of band items will still be processed).
	//
	// Valid value for Response.Reswap().
	SwapNone swapValue = "none"
)

// StatusCode sets the HTTP response header of this response.
//
// If StatusCode is not called, the default status code will be 200 OK.
func (r Response) StatusCode(statusCode int) Response {
	r.setStatusCode(statusCode)
	return r
}

// Internal method for StatusCode
func (r *Response) setStatusCode(statusCode int) {
	r.statusCode = statusCode
}

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
// Sets the same header as Response.PreventPushURL(), overwriting previous set headers.
//
// Sets the 'HX-Push-Url' header.
func (r Response) PushURL(url string) Response {
	r.headers[HeaderPushURL] = url
	return r
}

// PreventPushURL prevents the browser’s history from being updated.
//
// Sets the same header as Response.PushURL(), overwriting previous set headers.
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
// Sets the same header as Response.PreventReplaceURL(), overwriting previous set headers.
//
// Sets the 'HX-Replace-Url' header.
func (r Response) ReplaceURL(url string) Response {
	r.headers[HeaderReplaceUrl] = url
	return r
}

// PreventReplaceURL prevents the browser’s current URL from being updated.
//
// Sets the same header as Response.ReplaceURL(), overwriting previous set headers.
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
	// Just a bare interface for restricting the triggerer types
	triggerer interface {
		htmxTrigger()
	}

	// Unexported with a public constructor function for type safety reasons
	trigger string
	// Unexported with a public constructor function for type safety reasons
	triggerValue struct {
		eventName string
		value     string
	}
	// Unexported with a public constructor function for type safety reasons
	triggerKeyValue struct {
		eventName string
		value     map[string]string
	}
)

// Trigger satisfies htmx.trigger
func (t trigger) htmxTrigger() {}

// TriggerValue satisfies htmx.trigger
func (t triggerValue) htmxTrigger() {}

// TriggerKeyValue satisfies htmx.trigger
func (t triggerKeyValue) htmxTrigger() {}

// Example:
//
//	TriggerWithValue("myEvent")
//
// Output header:
//
//	HX-Trigger: myEvent
func Trigger(eventName string) trigger {
	return trigger(eventName)
}

// Example:
//
//	TriggerValue("showMessage", "Here Is A Message")
//
// Output header:
//
//	HX-Trigger: {"showMessage":"Here Is A Message"}
func TriggerValue(eventName string, value string) triggerValue {
	return triggerValue{
		eventName: eventName,
		value:     value,
	}
}

// Example:
//
//	TriggerWithValue("showMessage", map[string]string{
//	  "level": "info",
//	  "message": "Here Is A Message",
//	})
//
// Output header:
//
//	HX-Trigger: {"showMessage":{"level" : "info", "message" : "Here Is A Message"}}
func TriggerKeyValue(eventName string, value map[string]string) triggerKeyValue {
	return triggerKeyValue{
		eventName: eventName,
		value:     value,
	}
}

// **( TODO function, not working yet )**
//
// AddTrigger adds a trigger for events that trigger as soon as the response is received.
//
// This can be called multiple times so you can add multiple triggers for different events.
//
// Sets the 'HX-Trigger' header.
func (r Response) AddTrigger(trigger triggerer) Response {
	// TODO: AddTrigger
	return r
}

// **( TODO function, not working yet )**
//
// AddTriggerAfterSettle adds a trigger for events that trigger after the settling step.
//
// This can be called multiple times so you can add multiple triggers for different events.
//
// Sets the 'HX-Trigger-After-Settle' header.
func (r Response) AddTriggerAfterSettle(trigger triggerer) Response {
	// TODO: AddTriggerAfterSettle
	return r
}

// **( TODO function, not working yet )**
//
// AddTriggerAfterSwap adds a trigger for events that trigger after the swap step.
//
// This can be called multiple times so you can add multiple triggers for different events.
//
// Sets the 'HX-Trigger-After-Swap' header.
func (r Response) AddTriggerAfterSwap(trigger triggerer) Response {
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
