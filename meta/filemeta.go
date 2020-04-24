package meta

import (
	mydb "StorePanAPI/db"
)

type FileMeta struct {
	FileShal string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

/**
 * 以文件hash值为key存储
 */
func SaveFileMeta(meta FileMeta) {
	fileMetas[meta.FileShal] = meta
}

/**
 * 存储到数据库
 */
func SaveFileMetaDB(meta FileMeta) bool {
	return mydb.UploadSQL(meta.FileShal, meta.FileName, meta.FileSize, meta.Location)
}

/**
 * 更新到数据库
 */
func UpdateFileMetaDB(meta FileMeta) bool {
	return mydb.UpdateSQL(meta.FileShal, meta.FileName, meta.FileSize, meta.Location)
}

/**
 * 根据hash值获取文件
 */
func GetFileMeta(fileShal string) FileMeta {
	return fileMetas[fileShal]
}

/**
 * 根据hash值获取文件
 */
func GetFileMetaDB(fileShal string) (FileMeta, error) {
	tabFile, err := mydb.GetFileMetaSQL(fileShal)
	if err != nil {
		return FileMeta{}, err
	}
	fileMeta := FileMeta{
		FileShal: tabFile.Filehash.String,
		FileName: tabFile.Filename.String,
		FileSize: tabFile.Filesize.Int64,
		Location: tabFile.Fileaddr.String,
	}
	return fileMeta, nil
}

/**
 * 根据hash值获取文件
 */
func GetAllFileMetaDB() ([]FileMeta, error) {
	tabFileSlice, err := mydb.GetAllFileMetaSQL()
	if err != nil {
		return make([]FileMeta, 0), err
	}
	var fileMetas []FileMeta
	for _, tabFile := range *tabFileSlice {
		fileMeta := FileMeta{
			FileShal: tabFile.Filehash.String,
			FileName: tabFile.Filename.String,
			FileSize: tabFile.Filesize.Int64,
			Location: tabFile.Fileaddr.String,
		}
		fileMetas = append(fileMetas, fileMeta)
	}
	return fileMetas, nil
}

/**
 * 根据hash值移除文件
 */
func RemoveFileMetaDB(fileShal string) bool {
	// delete(fileMetas,fileShal)
	return mydb.DeleteFileMetaSQL(fileShal)
}
