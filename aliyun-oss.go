package dbHelper

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var AliYunOSSClient map[string]*oss.Bucket

func initAliYunOSSClient() {
	AliYunOSSClient = make(map[string]*oss.Bucket)
	for _, v := range Conf.AliYunOSS {
		InfoF("[AliYunOSS]连接桶%s", v.BucketName)
		err := connAliYunOSSClient(v)
		if err != nil {
			panic(err)
		}
	}
}

func GetAliYunOSSClient(tag string) *oss.Bucket {
	m, ok := AliYunOSSClient[tag]
	if !ok {
		panic("[AliYunOSS] 未init")
	}
	return m
}

func connAliYunOSSClient(conf *AliYunOSS) error {

	client, err := oss.New(conf.Endpoint, conf.AccessKeyId, conf.AccessKeySecret)
	if err != nil {
		ErrorF("创建OSS客户端失败: %v", err)
		return err
	}

	// 获取存储空间
	bucket, err := client.Bucket(conf.BucketName)
	if err != nil {
		ErrorF("获取存储空间失败: %v", err)
		return err
	}

	AliYunOSSClient[conf.Tag] = bucket
	return nil
}
