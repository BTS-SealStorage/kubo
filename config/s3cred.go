package config

type S3Credential struct {
	ID              string `json:"id,omitempty"`
	AccessKeyID     string `json:"aws_access_key_id,omitempty"`
	SecretAccessKey string `json:"aws_secret_access_key,omitempty"`
	Region          string `json:"region,omitempty"`
	Bucket          string `json:"bucket,omitempty"`
	Endpoint        string `json:"endpoint,omitempty"`
}
