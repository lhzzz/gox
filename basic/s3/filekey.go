package s3

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/minio/minio-go/v6"
)

type filekeyMgr struct {
	mio    *minio.Client
	Bucket string
	Path   string
}

type FilekeyMgr interface {
	String() string
	ReadAll() ([]byte, error)
	SaveFile(localFilePath string) error
	Put(reader io.Reader, size int64, contentType string) (n int64, err error)
	FPut(localFilePath string, contentType string) (n int64, err error)
	CopyTo(filekey string) error
	Delete() error
}

var errMinioNil = fmt.Errorf("minio client is nil")

func NewFilekey(bucket, path string) string {
	return NewFilekeyMgrWithArgs(bucket, path, nil).String()
}

func ParseFileKey(fileKey string) (bucket string, path string, err error) {
	pdata, err := base64.StdEncoding.DecodeString(fileKey)
	if err != nil {
		return "", "", fmt.Errorf("base64 DecodeString err : %v, filekey : %s", err, fileKey)
	}
	bucketPath := strings.Split(string(pdata), "_")
	if len(bucketPath) < 2 {
		return "", "", fmt.Errorf("invalid filekey : %s", fileKey)
	}
	bucket, path = bucketPath[0], strings.Join(bucketPath[1:], "_")
	return
}

func NewFilekeyMgrWithArgs(bucket, path string, client *minio.Client) FilekeyMgr {
	return &filekeyMgr{
		mio:    client,
		Bucket: bucket,
		Path:   path,
	}
}

func NewFilekeyMgr(filekey string, client *minio.Client) (FilekeyMgr, error) {
	bucket, path, err := ParseFileKey(filekey)
	if err != nil {
		return nil, err
	}
	return NewFilekeyMgrWithArgs(bucket, path, client), nil
}

func (f *filekeyMgr) String() string {
	return base64.StdEncoding.EncodeToString([]byte(f.Bucket + "_" + f.Path))
}

func (f *filekeyMgr) ReadAll() ([]byte, error) {
	if f.mio == nil {
		return nil, errMinioNil
	}
	object, err := f.mio.GetObject(f.Bucket, f.Path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()
	return ioutil.ReadAll(object)
}

//read from minio and save to local file
func (f *filekeyMgr) SaveFile(localFilePath string) error {
	if f.mio == nil {
		return errMinioNil
	}
	return f.mio.FGetObject(f.Bucket, f.Path, localFilePath, minio.GetObjectOptions{})
}

//upload data to minio filekey
func (f *filekeyMgr) Put(reader io.Reader, size int64, contentType string) (n int64, err error) {
	if f.mio == nil {
		return 0, errMinioNil
	}
	return f.mio.PutObject(f.Bucket, f.Path, reader, size, minio.PutObjectOptions{ContentType: contentType})
}

//upload localfile to minio filekey
func (f *filekeyMgr) FPut(localFilePath string, contentType string) (n int64, err error) {
	if f.mio == nil {
		return 0, errMinioNil
	}
	return f.mio.FPutObject(f.Bucket, f.Path, localFilePath, minio.PutObjectOptions{ContentType: contentType})
}

//copy data to another filekey
func (f *filekeyMgr) CopyTo(filekey string) error {
	if f.mio == nil {
		return errMinioNil
	}
	bucket, path, err := ParseFileKey(filekey)
	if err != nil {
		return err
	}
	dst, err := minio.NewDestinationInfo(bucket, path, nil, nil)
	if err != nil {
		return err
	}
	src := minio.NewSourceInfo(f.Bucket, f.Path, nil)
	return f.mio.CopyObject(dst, src)
}

//remove the filekey itself
func (f *filekeyMgr) Delete() error {
	if f.mio == nil {
		return errMinioNil
	}
	return f.mio.RemoveObject(f.Bucket, f.Path)
}
