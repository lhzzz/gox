package pprof

import (
	"net/http"
	"net/http/pprof"
)

// NewHandler new a pprof handler.
func NewPprofHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return mux
}

func Serve(addr string) {
	go func() {
		handler := NewPprofHandler()
		http.ListenAndServe(addr, handler)
	}()
}
