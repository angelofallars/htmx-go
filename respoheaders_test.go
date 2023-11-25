package htmx

import (
	"testing"
	"time"
)

func TestXxx(t *testing.T) {
	testCases := []struct {
		name     string
		strategy swapStrategy
		headers  string
	}{
		{
			name:     "no modifiers",
			strategy: SwapInnerHTML,
			headers:  "innerHTML",
		},
		{
			name:     "one modifiers",
			strategy: SwapInnerHTML.WithTransition(true),
			headers:  "innerHTML transition:true",
		},
		{
			name:     "multiple modifiers",
			strategy: SwapInnerHTML.WithTransition(true).WithSwap(5 * time.Second),
			headers:  "innerHTML transition:true swap:5s",
		},
	}

	for _, tc := range testCases {
		if headers := tc.strategy.swap(); headers != tc.headers {
			t.Errorf(`got "%v", want "%v"`, headers, tc.headers)
		}
	}
}
