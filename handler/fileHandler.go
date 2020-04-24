package handler

import (
	"StorePanAPI/common"
	"StorePanAPI/config"
	"StorePanAPI/meta"
	"StorePanAPI/mq"
	"StorePanAPI/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

/**
 * 上传文件处理方法
 */
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//return page
	} else if r.Method == "POST" {
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Fail to get file,error %s", err.Error())
			return
		}
		defer file.Close()
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "D:\\PanStoreSys\\" + head.Filename,
			UploadAt: time.Now().Format("2000-11-01,14:32:16"),
		}
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Fail to create file,error %s", err.Error())
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Fail to save file,error %s", err.Error())
			return
		}
		// 保存文件元数据
		newFile.Seek(0, 0)
		fileMeta.FileShal = utils.FileSha1(newFile)
		// meta.SaveFileMeta(fileMeta)
		meta.SaveFileMetaDB(fileMeta)
		data := mq.TransferData{
			FileHash:     fileMeta.FileShal,
			CurLocation:  fileMeta.Location,
			DestLocation: "D:\\PanStoreSys\\storage\\" + fileMeta.FileName,
		}
		pubData, _ := json.Marshal(data)
		mq.Publish(
			config.TransferExchangerName,
			config.TransferRoutingKey,
			pubData,
		)
		http.Redirect(w, r, "file/upload/success", http.StatusAccepted)
	}
}

/**
 * 上传成功处理方法
 */
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "upload success!")
}

/**
 * 查询文件元信息方法
 */
func QueryFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Fail to parse form data,error:%s", err.Error())
		return
	}
	fileHash := r.Form.Get("filehash")
	// fileMeta := meta.GetFileMeta(fileHash)
	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		ret := common.Result{}
		ret.ErrorWithMsg(500, "Fail to query database,error:"+err.Error())
		data, _ := json.Marshal(ret)
		_, _ = w.Write(data)
		return
	}
	ret := common.Result{}
	ret.SuccessWithData("", "查询成功", fileMeta)
	data, _ := json.Marshal(ret)
	_, _ = w.Write(data)
}

/**
 * 查询文件元信息方法
 */
func QueryAllFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// fileMeta := meta.GetFileMeta(fileHash)
		fileMetaSlice, err := meta.GetAllFileMetaDB()
		if err != nil {
			ret := common.Result{}
			ret.ErrorWithMsg(500, "Fail to query database,error:"+err.Error())
			data, _ := json.Marshal(ret)
			_, _ = w.Write(data)
			return
		}
		ret := common.Result{}
		ret.SuccessWithData("", "查询成功", fileMetaSlice)
		data, _ := json.Marshal(ret)
		_, _ = w.Write(data)
	}
}

/**
 * 下载文件方法
 */
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Fail to parse form data,error:%s", err.Error())
		return
	}
	fileHash := r.Form.Get("filehash")
	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		ret := common.Result{}
		ret.ErrorWithMsg(501, "查询数据库失败")
		data, _ := json.Marshal(ret)
		_, _ = w.Write(data)
		return
	}
	file, err := os.Open(fileMeta.Location)
	if err != nil {
		fmt.Printf("Fail to open file,error:%s , location:%s", err.Error(), fileMeta.Location)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		ret := &common.Result{
			Code:  http.StatusInternalServerError,
			Param: "",
			Msg:   "Fail to open file,error:" + err.Error() + "location:" + fileMeta.Location,
		}
		_, _ = w.Write(common.ToBytes(ret))
		fmt.Printf("Fail to open file,error:%s , location:%s", err.Error())
		return
	}
	w.Header().Set("Content-type", "application/octect-stream")
	w.Header().Set("Content-Description", "attachment;filename=\""+fileMeta.FileName+"\"")
	_, _ = w.Write(data)
}

/**
 * 修改文件方法
 */
func UpdateFileHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Fail to parse form data,error:%s", err.Error())
		return
	}
	fileHash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")

	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		ret := common.Result{}
		ret.ErrorWithMsg(501, "查询数据库失败")
		data, _ := json.Marshal(ret)
		_, _ = w.Write(data)
		return
	}
	fileMeta.FileName = filename
	meta.SaveFileMeta(fileMeta)

	ret := common.Result{}
	ret.SuccessWithData("", "更新成功", fileMeta)
	data, _ := json.Marshal(ret)
	_, _ = w.Write(data)
}

/**
 * 删除文件方法
 */
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Fail to parse form data,error:%s", err.Error())
		return
	}
	fileHash := r.Form.Get("filehash")
	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		ret := common.Result{}
		ret.ErrorWithMsg(501, "数据库查询失败")
		data, _ := json.Marshal(ret)
		_, _ = w.Write(data)
		return
	}
	err = os.Remove(fileMeta.Location)
	if err != nil {
		ret := common.Result{}
		ret.ErrorWithMsg(500, "Fail to delete fileMeta:"+err.Error())
		data, _ := json.Marshal(ret)
		_, _ = w.Write(data)
		return
	}
	isSuccess := meta.RemoveFileMetaDB(fileHash)
	if !isSuccess {
		ret := common.Result{}
		ret.ErrorWithMsg(501, "数据库操作失败")
		data, _ := json.Marshal(ret)
		_, _ = w.Write(data)
		return
	}
	ret := common.Result{}
	ret.SuccessWithoutData("删除成功")
	data, _ := json.Marshal(ret)
	_, _ = w.Write(data)
}
