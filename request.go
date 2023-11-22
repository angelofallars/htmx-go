package htmx

import "net/http"

// IsHTMXRequest returns true if the given request
// was sent by HTMX.
//
// Checks if 'HX-Request' is 'true'
func IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get(HeaderRequest) == "true"
}

// IsBoosted returns true if the given request
// was made via an element using 'hx-boost'.
//
// Checks if 'HX-Boosted' is 'true'
func IsBoosted(r *http.Request) bool {
	return r.Header.Get(HeaderBoosted) == "true"
}

// IsHistoryRestoreRequest returns true if the given request
// is for history restoration after a miss in the local history cache
//
// Checks if 'HX-History-Restore-Request' is 'true'
func IsHistoryRestoreRequest(r *http.Request) bool {
	return r.Header.Get(HeaderHistoryRestoreRequest) == "true"
}

// GetCurrentURL returns the current URL of the request.
//
// If header 'HX-Current-URL' is not found, returns false on the returned bool.
func GetCurrentURL(r *http.Request) (string, bool) {
	if _, ok := r.Header[HeaderCurrentURL]; !ok {
		return "", false
	}
	return r.Header.Get(HeaderCurrentURL), true
}

// GetPrompt returns the user response to an hx-prompt from a given request.
//
// If header 'HX-Prompt' is not found, returns false on the returned bool.
func GetPrompt(r *http.Request) (string, bool) {
	if _, ok := r.Header[HeaderPrompt]; !ok {
		return "", false
	}
	return r.Header.Get(HeaderPrompt), true
}

// GetTarget returns the CSS selector of the target element if it exists from a given request.
// If header 'HX-Target' is not found, returns false on the returned bool.
func GetTarget(r *http.Request) (string, bool) {
	if _, ok := r.Header[HeaderPrompt]; !ok {
		return "", false
	}
	return r.Header.Get(HeaderPrompt), true
}

// GetTriggerName returns the 'name' of the triggered element if it exists from a given request.
//
// If header 'HX-Trigger-Name' is not found, returns false on the returned bool.
func GetTriggerName(r *http.Request) (string, bool) {
	if _, ok := r.Header[HeaderPrompt]; !ok {
		return "", false
	}
	return r.Header.Get(HeaderPrompt), true
}

// GetTrigger returns the CSS selector of the triggered element if it exists from a given request.
//
// If header 'HX-Trigger' is not found, returns false on the returned bool.
func GetTrigger(r *http.Request) (string, bool) {
	if _, ok := r.Header[HeaderPrompt]; !ok {
		return "", false
	}
	return r.Header.Get(HeaderPrompt), true
}
