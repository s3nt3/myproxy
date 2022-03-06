package main

import (
	"github.com/s3nt3/myproxy/internal/auth"
	"github.com/s3nt3/myproxy/internal/logger"
	"github.com/s3nt3/myproxy/internal/server"
)

func init() {
	logger.InitLogger(true, false)
}

func main() {
	proxyServer := server.NewProxyServer(
		"127.0.0.1", 4000,
		auth.NewMockProxyAuth(),
		server.NewMockProxyServerConfig(),
		server.NewMockProxyServerHandler("127.0.0.1", 3306),
	)

	proxyServer.Run()
}
