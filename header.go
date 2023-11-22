package htmx

// Request headers
const (
	// Request header that indicates that the request is via an element using hx-boost.
	HeaderBoosted = "HX-Boosted"
	// Request header for the current URL of the browser.
	HeaderCurrentURL = "HX-Current-URL"
	// Request header that is “true” if the request is for history restoration after a miss in the local history cache.
	HeaderHistoryRestoreRequest = "HX-History-Restore-Request"
	// Request header for the user response to an hx-prompt.
	HeaderPrompt = "HX-Prompt"
	// Request header that is always “true” for HTMX requests.
	HeaderRequest     = "Hx-Request"
	HeaderTarget      = "HX-Target"
	HeaderTriggerName = "Hx-Trigger-Name"
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
