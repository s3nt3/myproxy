package server

import (
	"fmt"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/server"
	"github.com/s3nt3/myproxy/internal/logger"
	"github.com/s3nt3/myproxy/internal/utils"
)

type MockProxyServerConfig struct {
	Protocol *server.Server
}

func NewMockProxyServerConfig() ProxyServerConfig {
	return &MockProxyServerConfig{
		Protocol: server.NewDefaultServer(),
	}
}

func (mock *MockProxyServerConfig) GetProtocolConfig() *server.Server {
	return mock.Protocol
}

type MockHandler struct {
	db *client.Conn
}

func NewMockHandler(upstream string, user string, password string) *MockHandler {
	db, err := client.Connect(upstream, user, password, "")
	utils.FatalCheck(err, fmt.Sprintf("Connect to upstream %s", upstream))
	return &MockHandler{
		db: db,
	}
}

func (h *MockHandler) UseDB(dbName string) error {
	logger.Logger.Debugf("UseDB: USE %s", dbName)
	return h.db.UseDB(dbName)
}

func (h *MockHandler) HandleQuery(query string) (*mysql.Result, error) {
	logger.Logger.Debugf("HandleQuery: %s", query)
	return h.db.Execute(query)
}

func (h *MockHandler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	logger.Logger.Debugf("HandleFieldList: %s %s", table, fieldWildcard)
	return h.db.FieldList(table, fieldWildcard)
}

func (h *MockHandler) HandleStmtPrepare(query string) (int, int, interface{}, error) {
	logger.Logger.Debugf("HandleStmtPrepare: %s", query)
	return 0, 0, nil, fmt.Errorf("not supported now")
}

func (h *MockHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	logger.Logger.Debugf("HandleStmtExecute: %s, args: %+v", query, args)
	return h.db.Execute(query, args...)
}

func (h *MockHandler) HandleStmtClose(context interface{}) error {
	logger.Logger.Debugf("HandleStmtClose")
	return h.db.Close()
}

func (h *MockHandler) HandleOtherCommand(cmd byte, data []byte) error {
	logger.Logger.Debugf("HandleOtherCommand: cmd(%x), data: %+v", cmd, data)
	return mysql.NewError(
		mysql.ER_UNKNOWN_ERROR,
		fmt.Sprintf("command %d is not supported now", cmd),
	)
}

type MockProxySeverHandler struct {
	Host string
	Port uint
}

func NewMockProxyServerHandler(host string, port uint) ProxyServerHandler {
	return &MockProxySeverHandler{
		Host: host, Port: port,
	}
}

func (mock *MockProxySeverHandler) GetProxyHandler() server.Handler {
	return NewMockHandler(fmt.Sprintf("%s:%d", mock.Host, mock.Port), "root", "")
}
