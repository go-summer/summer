package summer

import (
	"context"
	"net/http"
	"strings"
)

type ReHandlerFun func(context.Context, *http.Request) string

// Remux 路由结构体
type Remux struct {
	tree           Node
	handlerMapping map[string]http.HandlerFunc
	middleHandler  http.HandlerFunc
}

// CreateNewRemux 创建路由
func CreateNewRemux() *Remux {
	re := Remux{
		handlerMapping: make(map[string]http.HandlerFunc),
	}
	return &re
}

// ServeHttp 实现 net/http http.GetHandler interface
func (re *Remux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	re.middleHandler.ServeHTTP(w, r)
}

// AddMiddleware 添加环绕通知中间件
func (re *Remux) AddMiddleware(f func(handlerFunc http.HandlerFunc) http.HandlerFunc) {
	if re.middleHandler == nil {
		re.middleHandler = f(re.routeMiddleware())
	} else {
		re.middleHandler = f(re.middleHandler)
	}
}

// SetHandlerMapping 添加 handler 和 url 路径的映射
func (re *Remux) SetHandlerMapping(urlStr string, handlerFunc func(context.Context, *http.Request) string) {
	if re.middleHandler == nil {
		re.middleHandler = re.routeMiddleware()
	}
	re.tree.InsertNode(urlStr, Value(packagefun(ReHandlerFun(handlerFunc))))
	//re.handlerMapping[urlStr] = packagefun(ReHandlerFun(handlerFunc))
}

// getHandlerMapping 获取 url 对应的 handler
func (re *Remux) getHandlerMapping(urlStr string) http.HandlerFunc {
	return http.HandlerFunc(re.tree.FindNode(urlStr))
	//return re.handlerMapping[urlStr]
}

// defaultMiddleware 默认中间件用于 查找 url 对应 handler
func (re *Remux) routeMiddleware() http.HandlerFunc {
	f := func(w http.ResponseWriter, req *http.Request) {
		fun := re.getHandlerMapping(strings.Split(req.RequestURI, "?")[0])
		if fun != nil {
			fun.ServeHTTP(w, req)
		}
	}
	return http.HandlerFunc(f)
}

// packagefun 用于包装 GetHandler 确保返回值能够输出
func packagefun(fun ReHandlerFun) http.HandlerFunc {
	f := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(fun(req.Context(), req)))
	}
	return http.HandlerFunc(f)
}
