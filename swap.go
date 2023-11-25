package htmx

import (
	"strings"
	"time"
)

// SwapStrategy is an 'hx-swap' value that determines the swapping strategy
// of [htmx.Response.Reswap] and [LocationContext].
//
// SwapStrategy methods add modifiers to change the behavior of the swap.
type SwapStrategy string

const (
	// Replace the inner html of the target element.
	//
	// Valid value for [Response.Reswap].
	SwapInnerHTML SwapStrategy = "innerHTML"

	// Replace the entire target element with the response.
	//
	// Valid value for [Response.Reswap].
	SwapOuterHTML SwapStrategy = "outerHTML"

	// Insert the response before the target element.
	//
	// Valid value for [Response.Reswap].
	SwapBeforeBegin SwapStrategy = "beforebegin"

	// Insert the response before the first child of the target element.
	//
	// Valid value for [Response.Reswap].
	SwapAfterBegin SwapStrategy = "afterbegin"

	// Insert the response after the last child of the target element.
	//
	// Valid value for [Response.Reswap].
	SwapBeforeEnd SwapStrategy = "beforeend"

	// Insert the response after the target element.
	//
	// Valid value for [Response.Reswap].
	SwapAfterEnd SwapStrategy = "afterend"

	// Deletes the target element regardless of the response.
	//
	// Valid value for [Response.Reswap].
	SwapDelete SwapStrategy = "delete"

	// Does not append content from response (out of band items will still be processed).
	//
	// Valid value for [Response.Reswap].
	SwapNone SwapStrategy = "none"

	// Uses the default swap style (default in HTMX is [SwapInnerHTML]).
	// This value is useful for adding modifiers to the [SwapStrategy]
	// through methods
	// without changing the default swap style.
	//
	//
	// Valid value for [Response.Reswap].
	SwapDefault SwapStrategy = ""
)

func (s SwapStrategy) swapString() string {
	return string(s)
}

// join joins any amount of strings together with a space in between.
func join(elems ...string) string {
	// TrimSpace is needed because strings.Join inserts
	// the separator string (spaces) at the start of the string
	return strings.TrimSpace(strings.Join(elems, " "))
}

func (s SwapStrategy) cutPrefix(prefix string) string {
	words := strings.Split(s.swapString(), " ")
	filteredWords := make([]string, len(words))

	for _, word := range words {
		if strings.HasPrefix(word, prefix+":") {
			continue
		}
		filteredWords = append(filteredWords, word)
	}

	return join(filteredWords...)
}

func (s SwapStrategy) boolModifier(prefix string, b bool) SwapStrategy {
	v := s.cutPrefix(prefix)

	v = join(v, prefix+":")
	if b {
		v = v + "true"
	} else {
		v = v + "false"
	}

	return SwapStrategy(v)
}

func (s SwapStrategy) timeModifier(prefix string, duration time.Duration) SwapStrategy {
	v := s.cutPrefix(prefix)
	v = join(v, prefix+":"+duration.String())
	return SwapStrategy(v)
}

// Transition makes the swap use the new View Transitions API.
//
// Adds the 'transition:<true | false>' modifier.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (s SwapStrategy) Transition(shouldTransition bool) SwapStrategy {
	return s.boolModifier("transition", shouldTransition)
}

// IgnoreTitle prevents HTMX from updating the title of the page
// if there is a '<title>' tag in the response content.
//
// By default, HTMX updates the title.
//
// Adds the 'ignoreTitle:<true | false>' modifier.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (s SwapStrategy) IgnoreTitle(shouldIgnore bool) SwapStrategy {
	return s.boolModifier("ignoreTitle", shouldIgnore)
}

// FocusScroll enables focus scroll to automatically scroll
// to the focused element after a request.
//
// Adds the 'focusScroll:<true | false>' modifier.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (s SwapStrategy) FocusScroll(shouldFocus bool) SwapStrategy {
	return s.boolModifier("focusScroll", shouldFocus)
}

// After modifies the amount of time that HTMX will wait
// after receiving a response to swap the content.
//
// Adds the 'swap:<duration>' modifier.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (s SwapStrategy) After(duration time.Duration) SwapStrategy {
	return s.timeModifier("swap", duration)
}

// SettleAfter modifies the amount of time that HTMX will wait
// after the swap before executing the settle logic.
//
// Adds the 'settle:<duration>' modifier.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (s SwapStrategy) SettleAfter(duration time.Duration) SwapStrategy {
	return s.timeModifier("settle", duration)
}

