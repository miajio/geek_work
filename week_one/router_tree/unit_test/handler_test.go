package unit_test

import (
	"fmt"
	rt "router_tree"
	"strings"
	"testing"
)

func TestHttpHandler(t *testing.T) {
	hs := rt.New()
	hs.Handle("GET", "/index", func(ctx *rt.Context) {
		ctx.Writer.Write([]byte(rt.SuccessPage))
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
