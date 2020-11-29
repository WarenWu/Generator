package testhelper

import (
	"encoding/json"
	"errors"
	"net"
	"sync/atomic"
)

type TcpServer struct {
	listener net.Listener
	active   uint32
}

func NewTcpServer() *TcpServer {
	return &TcpServer{}
}

func (server *TcpServer) init(addr string) error {
	if !atomic.CompareAndSwapUint32(&server.active, 0, 1) {
		return nil
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		atomic.StoreUint32(&server.active, 0)
		return err
	}
	server.listener = listener
	return nil
}

func (server *TcpServer) Listen(addr string) error {
	err := server.init(addr)
	if err != nil {
		return err
	}
	go func() {
		for {
			if atomic.LoadUint32(&server.active) != 1 {
				break
			}
			conn, err := server.listener.Accept()
			if err != nil {
				//打印日志
				continue
			}
			go handleFunc(conn)
		}
	}()
	return nil
}

func handleFunc(conn net.Conn) {
	var resp Response
	var req Request
	var errMsg string
	readBytes, err := read(conn)
	if err != nil {
		errMsg += "TcpServer read error :" + err.Error()
	} else {
		err := json.Unmarshal(readBytes, req)
		if err != nil {
			errMsg += "TcpServer read error :" + err.Error()
		} else {
			resp.ID = req.ID
			resp.Result = op(req.Operands, req.Operator)
			resp.Formula = genFormula(req.Operands, req.Operator, resp.Result, true)
		}
	}
	if errMsg != "" {
		resp.Err = errors.New(errMsg)
	}
	writeBytes, err := json.Marshal(resp)
	if err != nil {
		//
	}

	err = write(conn, writeBytes)
	if err != nil {

	}
}

func (server *TcpServer) Close() {
	atomic.StoreUint32(&server.active, 0)
	server.listener.Close()
}
