package auth

import (
	"github.com/go-mysql-org/go-mysql/server"
)

type AuthPlugin interface {
	GetProvider() server.CredentialProvider
}
