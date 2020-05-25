package Test

import (
	"StorePanAPI/rpcService"
	"StorePanAPI/wrapper"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/consul"
	"net/http"
	"testing"
	"time"
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

// timeout middleware wraps the request context with a timeout
func timeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {

				// write response and abort the request
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}

			//cancel to clear resources after finished
			cancel()
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func timedHandler(duration time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {

		// get the underlying request context
		ctx := c.Request.Context()

		// create the response data type to use as a channel type
		type responseData struct {
			status int
			body   map[string]interface{}
		}

		// create a done channel to tell the request it's done
		doneChan := make(chan responseData)

		// here you put the actual work needed for the request
		// and then send the doneChan with the status and body
		// to finish the request by writing the response
		go func() {
			time.Sleep(duration)
			doneChan <- responseData{
				status: 200,
				body:   gin.H{"hello": "world"},
			}
		}()

		// non-blocking select on two channels see if the request
		// times out or finishes
		select {

		// if the context is done it timed out or was cancelled
		// so don't return anything
		case <-ctx.Done():
			c.JSON(300, gin.H{
				"msg": "test",
			})

			return
			// if the request finished then finish the request by
			// writing the response
		case res := <-doneChan:
			c.JSON(res.status, res.body)
		}
	}
}

func TestWeb(t *testing.T) {
	// create new gin without any middleware
	engine := gin.New()

	// add timeout middleware with 2 second duration
	engine.Use(timeoutMiddleware(time.Second * 2))

	// create a handler that will last 1 seconds
	engine.GET("/short", timedHandler(time.Second))

	// create a route that will last 5 seconds
	engine.GET("/long", timedHandler(time.Second*5))

	// run the server
	log.Fatal(engine.Run(":8099"))
}
