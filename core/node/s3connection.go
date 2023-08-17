package node

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

func (s3c *S3Connection) Connection() s3iface.S3API {
	return s3c.conn
}

func (s3c *S3Connection) FileInfo(url *url.URL) (int64, error) {
	if s3c.conn == nil {
		err := s3c.connect()
		if err != nil {
			return 0, err
		}
	}
	params := &s3.HeadObjectInput{
		Bucket: aws.String(s3c.cred.Bucket),
		Key:    aws.String(url.Path),
	}

	resp, err := s3c.conn.HeadObject(params)
	if err != nil {
		return 0, err
	}
	return *resp.ContentLength, nil
}

func (s3c *S3Connection) Download(url *url.URL) (io.ReadCloser, int64, error) {
	if s3c.conn == nil {
		err := s3c.connect()
		if err != nil {
			return nil, 0, err
		}
	}

	resp, err := s3c.conn.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3c.cred.Bucket),
		Key:    &url.Path})
	if err != nil {
		return nil, 0, err
	}

	return resp.Body, *resp.ContentLength, nil
}

func (s3c *S3Connection) Credential() config.S3Credential {
	return s3c.cred
}

func (s3c *S3Connection) connect() error {
	creds := credentials.NewStaticCredentials(s3c.cred.AccessKeyID, s3c.cred.SecretAccessKey, "")

	_, err := creds.Get()
	if err != nil {
		fmt.Errorf("failed to get S3 credentials")
	}

	cfg := &aws.Config{}
	cfg = cfg.WithCredentials(creds).WithRegion(s3c.cred.Region).WithEndpoint(s3c.cred.Endpoint)

	sess, err := session.NewSession(cfg)
	if err != nil {
		fmt.Errorf("failed to create S3 session")
	}

	s3c.conn = s3.New(sess)

	return nil
}

var _ s3api.S3Backend = (*S3Connection)(nil)
