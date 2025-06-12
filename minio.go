package dbHelper

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

var MinioClient map[string]*minio.Client

func initMinioClient() {
	MinioClient = make(map[string]*minio.Client)
	for _, v := range Conf.MinIOConf {
		InfoF("[Minio]连接 %s", v.Endpoint)
		err := connMinioClient(v)
		if err != nil {
			panic(err)
		}
	}
}

func GetMinioClient(tag string) *minio.Client {
	m, ok := MinioClient[tag]
	if !ok {
		panic("[Minio] 未init")
	}
	return m
}

func connMinioClient(conf *MinIOConf) error {
	// 创建 MinIO 客户端
	minioClient, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyId, conf.AccessKeySecret, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		Error("创建 MinIO 客户端失败:", err)
		return err
	}

	MinioClient[conf.Tag] = minioClient
	return nil
}

// MinioCreateBucket 创建存储桶
func MinioCreateBucket(client *minio.Client, ctx context.Context, bucketName string) error {
	// 检查存储桶是否已存在
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		// 创建存储桶
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		InfoF("成功创建存储桶: %s", bucketName)
	} else {
		InfoF("存储桶已存在: %s", bucketName)
	}

	return nil
}

// MinioUploadFile 上传文件
func MinioUploadFile(client *minio.Client, ctx context.Context, bucketName, objectName, filePath, contentType string) error {
	// 上传文件到指定的存储桶
	info, err := client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	InfoF("成功上传文件: %s，大小: %d 字节", objectName, info.Size)
	return nil
}

// MinioDownloadFile 下载文件
func MinioDownloadFile(client *minio.Client, ctx context.Context, bucketName, objectName, destinationPath string) error {
	// 下载文件到本地
	err := client.FGetObject(ctx, bucketName, objectName, destinationPath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	InfoF("成功下载文件到: %s", destinationPath)
	return nil
}

// MinioGenerateSignedURL 生成预签名 URL
func MinioGenerateSignedURL(client *minio.Client, ctx context.Context, bucketName, objectName string) (string, error) {
	// 生成一个小时后过期的 GET 请求预签名 URL
	reqParams := make(map[string][]string)
	signedURL, err := client.PresignedGetObject(ctx, bucketName, objectName, time.Hour, reqParams)
	if err != nil {
		return "", err
	}

	return signedURL.String(), nil
}

// MinioListObjects 列出存储桶中的对象
func MinioListObjects(client *minio.Client, ctx context.Context, bucketName string) error {
	// 创建一个对象迭代器
	doneCh := make(chan struct{})
	defer close(doneCh)

	// 列出存储桶中的所有对象
	objectCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	InfoF("存储桶 %s 中的对象列表:", bucketName)
	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}
		InfoF("  - %s (大小: %d 字节, 修改时间: %s)", object.Key, object.Size, object.LastModified)
	}

	return nil
}

// MinioDeleteObject 删除对象
func MinioDeleteObject(client *minio.Client, ctx context.Context, bucketName, objectName string) error {
	// 删除指定的对象
	err := client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	InfoF("成功删除对象: %s", objectName)
	return nil
}

// MinioDeleteBucket 删除存储桶
func MinioDeleteBucket(client *minio.Client, ctx context.Context, bucketName string) error {
	// 删除存储桶（必须为空）
	err := client.RemoveBucket(ctx, bucketName)
	if err != nil {
		return err
	}

	InfoF("成功删除存储桶: %s\n", bucketName)
	return nil
}
