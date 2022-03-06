package server

import (
	"fmt"
	"net"

	"github.com/go-mysql-org/go-mysql/server"

	"github.com/s3nt3/myproxy/internal/auth"
	"github.com/s3nt3/myproxy/internal/logger"
	"github.com/s3nt3/myproxy/internal/utils"
)

type ProxyServerConfig interface {
	GetProtocolConfig() *server.Server
}

type ProxyServerHandler interface {
	GetProxyHandler() server.Handler
}

type ProxyServer struct {
	Host string
	Port uint

	Auth    auth.AuthPlugin
	Config  ProxyServerConfig
	Handler ProxyServerHandler
}

func NewProxyServer(
	host string,
	port uint,
	authPlugin auth.AuthPlugin,
	config ProxyServerConfig,
	handler ProxyServerHandler) *ProxyServer {
	return &ProxyServer{
		Host:    host,
		Port:    port,
		Auth:    authPlugin,
		Config:  config,
		Handler: handler,
	}
}

func (s *ProxyServer) Run() {
	proxy, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	utils.FatalCheck(err, fmt.Sprintf("MySQL proxy server start, listening at %s:%d", s.Host, s.Port))

	for {
		conn, err := proxy.Accept()
		ok := utils.ErrorCheck(err, fmt.Sprintf("Accept connection from %s", conn.RemoteAddr().String()))

		if ok {
			proxyConn, err := server.NewCustomizedConn(conn,
				s.Config.GetProtocolConfig(),
				s.Auth.GetProvider(),
				s.Handler.GetProxyHandler())
			if err != nil {
				logger.Logger.Error(err)
			} else {
				go func(conn *server.Conn) {
					defer func(remoteAddr string) {
						logger.Logger.Debugf("Close connection from %s", remoteAddr)
					}(conn.RemoteAddr().String())

					for {
						if err := conn.HandleCommand(); err != nil {
							logger.Logger.Error(err)
							break
						}
					}
				}(proxyConn)
			}
		}
	}
}
