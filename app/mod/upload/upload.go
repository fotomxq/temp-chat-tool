package ModUpload

import (
	"crypto/sha1"
	"github.com/fotomxq/gobase/file"
	"github.com/fotomxq/gobase/filter"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"strconv"
	"strings"
	"time"
)

//上传文件类别封装
type ModUploadFileType struct {
	//文件尺寸
	Size int64
	//文件名称，含类别
	Name string
	//文件名称，不含类别
	OnlyName string
	//新的文件名称
	NewName string
	//文件类别
	Type string
	//创建时间
	CreateTime int64
	//存储路径
	Src string
	//SHA1摘要
	SHA1 string
}

//上传文件
// 可利用该方法，实现任意文件、目标得上传
// 注意，目标路径不要交给用户设计，否则将出现重大系统级漏洞，影响系统安全
// 存储后，默认根据创建"Unix时间戳_"结构设计文件名称
//param c *gin.Context
//param formName string 表单名称
//param targetSrc string 目标路径，末尾必须添加Sep
//param maxSize int64 文件最大大小，如果为0则不限制
//param filterType []string 文件类别限制
//param isRename bool 是否重新命名文件名称
//return FileUploadType 文件类型
//return error 错误信息
func ModUploadFile(c *gin.Context, formName string, targetSrc string, maxSize int64, filterType []string, isRename bool) (ModUploadFileType, error) {
	//初始化
	res := ModUploadFileType{}
	//获取文件
	formFile, header, err := c.Request.FormFile(formName)
	if err != nil {
		return res, err
	}
	defer formFile.Close()
	//判断文件尺寸
	if header.Size > maxSize && maxSize > 0 {
		err = errors.New("upload file size too lager.")
		return res, err
	}
	res.Size = header.Size
	//获取文件名
	res.Name = header.Filename
	//过滤Names，不允许非英文等特殊字符
	res.Name = filter.CheckFilterStr(res.Name, 1, 250)
	//继续拆解文件名、类型
	names := strings.Split(res.Name, ".")
	res.Type = names[len(names)-1]
	res.OnlyName = names[0]
	for k, v := range names {
		if k == 0 {
			continue
		}
		if k >= len(names)-1 {
			break
		}
		res.OnlyName += "." + v
	}
	//甄别filterType
	if len(filterType) > 0 {
		isOK := false
		for _, v := range filterType {
			if v == res.Type {
				isOK = true
			}
		}
		if isOK == false {
			err = errors.New("upload file type is ban.")
			return res, err
		}
	}
	//创建时间
	res.CreateTime = time.Now().Unix()
	//计算SHA1
	buf := make([]byte, res.Size)
	n, err := formFile.Read(buf)
	if err != nil {
		return res, err
	}
	dataByte := buf[:n]
	hash := sha1.New()
	if _, err := io.WriteString(hash, string(dataByte)); err != nil {
		return res, err
	}
	res.SHA1 = string(hash.Sum(nil))
	//构建新的文件名称
	if isRename == true {
		res.NewName = strconv.FormatInt(res.CreateTime, 10) + "_" + res.SHA1 + "." + res.Type
	} else {
		res.NewName = res.Name
	}
	//构建文件路径
	res.Src = targetSrc + res.NewName
	//文件不能已经存在
	if file.IsExist(res.Src) == true {
		return res, errors.New("upload file is exist.")
	}
	//创建并保存文件
	err = file.WriteFile(res.Src, dataByte)
	if err != nil {
		return res, err
	}
	//返回
	return res, nil
}
