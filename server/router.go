package server

import (
	"context"
	"github.com/xpfo-go/sql2api/util"
	"net/http"
	"sort"
	"sync"
)

var defaultServeMux *ServeMux

func init() {
	defaultServeMux = GetRouter()
}

func GetRouter() *ServeMux {
	if defaultServeMux != nil {
		return defaultServeMux
	}

	defaultServeMux = &ServeMux{
		m: make(map[string]muxEntry),
	}

	return defaultServeMux
}

type ServeMux struct {
	mu sync.RWMutex
	m  map[string]muxEntry
}

type muxEntry struct {
	pattern string
	method  string
	handler http.HandlerFunc
}

func (s *ServeMux) RegisterFunc(method, pattern string, handler http.HandlerFunc) {
	// TODO: format method

	pattern = s.formatPattern(pattern)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[pattern] = muxEntry{
		method:  method,
		handler: handler,
		pattern: pattern,
	}
}

func (s *ServeMux) GetApiList() []string {
	res := make([]string, 0, len(s.m))
	s.mu.Lock()
	defer s.mu.Unlock()

	for pattern := range s.m {
		res = append(res, pattern)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})
	return res
}

func (s *ServeMux) DeleteRouter(pattern string) {
	pattern = s.formatPattern(pattern)

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, pattern)
}

func (s *ServeMux) formatPattern(pattern string) string {
	if pattern[0] != '/' {
		pattern = "/" + pattern
	}
	return pattern
}

func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	handler, ok := s.m[r.URL.EscapedPath()]
	s.mu.Unlock()

	if !ok || handler.method != r.Method {
		http.NotFound(w, r)
		return
	}

	requestID := r.Header.Get(util.RequestIDHeaderKey)
	if requestID == "" || len(requestID) != 32 {
		requestID = util.GenUUID4()
	}

	ctx := context.WithValue(r.Context(), util.RequestIDHeaderKey, requestID)
	handler.handler(w, r.WithContext(ctx))
}
