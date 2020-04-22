package Wrapper

import (
	"StorePanAPI/Service"
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/client"
	"strconv"
)

type HystrixWrapper struct {
	client.Client
}

func NewProd(id int32, name string) *Service.ProdModel {
	return &Service.ProdModel{
		ProdId:   id,
		ProdName: name,
	}
}

// 降级方法
func DefaultProdList(rsp interface{}) {
	ret := make([]*Service.ProdModel, 0)
	var i int32
	for i = 0; i < 5; i++ {
		ret = append(ret, NewProd(20+i, "prod"+strconv.Itoa(int(i))))
	}
	result := rsp.(*Service.ProdResponse1)
	result.Data = ret
}

func (this *HystrixWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	cmdName := req.Service() + "." + req.Endpoint()
	// hystrix通用配置
	hystrixConfig := hystrix.CommandConfig{
		Timeout:                5000,
		RequestVolumeThreshold: 5,     // 达到几次请求后开始错误率计算
		ErrorPercentThreshold:  20,    // 错误率达到多少 直接执行降级方法
		SleepWindow:            10000, // 熔断器再次判断时间 10S

	}
	hystrix.ConfigureCommand(cmdName, hystrixConfig)
	return hystrix.Do(cmdName, func() error {
		return this.Client.Call(ctx, req, rsp)
	}, func(err error) error {
		// 执行降级方法
		DefaultProdList(rsp)
		return nil
	})
}

func NewHystrixWrapper(c client.Client) client.Client {
	return &HystrixWrapper{c}
}
