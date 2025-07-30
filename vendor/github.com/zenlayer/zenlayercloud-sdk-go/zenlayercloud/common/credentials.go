package common

type CredentialIface interface {
	GetSecretId() string
	GetToken() string
	GetSecretKey() string
}

type Credential struct {
	SecretKeyId       string
	SecretKeyPassword string
}

func NewCredential(secretKeyId, secretKeyPassword string) *Credential {
	return &Credential{
		SecretKeyId:       secretKeyId,
		SecretKeyPassword: secretKeyPassword,
	}
}

func (c *Credential) GetSecretKeyId() string {
	return c.SecretKeyId
}

func (c *Credential) GetSecretKeyPassword() string {
	return c.SecretKeyPassword
}
