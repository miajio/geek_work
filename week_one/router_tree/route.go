package routertree

import (
	"regexp"
	"strings"
)

// IRoute 路有树接口声明
type IRoute interface {
	Hit(ctx *Context)                                            // Hit http.handler命中方法
	Handle(method, relativePath string, handlerFunc HandlerFunc) // Handle 添加路由节点方法
}

// Route 路由
type Route struct {
	Nodes map[string]*Node
}

var _ IRoute = (*Route)(nil)

// Node 路由节点
type Node struct {
	path        string           // 路由路径
	children    map[string]*Node // 下级子节点路由
	handlerFunc HandlerFunc      // 执行业务逻辑
}

func NewRoute() *Route {
	return &Route{
		Nodes: map[string]*Node{},
	}
}

// routes Hit命中 如果请求信息命中则执行route路由对应的httpFunc方法,否则返回404
func (rt *Route) Hit(ctx *Context) {
	nodes, ok := rt.Nodes[ctx.Request.Method]
	if !ok {
		ctx.notFundPage()
		return
	}

	url := ctx.Request.URL.Path
	if strings.ReplaceAll(url, " ", "") == "" {
		if nodes.handlerFunc == nil {
			ctx.notFundPage()
			return
		}
		nodes.handlerFunc(ctx)
		return
	}

	segs := strings.Split(url[1:], "/")
	for _, seg := range segs {
		if seg == "" || strings.ReplaceAll(seg, " ", "") == "" {
			continue
		}
		node, ok := nodes.children[seg]
		if !ok {
			isNotFundPage := true
			for k, v := range nodes.children {
				is, err := regexp.MatchString(k, seg)
				if is && err == nil {
					node = v
					isNotFundPage = false
					break
				}
			}

			if isNotFundPage {
				tpNode, ok := nodes.children["*"]
				if ok {
					node = tpNode
					continue
				}

				ctx.notFundPage()
				return
			}
		}
		nodes = node
	}
	if nodes == nil || nodes.handlerFunc == nil {
		ctx.notFundPage()
		return
	}
	nodes.handlerFunc(ctx)
}

func (ctx *Context) notFundPage() {
	ctx.Writer.Write([]byte(UndfindPage))
	ctx.Writer.WriteHeader(404)
}

// handle 向路由集中写入路由数据
func (rt *Route) Handle(method, relativePath string, handlerFunc HandlerFunc) {
	root, ok := rt.Nodes[method]

	if !ok || strings.ReplaceAll(relativePath, " ", "") == "" {
		// 根节点不存在
		root = &Node{
			path: "/",
		}
		rt.Nodes[method] = root
	}

	// 切割path
	var segs []string
	if strings.ReplaceAll(relativePath, " ", "") == "" {
		root.handlerFunc = handlerFunc
		return
	}
	if relativePath[0:1] == "/" {
		segs = strings.Split(relativePath[1:], "/")
	} else {
		segs = strings.Split(relativePath, "/")
	}

	for _, seg := range segs {
		if seg == "" || strings.ReplaceAll(seg, " ", "") == "" {
			continue
		}
		if strings.ReplaceAll(seg, "*", "") == "" {
			seg = "*"
		}

		children := root.childrenCreate(seg)
		root = children
	}
	root.handlerFunc = handlerFunc
}

// childrenCreate 创建子节点方法
func (n *Node) childrenCreate(seg string) *Node {
	if n.children == nil {
		n.children = map[string]*Node{}
	}
	res, ok := n.children[seg]
	if !ok {
		res = &Node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}
