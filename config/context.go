package config

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync"
)

type contextStruct struct {
	mutex          sync.Mutex
	reverseProxies map[ReverseProxy]*httputil.ReverseProxy
}

func newContext() *contextStruct {
	return &contextStruct{
		mutex:          sync.Mutex{},
		reverseProxies: make(map[ReverseProxy]*httputil.ReverseProxy),
	}
}

func RenewContext() {
	context = newContext()
}

func (context *contextStruct) getProxy(key ReverseProxy) *httputil.ReverseProxy {
	context.mutex.Lock()
	defer context.mutex.Unlock()
	val := context.reverseProxies[key]
	if val == nil {
		u, _ := url.Parse(fmt.Sprintf("http://%s:%d", key.Host, key.Port))
		val = httputil.NewSingleHostReverseProxy(u)
		context.reverseProxies[key] = val
	}
	return val
}

var context = newContext()
