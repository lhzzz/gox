package s3

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("127.0.0.1", "", "")
	assert.Nil(t, err)

	bs, err := client.ListBuckets()
	assert.Nil(t, err)
	t.Log(bs)
}

func TestFilekey(t *testing.T) {
	bucket := "private"
	path := "/test/1.png"

	fk := NewFilekey(bucket, path)
	b, p, err := ParseFileKey(fk)
	assert.Nil(t, err)
	assert.Equal(t, bucket, b)
	assert.Equal(t, path, p)
}

func TestFilekeyMgr(t *testing.T) {
	client, err := NewClient("127.0.0.1", "", "")
	assert.Nil(t, err)

	fkm := NewFilekeyMgrWithArgs("private", "/test/1.txt", client)
	n, err := fkm.Put(bytes.NewBuffer([]byte("12345")), 5, "text/plain")
	assert.Nil(t, err)
	assert.Equal(t, n, int64(5))

	data, err := fkm.ReadAll()
	assert.Nil(t, err)
	assert.Equal(t, data, []byte("12345"))

	t.Log(fkm.String())

	err = fkm.Delete()
	assert.Nil(t, err)

	data, err = fkm.ReadAll()
	t.Log(data, err)
}
