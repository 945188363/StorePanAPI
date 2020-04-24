package main

import (
	"StorePanAPI/config"
	"StorePanAPI/meta"
	"StorePanAPI/mq"
	"encoding/json"
	"io"
	"log"
	"os"
)

func ProcessTransfer(msg []byte) bool {
	pubDate := mq.TransferData{}
	// 反序列化传输文件信息
	err := json.Unmarshal(msg, pubDate)
	if err != nil {
		log.Fatal(err)
		return false
	}
	// 打开源文件
	file, err := os.Open(pubDate.CurLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 在目标地址创建新的文件来拷贝
	newFile, err := os.Create(pubDate.DestLocation)
	if err != nil {
		log.Printf("Fail to create file,error %s", err.Error())
		return false
	}
	defer newFile.Close()

	// 转移存储到新的存储地址
	fileMeta := meta.FileMeta{
		FileShal: pubDate.FileHash,
		FileName: file.Name(),
		Location: pubDate.DestLocation,
	}
	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		log.Printf("Fail to save file,error %s", err.Error())
		return false
	}
	// 更新文件地址信息
	meta.UpdateFileMetaDB(fileMeta)

	return true
}

func main() {
	log.Println("Starting listen mq msg to transfer file")
	mq.StartConsume(
		config.TransferExchangerName,
		"transfer_file",
		ProcessTransfer)
}
