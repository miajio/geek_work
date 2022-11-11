package unit_test

import (
	rt "router_tree"
	"testing"
)

func TestHttpHandler(t *testing.T) {
	hs := rt.New()
	hs.Route("GET", "/index", func(ctx *rt.Context) {

	})
	hs.Start(":8088")
}
