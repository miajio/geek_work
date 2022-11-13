package routertree

import (
	"net/http"
)

// IEngine 基于http.Handler的web http服务接口(方法声明)
type IEngine interface {
	Handle(method, path string, httpFunc HandlerFunc)
	Clean()
	Start(addr string) error
	StartTLS(addr, certFile, keyFile string) error
}

// Engine 基于http.Handler的web http服务
type Engine struct {
	route *Route // 依据路由集组合成的路由树
}

// HandlerFunc 处理http请求的功能函数
type HandlerFunc func(ctx *Context)

var _ IEngine = (*Engine)(nil)

// New 创建一个Engine
func New() *Engine {
	return &Engine{
		route: NewRoute(),
	}
}

// handler
func (h *Engine) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.route.Hit(&Context{
			Writer:  w,
			Request: r,
		})
	}
}

// Route 添加路由
func (h *Engine) Handle(method, url string, handlerFunc HandlerFunc) {
	h.route.Handle(method, url, handlerFunc)
}

// Clean 重置Engine
func (h *Engine) Clean() {
	h.route = NewRoute()
}

// Start
func (h *Engine) Start(addr string) error {
	return http.ListenAndServe(addr, h.handler())
}

// StartTLS
func (h *Engine) StartTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, h.handler())
}
