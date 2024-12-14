package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Webserver struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebserverPort string
}

func NewWebServer(webserverPort string) *Webserver {
	return &Webserver{
		WebserverPort: webserverPort,
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
	}
}

func (s *Webserver) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *Webserver) Start() {
	s.Router.Use(middleware.Logger)
	s.Router = chi.NewRouter()
	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}

	if err := http.ListenAndServe(s.WebserverPort, s.Router); err != nil {
		panic(err.Error())
	}

}
