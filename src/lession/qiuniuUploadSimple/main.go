package main

import (
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/api.v7/auth/qbox"
	"fmt"
	"context"
)

func main() {

	//simeUploadFile()

	coversimeUploadFile()

	customPutRet()
}

/**
 * 简单的上传文件
 */
func simeUploadFile() {
	accessKey := "HlE45UT8wRJBPWBb4HIup2dKn33cWcBaq6Wo-jye"
	secretKey := "IqPCJAY-0Q90VX9vF7BNSg2a_uzGlVH8TwvOi_j0"

	localFile := "persion_log.jpg"

	key := "persion_log.jpg"

	bucket := "drustydatarepo"

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}

	mac := qbox.NewMac(accessKey, secretKey)

	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}

	cfg.Zone = &storage.ZoneHuanan

	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}

	putExTtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExTtra)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ret.Key, ret.Hash)
}

/**
 * 覆盖已存在的文件
*/
func coversimeUploadFile() {

	accessKey := "HlE45UT8wRJBPWBb4HIup2dKn33cWcBaq6Wo-jye"
	secretKey := "IqPCJAY-0Q90VX9vF7BNSg2a_uzGlVH8TwvOi_j0"
	localFile := "persion_log.jpg"
	bucket := "drustydatarepo"
	key := "persion_log.jpg"
	keyToOverwrite := "persion_log.jpg"

	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", bucket, keyToOverwrite),
	}

	mac := qbox.NewMac(accessKey, secretKey)

	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}

	cfg.Zone = &storage.ZoneHuanan

	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}

	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)
}

type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

/**
 * 自定义返回结构体
 */
func customPutRet() {

	accessKey := "HlE45UT8wRJBPWBb4HIup2dKn33cWcBaq6Wo-jye"
	secretKey := "IqPCJAY-0Q90VX9vF7BNSg2a_uzGlVH8TwvOi_j0"
	localFile := "persion_log.jpg"
	bucket := "drustydatarepo"
	key := "persion_log.jpg"

	putPolicy := storage.PutPolicy{
		Scope:      bucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}

	mac := qbox.NewMac(accessKey, secretKey)

	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}

	cfg.Zone = &storage.ZoneHuanan

	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)

	ret := MyPutRet{}

	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ret)
}
