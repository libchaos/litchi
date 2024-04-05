package render_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/libchaos/litchi/render"

	"github.com/go-chi/chi/v5"
)

func DeprecatedHelloWorld(w http.ResponseWriter, r *http.Request) {
	render.PermanentRedirect(r, "/new").Write(w)
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	render.OK("Hello, World!").Write(w)
}

func Example_redirecting() {
	r := chi.NewRouter()
	r.Get("/", DeprecatedHelloWorld)
	r.Get("/new", HelloWorld)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(res, req)

	fmt.Println(res.Header(), res.Code)
	fmt.Println(res.Body)

	// Output:
	// map[Content-Type:[text/html; charset=utf-8] Location:[/new]] 308
	// <a href="/new">Permanent Redirect</a>.
}
