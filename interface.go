package summer

import (
	"context"
	"net/http"
)

// Controller 接口定义
type Controller interface {
	GetPath() string
	GetHandler() func(context.Context, *http.Request) string
}
