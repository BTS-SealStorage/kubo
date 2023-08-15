package node

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	s3api "github.com/ipfs/boxo/s3connection"

	config "github.com/ipfs/kubo/config"
	"io"
	"net/url"
)

type S3Connection struct {
	cred config.S3Credential
	conn s3iface.S3API
}

func (s3c S3Connection) Connection() s3iface.S3API {
	return s3c.conn
}

func (s3c S3Connection) FileInfo(url url.URL) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s3c S3Connection) Download(url url.URL) io.ReadCloser {
	//TODO implement me
	panic("implement me")
}

func (s3c S3Connection) Credential() config.S3Credential {
	return s3c.cred
}

func (s3c S3Connection) Connect() error {
	//TODO implement me
	panic("implement me")
}

var _ s3api.S3Backend = (*S3Connection)(nil)
