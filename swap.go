package htmx

type (
	// Interface to define valid 'hx-swap' values.
	swapper interface {
		swap() string
	}
	// `hx-swap` value to determine the strategy of swapping.
	swapStrategy string
)

func (s swapStrategy) swap() string {
	return string(s)
}

const (
	// Replace the inner html of the target element.
	//
	// Valid value for Response.Reswap().
	SwapInnerHTML swapStrategy = "innerHTML"

	// Replace the entire target element with the response.
	//
	// Valid value for Response.Reswap().
	SwapOuterHTML swapStrategy = "outerHTML"

	// Insert the response before the target element.
	//
	// Valid value for Response.Reswap().
	SwapBeforeBegin swapStrategy = "beforebegin"

	// Insert the response before the first child of the target element.
	//
	// Valid value for Response.Reswap().
	SwapAfterBegin swapStrategy = "afterbegin"

	// Insert the response after the last child of the target element.
	//
	// Valid value for Response.Reswap().
	SwapBeforeEnd swapStrategy = "beforeend"

	// Insert the response after the target element.
	//
	// Valid value for Response.Reswap().
	SwapAfterEnd swapStrategy = "afterend"

	// Deletes the target element regardless of the response.
	//
	// Valid value for Response.Reswap().
	SwapDelete swapStrategy = "delete"

	// Does not append content from response (out of band items will still be processed).
	//
	// Valid value for Response.Reswap().
	SwapNone swapStrategy = "none"
)
