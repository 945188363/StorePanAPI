package main

import (
	"StorePanAPI/handler"
	_ "StorePanAPI/handler"
	"StorePanAPI/middleware"
	_ "StorePanAPI/middleware"
	"StorePanAPI/rpcService"
	"StorePanAPI/wrapper"
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/consul"
	"io/ioutil"
	"net/http"
	_ "net/http"
)

func main() {
	// 上传文件基本接口
	http.HandleFunc("/file/upload", middleware.HeadMiddleware(handler.UploadFileHandler))
	http.HandleFunc("/file/upload/success", middleware.HeadMiddleware(handler.UploadSuccessHandler))
	http.HandleFunc("/file/query", middleware.HeadMiddleware(handler.QueryFileMetaHandler))
	http.HandleFunc("/file/queryAll", middleware.HeadMiddleware(handler.QueryAllFileMetaHandler))
	http.HandleFunc("/file/download", middleware.HeadMiddleware(handler.DownloadFileHandler))
	http.HandleFunc("/file/update", middleware.HeadMiddleware(handler.UpdateFileHandler))
	http.HandleFunc("/file/delete", middleware.HeadMiddleware(handler.DeleteFileHandler))

	// 分块上传接口
	http.HandleFunc("/file/mpupload/init", middleware.HeadMiddleware(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart", middleware.HeadMiddleware(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", middleware.HeadMiddleware(handler.CompleteUploadHandler))

	err := http.ListenAndServe(":8899", nil)
	if err != nil {
		fmt.Printf("Fail to start server,error %s", err.Error())
		return
	}
	// Test()
	// fmt.Printf(CallRpcAPI())
}

func Test() {
	resp, err := http.PostForm("https://mock.yonyoucloud.com/mock/6292/media/api/mediaInfo", nil)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))
}

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
