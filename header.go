// Package htmx provides utilities to build HTMX-powered web applications.
package htmx

// HTTP request headers
const (
	// Request header that is "true" if the request was made from an element using 'hx-boost'.
	HeaderBoosted = "HX-Boosted"
	// Request header for the current URL of the browser.
	HeaderCurrentURL = "HX-Current-URL"
	// Request header that is “true” if the request is for history restoration after a miss in the local history cache.
	HeaderHistoryRestoreRequest = "HX-History-Restore-Request"
	// Request header for the user response to an hx-prompt.
	HeaderPrompt = "HX-Prompt"
	// Request header that is always “true” for HTMX requests.
	HeaderRequest = "Hx-Request"
	// Request header of the id of the target element if it exists.
	HeaderTarget = "HX-Target"
	// Request header of the name of the triggered element if it exists.
	HeaderTriggerName = "Hx-Trigger-Name"
)

// Common HTTP headers
const (
	// As a request header: The ID of the triggered element if it exists.
	//
	// As a response header: Allows you to trigger client-side events.
	HeaderTrigger = "HX-Trigger"
)

// HTTP response headers
const (
	// Response header that allows you to do a client-side redirect that does not do a full page reload.
	HeaderLocation = "HX-Location"
	// Response header that pushes a new url into the history stack.
	HeaderPushURL = "HX-Push-Url"
	// Response header that can be used to do a client-side redirect to a new location.
	HeaderRedirect = "HX-Redirect"
	// Response header that if set to “true” the client-side will do a full refresh of the page.
	HeaderRefresh = "HX-Refresh"
	// Response header that replaces the current URL in the location bar.
	HeaderReplaceUrl = "HX-Replace-Url"
	// Response header that allows you to specify how the response will be swapped.
	HeaderReswap = "HX-Reswap"
	// Response header that uses a CSS selector that updates the target of the content update to a
	// different element on the page.
	HeaderRetarget = "HX-Retarget"
	// Response header that uses a CSS selector that allows you to choose which
	// part of the response is used to be swapped in. Overrides an existing hx-select
	// on the triggering element.
	HeaderReselect = "HX-Reselect"
	// Response header that allows you to trigger client-side events after the settle step.
	HeaderTriggerAfterSettle = "HX-Trigger-After-Settle"
	// Response header that allows you to trigger client-side events after the swap step.
	HeaderTriggerAfterSwap = "HX-Trigger-After-Swap"
)

// 286 Stop Polling
//
// HTTP status code that tells HTMX to stop polling from a server response.
//
// For more info, see https://htmx.org/docs/#load_polling
const StatusStopPolling int = 286
