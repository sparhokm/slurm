package minio

type Config interface {
	GetEndpoint() string
	GetAccessKey() string
	GetSecretKey() string
	GetUseSSL() bool
	GetBucket() string
}
