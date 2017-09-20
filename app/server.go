package app

import (
	"fmt"
	"log"
	"net/http"
	"requests-counter/config"
	"requests-counter/handlers"
)

type (
	Server struct {
		dispatcher *Dispatcher
		saver      *Saver
		errorChan  chan error
	}
)

func NewServer(dispatcher *Dispatcher) *Server {
	return &Server{
		dispatcher: dispatcher,
		errorChan:  make(chan error, 10),
	}
}

func (s *Server) Start() {
	listenAddr := config.Configuration.Server.Host + config.Configuration.Server.Port
	println(fmt.Sprintf("Server started...\nListing to http://%s", listenAddr))

	s.configureRoutes()

	go func() {
		s.errorChan <- http.ListenAndServe(listenAddr, nil)
	}()
	go s.monitorErrors()
}

func (s *Server) configureRoutes() {
	routes := []string{
		"/metrics",
		"/stats",
		"/health",
		"/status",
	}

	for _, route := range routes {
		http.Handle(route, s.getHandler(route))
	}
}

func (s *Server) getHandler(route string) http.Handler {
	var handler http.Handler

	switch route {
	case "/metrics", "/stats":
		handler = handlers.NewMetricsHandler(s.dispatcher.requestChan, s.dispatcher.responseChan)

	case "/health", "/status":
		handler = handlers.NewHealthHandler()
	}

	return handler
}

func (s *Server) monitorErrors() {
	for {
		select {
		case err := <-s.errorChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (s *Server) Stop() {
	s.dispatcher.Stop()
}
