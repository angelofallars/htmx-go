package htmx

import (
	"encoding/json"
	"strings"
)

const (
	trueString  = "true"
	falseString = "false"
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
// you can use [htmx.Response.LocationWithContext].
//
// Sets the 'HX-Location' header.
//
// For more info, see https://htmx.org/headers/hx-location/
func (r Response) Location(path string) Response {
	r.headers[HeaderLocation] = path
	return r
}

// A context object that is used by [htmx.Response.LocationWithContext]
// to finely determine the behavior of client-side redirection.
//
// In the browser, these are the parameters that will be used by 'htmx.ajax()'.
//
// For more info, see https://htmx.org/headers/hx-location/
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
	Swap SwapStrategy
	// Values to submit with the request.
	Values map[string]string
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
	Swap    string            `json:"swap,omitempty"`
	Values  map[string]string `json:"values,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Select  string            `json:"select,omitempty"`
}

// LocationWithContext allows you to do a client-side redirect that does not do a full page reload,
// redirecting to a specific target on the page with the given context.
//
// For simple redirects, you can just use [htmx.Response.Location].
//
// Sets the 'HX-Location' header.
//
// For more info, see https://htmx.org/headers/hx-location/
func (r Response) LocationWithContext(path string, ctx LocationContext) Response {
	// Replace the error at the start because the last errors shouldn't really matter
	r.locationWithContextErr = make([]error, 0)

	c := locationContext{
		Path:    path,
		Source:  ctx.Source,
		Event:   ctx.Event,
		Handler: ctx.Handler,
		Target:  ctx.Target,
		Swap:    ctx.Swap.swapString(),
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
// Sets the same header as [htmx.Response.PreventPushURL], overwriting previous set headers.
//
// Sets the 'HX-Push-Url' header.
//
// For more info, see https://htmx.org/headers/hx-push-url/
func (r Response) PushURL(url string) Response {
	r.headers[HeaderPushURL] = url
	return r
}

// PreventPushURL prevents the browser’s history from being updated.
//
// Sets the same header as [htmx.Response.PushURL], overwriting previous set headers.
//
// Sets the 'HX-Push-Url' header.
//
// For more info, see https://htmx.org/headers/hx-push-url/
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
// Sets the same header as [htmx.Response.PreventReplaceURL], overwriting previous set headers.
//
// Sets the 'HX-Replace-Url' header.
//
// For more info, see https://htmx.org/headers/hx-replace-url/
func (r Response) ReplaceURL(url string) Response {
	r.headers[HeaderReplaceUrl] = url
	return r
}

// PreventReplaceURL prevents the browser’s current URL from being updated.
//
// Sets the same header as [htmx.Response.ReplaceURL], overwriting previous set headers.
//
// Sets the 'HX-Replace-Url' header to 'false'.
//
// For more info, see https://htmx.org/headers/hx-replace-url/
func (r Response) PreventReplaceURL() Response {
	r.headers[HeaderReplaceUrl] = falseString
	return r
}

// Reswap allows you to specify how the response will be swapped.
//
// Sets the 'HX-Reswap' header.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (r Response) Reswap(s SwapStrategy) Response {
	r.headers[HeaderReswap] = s.swapString()
	return r
}

// Retarget accepts a CSS selector that updates the target of the content update to a different element on the page. Overrides an existing 'hx-select' on the triggering element.
//
// Sets the 'HX-Retarget' header.
//
// For more info, see https://htmx.org/attributes/hx-target/
func (r Response) Retarget(cssSelector string) Response {
	r.headers[HeaderRetarget] = cssSelector
	return r
}

// Reselect accepts a CSS selector that allows you to choose which part of the response is used to be swapped in.
// Overrides an existing hx-select on the triggering element.
//
// Sets the 'HX-Reselect' header.
//
// For more info, see https://htmx.org/attributes/hx-select/
func (r Response) Reselect(cssSelector string) Response {
	r.headers[HeaderReselect] = cssSelector
	return r
}

type (
	// EventTrigger gives an HTMX response directives to
	// triggers events on the client side.
	EventTrigger interface {
		htmxTrigger()
	}

	// Unexported with a public constructor function for type safety reasons
	triggerPlain string
	// Unexported with a public constructor function for type safety reasons
	triggerDetail struct {
		eventName string
		value     string
	}
	// Unexported with a public constructor function for type safety reasons
	triggerObject struct {
		eventName string
		object    any
	}
)

// trigger satisfies htmx.EventTrigger
func (t triggerPlain) htmxTrigger() {}

// triggerDetail satisfies htmx.EventTrigger
func (t triggerDetail) htmxTrigger() {}

// triggerObject satisfies htmx.EventTrigger
func (t triggerObject) htmxTrigger() {}

// Trigger returns an event trigger with no additional details.
//
// Example:
//
//	htmx.Trigger("myEvent")
//
// Output header:
//
//	HX-Trigger: myEvent
//
// For more info, see https://htmx.org/headers/hx-trigger/
func Trigger(eventName string) triggerPlain {
	return triggerPlain(eventName)
}

// TriggerDetail returns an event trigger with one detail string.
// Will be encoded as JSON.
//
// Example:
//
//	htmx.TriggerDetail("showMessage", "Here Is A Message")
//
// Output header:
//
//	HX-Trigger: {"showMessage":"Here Is A Message"}
//
// For more info, see https://htmx.org/headers/hx-trigger/
func TriggerDetail(eventName string, detailValue string) triggerDetail {
	return triggerDetail{
		eventName: eventName,
		value:     detailValue,
	}
}

// TriggerObject returns an event trigger with a given detail object that **must** be serializable to JSON.
//
// Structs with JSON tags can work, and so does `map[string]string` values which are safe to serialize.
//
// Example:
//
//	htmx.TriggerObject("showMessage", map[string]string{
//	  "level": "info",
//	  "message": "Here Is A Message",
//	})
//
// Output header:
//
//	HX-Trigger: {"showMessage":{"level" : "info", "message" : "Here Is A Message"}}
//
// For more info, see https://htmx.org/headers/hx-trigger/
func TriggerObject(eventName string, detailObject any) triggerObject {
	return triggerObject{
		eventName: eventName,
		object:    detailObject,
	}
}

// triggersToString converts a slice of triggers into a header value
// for headers like 'HX-Trigger'.
func triggersToString(triggers []EventTrigger) (string, error) {
	simpleEvents := make([]string, 0)
	detailEvents := make(map[string]any)

	for _, t := range triggers {
		switch v := t.(type) {
		case triggerPlain:
			simpleEvents = append(simpleEvents, string(v))
		case triggerObject:
			detailEvents[v.eventName] = v.object
		case triggerDetail:
			detailEvents[v.eventName] = v.value
		}
	}

	if len(detailEvents) == 0 {
		return strings.Join(simpleEvents, ", "), nil
	} else {
		for _, evt := range simpleEvents {
			detailEvents[evt] = ""
		}

		bytes, err := json.Marshal(detailEvents)
		if err != nil {
			return "", err
		}

		return string(bytes), nil
	}
}

// AddTrigger adds trigger(s) for events that trigger as soon as the response is received.
//
// This can be called multiple times so you can add as many triggers as you need.
//
// Sets the 'HX-Trigger' header.
//
// For more info, see https://htmx.org/headers/hx-trigger/
func (r Response) AddTrigger(trigger ...EventTrigger) Response {
	r.initTriggers()
	r.triggers = append(r.triggers, trigger...)
	return r
}

// AddTriggerAfterSettle adds trigger(s) for events that trigger after the settling step.
//
// This can be called multiple times so you can add as many triggers as you need.
//
// Sets the 'HX-Trigger-After-Settle' header.
//
// For more info, see https://htmx.org/headers/hx-trigger/
func (r Response) AddTriggerAfterSettle(trigger ...EventTrigger) Response {
	r.initTriggersAfterSettle()
	r.triggersAfterSettle = append(r.triggersAfterSettle, trigger...)
	return r
}

// AddTriggerAfterSwap adds trigger(s) for events that trigger after the swap step.
//
// This can be called multiple times so you can add as many triggers as you need.
//
// Sets the 'HX-Trigger-After-Swap' header.
//
// For more info, see https://htmx.org/headers/hx-trigger/
func (r Response) AddTriggerAfterSwap(trigger ...EventTrigger) Response {
	r.initTriggersAfterSwap()
	r.triggersAfterSwap = append(r.triggersAfterSwap, trigger...)
	return r
}

// Lazily init the triggers slice because not all responses
// use triggers
func (r *Response) initTriggers() {
	if r.triggers == nil {
		r.triggers = make([]EventTrigger, 0)
	}
}

// Lazily init the triggersAfterSettle slice because not all responses
// use triggers
func (r *Response) initTriggersAfterSettle() {
	if r.triggersAfterSettle == nil {
		r.triggersAfterSettle = make([]EventTrigger, 0)
	}
}

// Lazily init the triggersAfterSwap slice because not all responses
// use triggers
func (r *Response) initTriggersAfterSwap() {
	if r.triggersAfterSwap == nil {
		r.triggersAfterSwap = make([]EventTrigger, 0)
	}
}
