package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/yezzey-gp/yproxy/config"
)

type SessionPool interface {
	GetSession() *s3.S3
}

type S3SessionPool struct {
	cnf config.Storage
}

// TODO : unit tests
func (sp *S3SessionPool) createSession() (*session.Session, error) {
	s, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	provider := &credentials.StaticProvider{Value: credentials.Value{
		AccessKeyID:     sp.cnf.AccessKeyId,
		SecretAccessKey: sp.cnf.SecretAccessKey,
	}}
	providers := make([]credentials.Provider, 0)
	providers = append(providers, provider)
	providers = append(providers, defaults.CredProviders(s.Config, defaults.Handlers())...)
	newCredentials := credentials.NewCredentials(&credentials.ChainProvider{
		VerboseErrors: aws.BoolValue(s.Config.CredentialsChainVerboseErrors),
		Providers:     providers,
	})

	s.Config.WithCredentials(newCredentials)
	return s, err
}

func (s *S3SessionPool) GetObject() (*s3.S3, error) {

	sess, err := s.createSession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new session")
	}
	return s3.New(sess), nil
}
