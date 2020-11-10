package redis

import (
	"bytes"
	"fmt"
	"log"
)

// NewLoggingConn returns a logging wrapper around a connection.
func NewLoggingConn(conn Conn, logger *log.Logger, prefix string) Conn {
	if prefix != "" {
		prefix = prefix + "."
	}
	return &loggingConn{conn, logger, prefix}
}

type loggingConn struct {
	Conn
	logger *log.Logger
	prefix string
}

func (c *loggingConn) Close() error {
	err := c.Conn.Close()
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%sClose() -> (%v)", c.prefix, err)
	c.logger.Output(2, buf.String())
	return err
}

func (c *loggingConn) printValue(buf *bytes.Buffer, v interface{}) {
	const chop = 32
	switch v := v.(type) {
	case []byte:
		if len(v) > chop {
			fmt.Fprintf(buf, "%q...", v[:chop])
		} else {
			fmt.Fprintf(buf, "%q", v)
		}
	case string:
		if len(v) > chop {
			fmt.Fprintf(buf, "%q...", v[:chop])
		} else {
			fmt.Fprintf(buf, "%q", v)
		}
	case []interface{}:
		if len(v) == 0 {
			buf.WriteString("[]")
		} else {
			sep := "["
			fin := "]"
			if len(v) > chop {
				v = v[:chop]
				fin = "...]"
			}
			for _, vv := range v {
				buf.WriteString(sep)
				c.printValue(buf, vv)
				sep = ", "
			}
			buf.WriteString(fin)
		}
	default:
		fmt.Fprint(buf, v)
	}
}

func (c *loggingConn) print(method, commandName string, args []interface{}, reply interface{}, err error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s%s(", c.prefix, method)
	if method != "Receive" {
		buf.WriteString(commandName)
		for _, arg := range args {
			buf.WriteString(", ")
			c.printValue(&buf, arg)
		}
	}
	buf.WriteString(") -> (")
	if method != "Send" {
		c.printValue(&buf, reply)
		buf.WriteString(", ")
	}
	fmt.Fprintf(&buf, "%v)", err)
	c.logger.Output(3, buf.String())
}

func (c *loggingConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	reply, err := c.Conn.Do(commandName, args...)
	c.print("Do", commandName, args, reply, err)
	return reply, err
}

func (c *loggingConn) Send(commandName string, args ...interface{}) error {
	err := c.Conn.Send(commandName, args...)
	c.print("Send", commandName, args, nil, err)
	return err
}

func (c *loggingConn) Receive() (interface{}, error) {
	reply, err := c.Conn.Receive()
	c.print("Receive", "", nil, reply, err)
	return reply, err
}
