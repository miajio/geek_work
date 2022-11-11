package routertree

import "net/http"

// Context 全局上下文
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
}
