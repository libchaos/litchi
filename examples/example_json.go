package main

import (
	"net/http"

	"github.com/libchaos/litchi/bind"
	"github.com/libchaos/litchi/render"
	"github.com/libchaos/litchi/validate"

	"github.com/go-chi/chi/v5"
)

type Stu struct {
	Name string `json:"name"`
}

func (s *Stu) Validate() []validate.Field {
	return []validate.Field{
		validate.MinLength(&s.Name, 3),
		validate.MaxLength(&s.Name, 10),
	}
}

func main() {
	r := chi.NewRouter()

	std := Stu{
		Name: "libchaos",
	}

	r.Get("/example", func(w http.ResponseWriter, r *http.Request) {
		render.OK(std).Write(w)
	})

	r.Post("/json", func(w http.ResponseWriter, r *http.Request) {
		stu, err := bind.Body[Stu](r)
		if err != nil {
			render.BadRequest(err).Write(w)
			return
		}

		render.OK(stu).Write(w)
	})

	http.ListenAndServe(":3000", r)
}
