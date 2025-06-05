package dbHelper

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

var TencentCOSClient map[string]*cos.Client

func initTencentCOSClient() {
	TencentCOSClient = make(map[string]*cos.Client)
	for _, v := range Conf.TenCentCOS {
		InfoF("[TencentCOS]连接桶%s", v.BucketURL)
		err := connTencentCOSClient(v)
		if err != nil {
			panic(err)
		}
	}
}

func connTencentCOSClient(conf *TenCentCOS) error {
	bucketUrlDev, err := url.Parse(conf.BucketURL)
	if err != nil {
		return err
	}
	bDev := &cos.BaseURL{
		BucketURL: bucketUrlDev,
	}
	TencentCOSClient[conf.Tag] = cos.NewClient(bDev, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.SecretId,
			SecretKey: conf.SecretKey,
		},
	})
	return nil
}

func GetTencentCOSClient(tag string) *cos.Client {
	m, ok := TencentCOSClient[tag]
	if !ok {
		panic("[TencentCOS] 未init")
	}
	return m
}

// TencentCOSCheckIsExist 检查文件是否存在
func TencentCOSCheckIsExist(cosClient *cos.Client, keyName string) bool {
	if len(keyName) == 0 {
		return false
	}
	ok, err := cosClient.Object.IsExist(context.Background(), keyName)
	if err == nil && ok {
		return true
	} else if err != nil {
		ErrorF("TencentCOSCheckIsExist err: %v, keyName: %v", err, keyName)
		return false
	} else {
		return false
	}
}
