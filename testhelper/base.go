package testhelper

import (
	"bufio"
	"bytes"
	"net"
	"strconv"
)

const DELIM = '\n'

var operators = []string{"+", "-", "*", "/"}

type Request struct {
	ID       int64
	Operands []int32
	Operator string
}

type Response struct {
	ID      int64
	Formula string
	Result  int64
	Err     error
}

func op(operands []int32, operator string) int64 {
	var result int64
	switch {
	case operator == "+":
		for _, v := range operands {
			if result == 0 {
				result = int64(v)
			} else {
				result += int64(v)
			}
		}
	case operator == "-":
		for _, v := range operands {
			if result == 0 {
				result = int64(v)
			} else {
				result -= int64(v)
			}
		}
	case operator == "*":
		for _, v := range operands {
			if result == 0 {
				result = int64(v)
			} else {
				result *= int64(v)
			}
		}
	case operator == "/":
		for _, v := range operands {
			if result == 0 {
				result = int64(v)
			} else {
				result /= int64(v)
			}
		}
	}
	return result
}

// genFormula 会根据参数生成字符串形式的公式。
func genFormula(operands []int32, operator string, result int64, equal bool) string {
	var buff bytes.Buffer
	n := len(operands)
	for i := 0; i < n; i++ {
		if i > 0 {
			buff.WriteString(" ")
			buff.WriteString(operator)
			buff.WriteString(" ")
		}

		buff.WriteString(strconv.FormatInt(int64(operands[i]), 10))
	}
	if equal {
		buff.WriteString(" = ")
	} else {
		buff.WriteString(" != ")
	}
	buff.WriteString(strconv.FormatInt(int64(result), 10))
	return buff.String()
}

func read(conn net.Conn) ([]byte, error) {
	readByte := make([]byte, 1)
	var buffer bytes.Buffer
	for {
		_, err := conn.Read(readByte)
		if err != nil {
			return nil, err
		}
		if readByte[0] == DELIM {
			break
		} else {
			buffer.WriteByte(readByte[0])
		}
	}
	return buffer.Bytes(), nil
}

func write(conn net.Conn, writeBytes []byte) error {
	buffer := bufio.NewWriter(conn)
	_, err := buffer.Write(writeBytes)
	if err == nil {
		err = buffer.WriteByte(DELIM)
	}
	if err == nil {
		buffer.Flush()
	}
	return err
}
