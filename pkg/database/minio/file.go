package minio

import (
	"context"
	"gin-template-app/pkg/log"
	"io"
	"time"

	"github.com/minio/minio-go/v6"
)

// UploadFile 上传文件
func (c Client) UploadFile(ctx context.Context, bucketName string, fileName string, data io.Reader) (err error) {
	_, err = c.c.PutObjectWithContext(ctx, bucketName, fileName, data, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		log.Error("UploadFile() error(%v)", err)
	}
	return
}

// UploadFileFromMinio 直接上传文件到minio
func (c Client) UploadFileFromMinio(ctx context.Context, bucketName, fileName string) (form map[string]string, err error) {
	// Initialize policy condition config.
	policy := minio.NewPostPolicy()
	// Apply upload policy restrictions:
	policy.SetBucket(bucketName)
	policy.SetKey(fileName)
	// expires in 10 min
	policy.SetExpires(time.Now().Add(time.Hour * 1).UTC())
	// // Only allow 'png' images. => image/png
	// policy.SetContentType("image/png")
	// Only allow content size in range 1KB to 10MB.
	policy.SetContentLengthRange(1024, 10*1024*1024)
	// Get the POST form key/value object:
	url, form, err := c.c.PresignedPostPolicy(policy)
	if err != nil {
		log.Error("UploadFileFromMinio() error(%v)", err)
		return
	}

	// fmt.Printf("curl ")
	// for k, v := range form {
	// 	fmt.Printf("-F %s=%s ", k, v)
	// }
	// fmt.Printf("-F file=@/etc/bash.bashrc ")
	// fmt.Printf("%s\n", url)

	form["url"] = url.String()
	return
}

// DeleteFile 删除文件
func (c Client) DeleteFile(ctx context.Context, bucketName string, fileName string) (err error) {
	err = c.c.RemoveObject(bucketName, fileName)
	if err != nil {
		log.Error("DeleteFile() error(%v)", err)
	}
	return
}

// DownloadFile 下载文件
func (c Client) DownloadFile(ctx context.Context, bucketName string, fileName string) (data io.Reader, err error) {
	object, err := c.c.GetObject(bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Error("GetObject() error:(%v)", err)
		return
	}
	data = object
	return
}

// DownloadLargeFile 下载大文件
func (c Client) DownloadLargeFile(ctx context.Context, bucketName string, fileName string) (path string, err error) {
	// reqParams := make(url.Values)
	// reqParams.Set("response-content-disposition", "attachment; filename=\""+fileName+"\"")
	// Generates a presigned url which expires in 1 hour.
	presignedURL, err := c.c.PresignedGetObject(bucketName, fileName, time.Second*60*60, nil)
	if err != nil {
		log.Error("DownloadLargeFile() error:(%v)", err)
		return
	}
	path = presignedURL.String()
	return
}

// NewBucket 新建一个桶储存
func (c Client) NewBucket(bucketName string, location string) (err error) {
	exists, err := c.c.BucketExists(bucketName)
	if err != nil {
		log.Error("BucketExists() error:(%v)", err)
		return
	}
	if !exists {
		err = c.c.MakeBucket(bucketName, location)
		if err != nil {
			log.Error("MakeBucket() error:(%v)", err)
			return
		}
	}
	return
}