type (
	// Direction is a value for the [htmx.SwapStrategy] 'scroll' and 'show' modifier methods.
	//
	// Possible values are [htmx.Top] and [htmx.Bottom].
	Direction interface {
		dirString() string
	}

	direction string
)

func (d direction) dirString() string {
	return string(d)
}

// Direction values for the [htmx.SwapStrategy] 'scroll' and 'show' modifier methods.
const (
	// Direction value for the 'scroll' and 'show' swap modifier methods.
	Top direction = "top"
	// Direction value for the 'scroll' and 'show' swap modifier methods.
	Bottom direction = "bottom"
)

// Scroll changes the scrolling behavior based on the given direction.
//
// Scroll([htmx.Top]) scrolls to the top of the swapped-in element.
//
// Scroll([htmx.Bottom]) scrolls to the bottom of the swapped-in element.
//
// Adds the 'scroll:<direction ("top" | "bottom")>' modifier.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (s SwapStrategy) Scroll(direction Direction) SwapStrategy {
	v := s.cutPrefix("scroll")
	mod := "scroll:" + direction.dirString()
	return SwapStrategy(join(v, mod))
}

// ScrollOn changes the scrolling behavior based on the given direction and CSS selector.
//
// ScrollOn(cssSelector, [htmx.Top]) scrolls to the top of the element found by the selector.
//
// ScrollOn(cssSelector, [htmx.Bottom]) scrolls to the bottom of the element found by the selector.
//
// Adds the 'scroll:<cssSelector>:<direction ("top" | "bottom")>' modifier.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (s SwapStrategy) ScrollOn(cssSelector string, direction Direction) SwapStrategy {
	v := s.cutPrefix("scroll")
	mod := "scroll:" + cssSelector + ":" + direction.dirString()
	return SwapStrategy(join(v, mod))
}

// ScrollWindow changes the scrolling behavior based on the given direction.
//
// ScrollWindow([htmx.Top]) scrolls to the very top of the window.
//
// ScrollWindow([htmx.Bottom]) scrolls to the very bottom of the window.
//
// Adds the 'scroll:window:<direction ("top" | "bottom")>' modifier.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (s SwapStrategy) ScrollWindow(direction Direction) SwapStrategy {
	v := s.cutPrefix("scroll")
	mod := "scroll:window:" + direction.dirString()
	return SwapStrategy(join(v, mod))
}

// Show changes the show behavior based on the given direction.
//
// Show([htmx.Top]) shows the top of the swapped-in element.
//
// Show([htmx.Bottom]) shows the bottom of the swapped-in element.
//
// Adds the 'show:<direction ("top" | "bottom")>' modifier.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (s SwapStrategy) Show(direction Direction) SwapStrategy {
	v := s.cutPrefix("show")
	mod := "show:" + direction.dirString()
	return SwapStrategy(join(v, mod))
}

// ShowOn changes the show behavior based on the given direction and CSS selector.
//
// ShowOn(cssSelector, [htmx.Top]) shows the top of the element found by the selector.
//
// ShowOn(cssSelector, [htmx.Bottom]) shows the bottom of the element found by the selector.
//
// Adds the 'show:<cssSelector>:<direction ("top" | "bottom")>' modifier.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (s SwapStrategy) ShowOn(cssSelector string, direction Direction) SwapStrategy {
	v := s.cutPrefix("show")
	mod := "show:" + cssSelector + ":" + direction.dirString()
	return SwapStrategy(join(v, mod))
}

// ShowWindow changes the show behavior based on the given direction.
//
// ScrollWindow([htmx.Top]) shows the very top of the window.
//
// ScrollWindow([htmx.Bottom]) shows the very bottom of the window.
//
// Adds the 'show:window:<direction ("top" | "bottom")>' modifier.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (s SwapStrategy) ShowWindow(direction Direction) SwapStrategy {
	v := s.cutPrefix("show")
	mod := "show:window:" + direction.dirString()
	return SwapStrategy(join(v, mod))
}

// ShowNone disables 'show'.
//
// Adds the 'show:none' modifier.
//
// For more info, see https://htmx.org/attributes/hx-swap/
func (s SwapStrategy) ShowNone() SwapStrategy {
	v := s.cutPrefix("show")
	mod := "show:none"
	return SwapStrategy(join(v, mod))
}
