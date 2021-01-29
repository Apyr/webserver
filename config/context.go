package config

import (
	"log"
	"net/http/httputil"
	"net/url"
	"sync"
)

type contextStruct struct {
	mutex          sync.Mutex
	reverseProxies map[string]*httputil.ReverseProxy
}

func newContext() *contextStruct {
	return &contextStruct{
		mutex:          sync.Mutex{},
		reverseProxies: make(map[string]*httputil.ReverseProxy),
	}
}

func RenewContext() {
	context = newContext()
}

func (context *contextStruct) getProxy(key string) *httputil.ReverseProxy {
	context.mutex.Lock()
	defer context.mutex.Unlock()
	val := context.reverseProxies[key]
	if val == nil {
		u, err := url.Parse(key)
		if err != nil {
			log.Fatalf("Reverse proxy url parsing error: %s", err)
		}
		val = httputil.NewSingleHostReverseProxy(u)
		context.reverseProxies[key] = val
	}
	return val
}

var context = newContext()
