package s3

import (
	"github.com/minio/minio-go/v6"
)

func NewClient(endpoint, accessKeyId, accessKey string) (*minio.Client, error) {
	return minio.New(endpoint, accessKeyId, accessKey, false)
}
