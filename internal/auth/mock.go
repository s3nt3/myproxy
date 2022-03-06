package auth

import (
	"github.com/go-mysql-org/go-mysql/server"
)

type MockProxyAuth struct {
	provider server.CredentialProvider
}

func NewMockProxyAuth() *MockProxyAuth {
	provider := &server.InMemoryProvider{}
	provider.AddUser("root", "proxy")

	return &MockProxyAuth{
		provider: provider,
	}
}

func (mock *MockProxyAuth) GetProvider() server.CredentialProvider {
	return mock.provider
}
