package tee

import (
	"log"
	"net/http"
	"strings"
	"sync"
)

type Handler func(context *Context)

var GroupRouteInndex = 1

// 通过引擎来控制路由前缀树入口
type Engine struct {
	CurNode *TreeNode // 路由当前节点，比如在返回group函数后添加路由，每个group函数返回新的engine的当前节点
	// 的是前缀树对应的当前路由，group后添加的路由，从当前节点开始添加
	RootNode     *TreeNode   // 路由根节点
	HandlerSlice [][]Handler // 中间件，控制器方法数组
	pool         sync.Pool
	PreIndex     int
	Index        int
	PreEngine    *Engine
	UrlParamsMap map[int][]string
}

// 前缀树节点，比如路由为/tee/api/:type/qq，那么路由会拆解成tee，api，:type，qq，四个节点,
// 为了支持路径参数，把带:的路由（例如:type)统一存储成/,并且用map存储执行函数索引和参数对应关系，存储到PathParams
type TreeNode struct {
	PathUrl  map[string]*TreeNode //当前路由对应下一个路由节点的url为key，下一个路由节点为value
	UrlValue string               // 当前路由url，例如api或qq
	// 当前路由为:type,:id等路径方式，存储参数，key为带有路由对应发放，
	// value为带有路径参数的url
	End          bool      // 是否是路url由的重点
	Index        int       // 当前执行的路由方法索引
	PreIndex     int       // 当前节点url对应的中间件
	
	PreEngine    *Engine
}

var RouteUrlParamsMap = make(map[int][]string)
var routeMap = make(map[string][]Handler)

// 前缀树路由根节点
var Root = &TreeNode{
	PathUrl: make(map[string]*TreeNode),
}

var Estart = GetNewEngine()

// 插入一个前缀树路由
func (curNode *TreeNode) Insert(routeStringSlice []string, index int, handlerIndex int, e *Engine, routeindex int) *TreeNode {
	curV := routeStringSlice[index]
	// 路径参数,保存为 /
	if []byte(routeStringSlice[index])[0] == ':' {
		curV = "/"
	}
	v, ok := curNode.PathUrl[curV]
	if ok {
		curNode = v
	} else {
		newNode := &TreeNode{
			PathUrl:      make(map[string]*TreeNode),
			UrlValue:     curV,
			
		}

		// 上一个节点的map指向刚创建的节点
		curNode.PathUrl[newNode.UrlValue] = newNode
		curNode = newNode
	}
	// 如果是路由参数，对应方法索引指向该路径参数
	if curV == "/" {
		RouteUrlParamsMap[routeindex] = append(RouteUrlParamsMap[routeindex], string([]byte(routeStringSlice[index])[1:]))
	}

	// 路由插入前缀树完毕
	if index+1 == len(routeStringSlice) {
		// 如果是group，代表路由还没结束
		if handlerIndex == 0 && !curNode.End {
			curNode.End = false

			return curNode
		}
		curNode.PreIndex = routeindex
		// 如果是addRoute，代表路由到此结束，End设置为true，并保存路由索引
		curNode.PreEngine = e
		curNode.End = true
		curNode.Index = handlerIndex
		return curNode
	}

	return curNode.Insert(routeStringSlice, index+1, handlerIndex, e, routeindex)

}

// 要实现框架，需要实现监听的serveHTTP方法
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	routeString := req.RequestURI + "/" + req.Method

	if v, ok := routeMap[routeString]; ok {
		handler, routeParam := v, make(map[string]interface{})
		context := Estart.pool.Get().(*Context)
		context.Req = req
		context.Res = w
		context.RouteParamMap = routeParam
		context.Index = -1
		context.HandlerSlice = handler
		context.Next()
		w = context.Res
		Estart.pool.Put(context)
	} else {
		isMatch, handler, routeParam := Estart.MatchRoute(routeString)
		if !isMatch {
			return
		}
		context := Estart.pool.Get().(*Context)
		context.Req = req
		context.Res = w
		context.RouteParamMap = routeParam
		context.Index = -1
		context.HandlerSlice = handler

		context.Next()
		w = context.Res
		Estart.pool.Put(context)
	}

	// 获取匹配路由方法

}

// 注册路由组到前缀树
func (e *Engine) Group(routeString string) *Engine {
	return e.GroupHandleIndex(routeString, 0)
}

// 路由注册方法
func (e *Engine) GroupHandleIndex(routeString string, handleIndex int) *Engine {
	routeStringSlice := strings.Split(routeString, "/")
	start := 0
	end := len(routeStringSlice)
	if len(routeStringSlice[start]) == 0 {
		start = 1
	}
	if len(routeStringSlice[end-1]) == 0 {
		end = end - 1
	}

	newE := &Engine{
		CurNode:   e.CurNode,
		RootNode:  Root,
		PreEngine: e,
		Index:     GroupRouteInndex + 1,
	}
	GroupRouteInndex++
	curIndex := GroupRouteInndex
	// 插入前缀路由，返回新的引擎
	newE.CurNode = newE.CurNode.Insert(routeStringSlice[start:end], 0, handleIndex, newE, curIndex)

	return newE
}

