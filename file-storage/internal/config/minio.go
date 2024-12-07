package config

type Minio struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_ket"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

func (m Minio) GetEndpoint() string {
	return m.Endpoint
}

func (m Minio) GetAccessKey() string {
	return m.AccessKey
}

func (m Minio) GetSecretKey() string {
	return m.SecretKey
}

func (m Minio) GetBucket() string {
	return m.Bucket
}

func (m Minio) GetUseSSL() bool {
	return m.UseSSL
}
