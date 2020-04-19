package main

import (
	"StorePanAPI/handler"
	"StorePanAPI/middleware"
	"fmt"
	"net/http"
)

func main() {
	// 上传文件基本接口
	http.HandleFunc("/file/upload",middleware.HeadMiddleware(handler.UploadFileHandler))
	http.HandleFunc("/file/upload/success",middleware.HeadMiddleware(handler.UploadSuccessHandler))
	http.HandleFunc("/file/query",middleware.HeadMiddleware(handler.QueryFileMetaHandler))
	http.HandleFunc("/file/queryAll",middleware.HeadMiddleware(handler.QueryAllFileMetaHandler))
	http.HandleFunc("/file/download",middleware.HeadMiddleware(handler.DownloadFileHandler))
	http.HandleFunc("/file/update",middleware.HeadMiddleware(handler.UpdateFileHandler))
	http.HandleFunc("/file/delete",middleware.HeadMiddleware(handler.DeleteFileHandler))

	// 分块上传接口
	http.HandleFunc("/file/mpupload/init", middleware.HeadMiddleware(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart", middleware.HeadMiddleware(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", middleware.HeadMiddleware(handler.CompleteUploadHandler))

	err := http.ListenAndServe(":8899",nil)
	if err != nil{
		fmt.Printf("Fail to start server,error %s",err.Error())
		return
	}
}