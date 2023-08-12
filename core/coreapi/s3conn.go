package coreapi

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/ipfs/boxo/s3connection"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ipfs/kubo/tracing"
)

type S3connAPI CoreAPI

func (s3conn *S3connAPI) Cred() s3connection.S3Credential {
	return s3conn.s3conn.Cred() // TODO
}
func (s3conn *S3connAPI) Conn() s3iface.S3API {
	return nil //TODO
}

func (s3conn *S3connAPI) Connect() error {
	_, span := tracing.Span(s3conn.nctx, "CoreAPI.S3ConnAPI", "Add",
		trace.WithAttributes(attribute.String("bucket", s3conn.Cred().Bucket)))
	defer span.End()

	//TODO
	return nil
}
