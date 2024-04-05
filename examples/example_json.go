package main

import (
	"net/http"

	"github.com/libchaos/litchi/render"

	"github.com/go-chi/chi/v5"
)

type Stu struct {
	Name string `json:"name"`
}

func main() {
	r := chi.NewRouter()

	std := Stu{
		Name: "libchaos",
	}

	r.Get("/example", func(w http.ResponseWriter, r *http.Request) {
		render.OK(std).Write(w)
	})

	http.ListenAndServe(":3000", r)
}
