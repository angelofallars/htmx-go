package htmx

import (
	"testing"
	"time"
)

func TestSwapStrategy_SwapString(t *testing.T) {
	testCases := []struct {
		name         string
		swapStrategy SwapStrategy
		result       string
	}{
		{
			name:         "no modifier",
			swapStrategy: SwapInnerHTML,
			result:       "innerHTML",
		},
		{
			name:         "one modifier",
			swapStrategy: SwapInnerHTML.Transition(true),
			result:       "innerHTML transition:true",
		},
		{
			name: "many modifiers",
			swapStrategy: SwapInnerHTML.Transition(true).
				IgnoreTitle(true).
				FocusScroll(true).
				After(5*time.Second).
				SettleAfter(5*time.Second).
				Scroll(Top).
				ScrollOn("#another-div", Top).
				ScrollWindow(Top).
				Show(Top).
				ShowOn("#another-div", Top).ShowWindow(Top).
				ShowNone(),
			result: "innerHTML transition:true ignoreTitle:true focusScroll:true swap:5s settle:5s scroll:window:top show:none",
		},
	}

	for _, tc := range testCases {
		if result := tc.swapStrategy.swapString(); result != tc.result {
			t.Errorf(`got: "%v", want: "%v"`, result, tc.result)
		}
	}
}
