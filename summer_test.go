package summer

import (
	"context"
	"net/http"
	"testing"
)

func TestApplication_Run(t *testing.T) {
	App.SetController("/hello", func(ctx context.Context, request *http.Request) string {
		return "hello"
	})
	App.SetPort("8080")
	App.Run()
}
