package routertree

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter_Handler(t *testing.T) {
	var mockHandler HandlerFunc = func(ctx *Context) {}
	testParams := []struct {
		method      string
		path        string
		handlerFunc HandlerFunc
	}{
		{
			method:      http.MethodPost,
			path:        "",
			handlerFunc: mockHandler,
		},
		{
			method:      http.MethodPost,
			path:        "/user",
			handlerFunc: mockHandler,
		},
		{
			method:      http.MethodPost,
			path:        "/user/home",
			handlerFunc: mockHandler,
		},
		{
			method:      http.MethodPost,
			path:        "//index///",
			handlerFunc: mockHandler,
		},
		{
			method:      http.MethodGet,
			path:        "/",
			handlerFunc: mockHandler,
		},
	}

	r := NewRoute()
	for _, route := range testParams {
		r.Handle(route.method, route.path, mockHandler)
	}

	wantRoute := &Route{
		Nodes: map[string]*Node{
			http.MethodPost: {
				path:        "/",
				handlerFunc: mockHandler,
				children: map[string]*Node{
					"user": {
						path:        "user",
						handlerFunc: mockHandler,
						children: map[string]*Node{
							"home": {
								path:        "home",
								handlerFunc: mockHandler,
							},
						},
					},
					"index": {
						path:        "index",
						handlerFunc: mockHandler,
					},
				},
			},
			http.MethodGet: {
				path:        "/",
				handlerFunc: mockHandler,
				children:    map[string]*Node{},
			},
		},
	}
	msg, ok := wantRoute.Equal(r)
	assert.True(t, ok, msg)
}

func (r *Route) Equal(y *Route) (string, bool) {
	for k, v := range r.Nodes {
		det, ok := y.Nodes[k]
		if !ok {
			return fmt.Sprintf("找不到对应的http method"), false
		}
		msg, ok := v.Equal(det)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

func (n *Node) Equal(y *Node) (string, bool) {
	if n.path != y.path {
		return "路径不一致", false
	}
	if len(n.children) != len(y.children) {
		return "子节点数量不一致", false
	}

	nHandler := reflect.ValueOf(n.handlerFunc)
	yHandler := reflect.ValueOf(y.handlerFunc)

	if nHandler != yHandler {
		return "handler不一致", false
	}

	for p, c := range n.children {
		dst, ok := y.children[p]
		if !ok {
			return "子节点不存在", false
		}
		msg, ok := c.Equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

func TestHttpHandler(t *testing.T) {
	hs := New()
	hs.Handle("GET", "/index/-?[1-9]\\d*", func(ctx *Context) {
		ctx.Writer.Write([]byte(SuccessPage))
		ctx.Writer.WriteHeader(200)
		ctx.Writer.Header().Add("Content-Type", "text/html")
	})
	hs.Start(":8088")
}

func TestRouterGroup(t *testing.T) {
	url := "user/home/123"
	segs := strings.Split(url, "/")
	fmt.Println(segs)
}

func TestZZ(t *testing.T) {
	str := "18"
	matched, err := regexp.MatchString("-?[1-9]\\d*", str)
	fmt.Println(matched, err)
}
