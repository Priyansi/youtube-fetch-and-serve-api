package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	srv *http.Server
}

func Get() *Server {
	return &Server{
		srv: &http.Server{},
	}
}

func (s *Server) SetAddrPort(addrPort string) *Server {
	s.srv.Addr = addrPort
	return s
}

func (s *Server) SetLogger(logger *log.Logger) *Server {
	s.srv.ErrorLog = logger
	return s
}

func (s *Server) SetRouter(router *httprouter.Router) *Server {
	s.srv.Handler = router
	return s
}

func (s *Server) Start() error {
	if len(s.srv.Addr) == 0 {
		return errors.New("Start: Address port is not initialized")
	}

	if s.srv.Handler == nil {
		return errors.New("Start: Handler is not initialized")
	}

	return s.srv.ListenAndServe()
}

func (s *Server) Close() error {
	return s.srv.Close()
}
