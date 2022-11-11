package routertree

import (
	"net/http"
	"path"
)

// IHttpServer 基于http.Handler的web http服务接口(方法声明)
type IHttpServer interface {
	Route(method, path string, httpFunc HttpFunc)
	Clean()
	Start(addr string) error
	StartTLS(addr, certFile, keyFile string) error
}

// HttpServer 基于http.Handler的web http服务
type HttpServer struct {
	routes routes
}

// HttpFunc 处理http请求的功能函数
type HttpFunc func(ctx *Context)

// route 路由
type route struct {
	group    string   // 路由分组 - 内置分组方式,直接依据path切割分组
	method   string   // http 方法
	url      string   // 路由地址
	httpFunc HttpFunc // 执行路由的方法
}

// routes 路由集
type routes []*route

// routes Hit命中 如果请求信息命中则执行route路由对应的httpFunc方法,否则返回404
func (rs routes) Hit(ctx *Context) {
	for _, r := range rs {
		if r.method == ctx.Request.Method {
			if r.url == ctx.Request.URL.Path {
				r.httpFunc(ctx)
				return
			}
		}
	}
	ctx.Writer.Write([]byte(undfindPage))
}

var _ IHttpServer = (*HttpServer)(nil)

func New() *HttpServer {
	return &HttpServer{
		routes: make(routes, 0),
	}
}

// handleDo excute route httpfunc
func (h *HttpServer) handleDo(ctx *Context) {
	h.routes.Hit(ctx)
}

// handler
func (h *HttpServer) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.handleDo(&Context{
			Writer:  w,
			Request: r,
		})
	}
}

// Route 添加路由
func (h *HttpServer) Route(method, url string, httpFunc HttpFunc) {
	h.routes = append(h.routes, &route{
		group:    path.Base(url),
		url:      url,
		method:   method,
		httpFunc: httpFunc,
	})
}

// Clean
func (h *HttpServer) Clean() {
	h.routes = make(routes, 0)
}

// Start
func (h *HttpServer) Start(addr string) error {
	return http.ListenAndServe(addr, h.handler())
}

// StartTLS
func (h *HttpServer) StartTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, h.handler())
}
