package summer

import (
	"context"
	"net/http"
)

// Controller 接口定义
type Controller interface {
	Path() string
	Handler(context.Context, *http.Request) string
}
