package oss

import (
	"../models"
	"../php2go"
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

/**
 * 分析文件，并返回文件信息
 */
func AnalysisFile(file multipart.File, header *multipart.FileHeader) (models.Files, error) {
	fileInfo := models.Files{}
	var err error
	if file == nil {
		file, err = header.Open()
		if err != nil {
			return fileInfo, err
		}
	}
	defer file.Close()
	fileNameSep := php2go.Explode(".", header.Filename)
	// 后缀名
	fileInfo.Suffix = fileNameSep[len(fileNameSep)-1]
	fileNameSep = fileNameSep[:len(fileNameSep)-1]
	// 文件名
	fileInfo.Name = php2go.Implode(".", fileNameSep)
	// 文件大小
	fileInfo.Size = strconv.FormatInt(header.Size, 10)
	fileSha1, err := php2go.Sha1FileSrc(file)
	if err != nil {
		return fileInfo, err
	}
	fileInfo.Hash = fileSha1
	sha1Arr := php2go.Split(fileSha1, 4)
	// token名
	fileInfo.TokenName = php2go.Md5(fileSha1)
	// 文件路径
	fileInfo.Path = "./uploads/" + php2go.Implode("/", sha1Arr) + "/"
	err = os.MkdirAll(fileInfo.Path, os.ModePerm)
	if err != nil {
		return fileInfo, err
	}
	fileInfo.Uri = fileInfo.Path + fileInfo.TokenName + "." + fileInfo.Suffix
	out, err := os.OpenFile(fileInfo.Uri, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fileInfo, err
	}
	defer out.Close()
	io.Copy(out, file)
	return fileInfo, nil
}
