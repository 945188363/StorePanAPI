package wrapper

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/client"
)

type LogWrapper struct {
	client.Client
}

func (this *LogWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	fmt.Println("日志测试")
	return this.Client.Call(ctx, req, rsp)
}

func NewLogWrapper(c client.Client) client.Client {
	return &LogWrapper{c}
}
