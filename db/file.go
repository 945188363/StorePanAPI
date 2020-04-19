package db

import (
	mydb "StorePanAPI/db/mysql"
	"database/sql"
	"fmt"
)
/**
 * 文件上传存储到数据库中
 */
func UploadSQL(filehash string, filename string, filesize int64, fileaddr string) bool  {
	stmt, err := mydb.DBConn().Prepare("INSERT ignore INTO file(file_sha1,file_name,file_size,file_addr,status) VALUES (?,?,?,?,1)")
	if err != nil {
		fmt.Printf("Fail to prepare statement,error:",err.Error())
		return false
	}
	defer stmt.Close()

	result, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Printf("Fail to execute sql, error:",err.Error())
		return false
	}
	if affected, err := result.RowsAffected();err ==nil {
		if affected > 0{
			fmt.Printf("File with hash:%s has benn uploaded successfully. ",filehash)
		}
		return true
	}
	return false
}

type TableFile struct {
	Filehash sql.NullString
	Filename sql.NullString
	Filesize sql.NullInt64
	Fileaddr sql.NullString
}

/**
 *  根据文件hash值获取数据库信息
 */
func GetFileMetaSQL(filehash string) (*TableFile,error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size,file_addr from file where file_sha1=? and status =1 limit 1")
	if err != nil {
		fmt.Printf("Fail to prepare statement,error:",err.Error())
		return nil,err
	}
	defer stmt.Close()
	tabFile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tabFile.Filehash, &tabFile.Filename, &tabFile.Filesize, &tabFile.Fileaddr)
	if err != nil {
		fmt.Printf("Fail to execute sql, error:",err.Error())
		return nil,err
	}
	return &tabFile,nil
}

/**
 *  获取数据库所有信息
 */
func GetAllFileMetaSQL() (*[]TableFile,error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size,file_addr from file where status =1 order by file_name")
	if err != nil {
		fmt.Printf("Fail to prepare statement,error:",err.Error())
		return nil,err
	}
	defer stmt.Close()
	var tabFileSlice []TableFile
	tabFile := TableFile{}
	row, _ := stmt.Query()
	count :=0
	for row.Next(){
		if err = row.Scan(&tabFile.Filehash, &tabFile.Filename, &tabFile.Filesize, &tabFile.Fileaddr);err !=nil{
			fmt.Printf("Fail to execute sql, error:",err.Error())
			return nil,err
		}
		tabFileSlice = append(tabFileSlice,tabFile)
		count++
	}

	return &tabFileSlice,nil
}

/**
 *  根据文件hash值删除数据库信息
 */
func DeleteFileMetaSQL(filehash string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"update file set status = 0 where file_sha1=?")
	if err != nil {
		fmt.Printf("Fail to prepare statement,error:",err.Error())
		return false
	}
	defer stmt.Close()
	result, err := stmt.Exec(filehash)
	if err != nil {
		fmt.Printf("Fail to execute sql, error:",err.Error())
		return false
	}

	if rw,err := result.RowsAffected();err == nil{
		if rw > 0{
			fmt.Printf("File with hash:%s has benn delete successfully. ",filehash)
		}
		return true
	}
	return false
}