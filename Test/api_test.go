package Test

import (
	"StorePanAPI/rpcService"
	"StorePanAPI/wrapper"
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/consul"
	"testing"
)

// 调用rpc服务
func CallRpcAPI() (string, error) {
	// 获取consul注册地址
	consulReg := consul.NewRegistry(
		registry.Addrs("localhost:8500"),
	)
	// 获取服务
	prodServiceClient := micro.NewService(
		micro.Name("ProdService.client"),
		micro.Registry(consulReg),
		micro.WrapClient(wrapper.NewLogWrapper),
		micro.WrapClient(wrapper.NewHystrixWrapper),
	)
	prodServiceClient.Init()
	prodService1 := rpcService.NewProdService1Service("ProdServiceRPC", prodServiceClient.Client())
	var req rpcService.ProdRequest1
	req.Size = 2
	prodResponse1, err := prodService1.GetProdList(context.Background(), &req)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return prodResponse1.String(), nil
}

// 使用micro插件调用rpc服务
func CallRpcAPI2(s selector.Selector) {
	// 获取服务
	prodServiceClient := micro.NewService(
		micro.Name("ProdService.client"),
		micro.Selector(s),
	)
	prodService1 := rpcService.NewProdService1Service("ProdServiceRPC", prodServiceClient.Client())
	var req rpcService.ProdRequest1
	req.Size = 2
	prodResponse1, err := prodService1.GetProdList(context.Background(), &req)
	if err != nil {
		log.Fatal(err)

	}
	fmt.Println(prodResponse1.Data)
}

func TestAPI4(t *testing.T) {
	fmt.Println(CallRpcAPI())
}

func TestAPI5(t *testing.T) {
	// 获取consul注册地址
	consulReg := consul.NewRegistry(
		registry.Addrs("localhost:8500"),
	)
	mySelector := selector.NewSelector(
		selector.Registry(consulReg),
		selector.SetStrategy(selector.Random),
	)
	CallRpcAPI2(mySelector)
}
