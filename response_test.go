package htmx

import (
	"html/template"
	"net/http"
	"testing"
)

func TestWrite(t *testing.T) {
	w := newMockResponseWriter()

	err := NewResponse().
		StatusCode(StatusStopPolling).
		Location("/profiles").
		Redirect("/pull").
		PushURL("/push").
		Refresh(true).
		ReplaceURL("/water").
		Retarget("#world").
		Reselect("#hello").
		AddTrigger(Trigger("myEvent")).
		Reswap(SwapInnerHTML.ShowOn("#swappy", Top)).
		Write(w)
	if err != nil {
		t.Errorf("an error occurred writing a response: %v", err)
	}

	if w.statusCode != StatusStopPolling {
		t.Errorf("wrong error code. want=%v, got=%v", StatusStopPolling, w.statusCode)
	}

	expectedHeaders := map[string]string{
		HeaderTrigger:    "myEvent",
		HeaderLocation:   "/profiles",
		HeaderRedirect:   "/pull",
		HeaderPushURL:    "/push",
		HeaderRefresh:    "true",
		HeaderReplaceUrl: "/water",
		HeaderRetarget:   "#world",
		HeaderReselect:   "#hello",
		HeaderReswap:     "innerHTML show:#swappy:top",
	}

	for k, v := range expectedHeaders {
		got := w.header.Get(k)
		if got != v {
			t.Errorf("wrong value for header %q. got=%q, want=%q", k, got, v)
		}
	}
}

func TestRenderHTML(t *testing.T) {
	text := `hello world!`

	w := newMockResponseWriter()

	_, err := NewResponse().Location("/conversation/message").RenderHTML(w, template.HTML(text))
	if err != nil {
		t.Errorf("an error occurred writing HTML: %v", err)
	}

	if got, want := w.Header().Get(HeaderLocation), "/conversation/message"; got != want {
		t.Errorf("wrong value for header %q. got=%q, want=%q", HeaderLocation, got, want)
	}

	if string(w.body) != text {
		t.Errorf("wrong response body. got=%q, want=%q", string(w.body), text)
	}
}

func TestMustRenderHTML(t *testing.T) {
	text := `hello world!`

	w := newMockResponseWriter()

	NewResponse().MustRenderHTML(w, template.HTML(text))
}

type mockResponseWriter struct {
	body       []byte
	statusCode int
	header     http.Header
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		header: http.Header{},
	}
}

func (mrw *mockResponseWriter) Header() http.Header {
	return mrw.header
}

func (mrw *mockResponseWriter) Write(b []byte) (int, error) {
	mrw.body = append(mrw.body, b...)
	return 0, nil
}

func (mrw *mockResponseWriter) WriteHeader(statusCode int) {
	mrw.statusCode = statusCode
}