// 匹配路由
func (e *Engine) MatchRoute(routeString string) (bool, []Handler, map[string]interface{}) {
	routeStringSlice := strings.Split(routeString, "/")
	e.CurNode = e.RootNode

	RouteParamMap := make(map[string]interface{})
	// isMatch为是否匹配，handerIndex为路由对应方法，handler为中间件
	isMatch, handlerIndex, _ := e.CurNode.Match(routeStringSlice, 1, RouteParamMap, 0)
	if !isMatch {
		return false, nil, RouteParamMap
	}

	return isMatch, e.HandlerSlice[handlerIndex-1], RouteParamMap
}

// 使用中间件
func (e *Engine) Use(handlers ...Handler) *Engine {

	if e.PreIndex != 0 {
		Estart.HandlerSlice[e.PreIndex-1] = append(Estart.HandlerSlice[e.PreIndex-1], handlers...)
	} else {
		Estart.HandlerSlice = append(Estart.HandlerSlice, handlers)
		e.PreIndex = len(Estart.HandlerSlice)
	}

	return e
}

// 匹配路由成功返回路径参数
func (curNode *TreeNode) Match(routeStringSlice []string, index int, RouteParamMap map[string]interface{}, urlindex int) (bool, int, int) {

	if len(routeStringSlice) == index {
		if curNode.End {
			return true, curNode.Index, curNode.PreIndex
		}
		return false, 0, 0
	}
	// 匹配不带路径参数的路由
	tempNode1, tempNode2 := curNode, curNode
	v, ok := tempNode1.PathUrl[routeStringSlice[index]]
	if ok {
		tempNode1 = v
		isMatch, handlerIndex, routeindex := tempNode1.Match(routeStringSlice, index+1, RouteParamMap, urlindex)
		if isMatch {
			return isMatch, handlerIndex, routeindex
		}
	}
	// 匹配带路径参数的路由
	v, ok = tempNode2.PathUrl["/"]
	if ok {
		urlindex++
		tempNode2 = v
		isMatch, handlerIndex, routeindex := tempNode2.Match(routeStringSlice, index+1, RouteParamMap, urlindex)
		if isMatch {

			RouteParamMap[Estart.UrlParamsMap[routeindex][urlindex-1]] = routeStringSlice[index]
			return isMatch, handlerIndex, routeindex
		}
		urlindex--
	}
	return false, 0, 0
}

// 增加路由
func (e *Engine) AddRoute(routeName string, handler ...Handler) *Engine {

	Estart.HandlerSlice = append(Estart.HandlerSlice, handler)
	e = e.GroupHandleIndex(routeName, len(Estart.HandlerSlice))
	return e
}

func NewEngine() *Engine {
	return &Engine{
		HandlerSlice: make([][]Handler, 0),
		RootNode:     Root,
		CurNode:      Root,
		UrlParamsMap: make(map[int][]string),
	}

}

func GetNewEngine() *Engine {
	e := NewEngine()
	e.pool.New = func() interface{} {
		return &Context{}
	}
	return e
}

func RouteInit(root *TreeNode, routeSlice []string, mode bool) {
	if root == nil {
		return
	}
	if root.End {
		newHandelers, newPathParam := CallBackTreeNode(root.PreEngine)
		Estart.HandlerSlice[root.Index-1] = append(newHandelers, Estart.HandlerSlice[root.Index-1]...)

		Estart.UrlParamsMap[root.PreIndex] = append(newPathParam, Estart.UrlParamsMap[root.PreIndex]...)

		if mode {
			routeMap["/"+strings.Join(routeSlice, "/")] = Estart.HandlerSlice[root.Index-1]
		}

	}

	for k, v := range root.PathUrl {
		if k != "/" {
			RouteInit(v, append(routeSlice, k), mode)
		} else {
			RouteInit(v, append(routeSlice, k), false)
		}

	}
}

func RouteDelete(root *TreeNode, mode bool) {

	if root.End {
		e := root.PreEngine
		if e != nil {
			if e.PreIndex > 0 {
				Estart.HandlerSlice[e.PreIndex-1] = nil
			}
			e = e.PreEngine
		}
		if mode {
			root = nil
		}
		return
	}

	for k, v := range root.PathUrl {
		if k != "/" {
			RouteDelete(v, mode)
		} else {
			RouteDelete(v, false)

		}

	}

}

func CallBackTreeNode(e *Engine) ([]Handler, []string) {
	if e == nil {
		return nil, nil
	}
	var handlers []Handler
	var strArr []string
	if e.PreIndex > 0 {
		handlers = Estart.HandlerSlice[e.PreIndex-1]

	}

	if e.Index > 0 {
		strArr = RouteUrlParamsMap[e.Index]
	}
	newHandelers, newPathParam := CallBackTreeNode(e.PreEngine)
	handlers = append(newHandelers, handlers...)
	strArr = append(newPathParam, strArr...)
	return handlers, strArr
}

// 启动程序，监听http方法
func Start(address string) {
	RouteInit(Root, nil, true)

	RouteDelete(Root, true)
	RouteUrlParamsMap = nil

	srv := &http.Server{
		Addr:         address,
		Handler:      Estart,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func (e *Engine) GET(routeName string, handlers ...Handler) {
	e.AddRoute(routeName+"/GET", handlers...)
}

func (e *Engine) POST(routeName string, handlers ...Handler) {
	e.AddRoute(routeName+"/POST", handlers...)
}

func (e *Engine) DELETE(routeName string, handlers ...Handler) {
	e.AddRoute(routeName+"/DELETE", handlers...)
}
