package http

import (
	"net/http"
)

type Server struct {
	Addr    string
	Scripts string
	Extra   string
}

func NewServer(addr string, path, extra string) *Server {
	return &Server{
		Addr:    addr,
		Scripts: path,
		Extra:   extra,
	}
}

func (server *Server) Serve() error {
	http.Handle("/ipxe", http.FileServer(http.Dir(server.Scripts)))

	if server.Extra != "" {
		http.Handle("/extra", http.FileServer(http.Dir(server.Extra)))
	}

	return http.ListenAndServe(server.Addr, nil)
}
