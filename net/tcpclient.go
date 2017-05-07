package net

import (
	"errors"
	"fmt"
	"log"
	"net"
	"github.com/carsonsx/tcptester/conf"
	"github.com/carsonsx/tcptester/util"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type TCPClient struct {
	Addr        string
	conn        net.Conn
	connected   bool
	CloseListen func()
	reader      Reader
	parser      Parser
	ReceivedCount int
}

func (c *TCPClient) Connect() error {

	if c.connected {
		return nil
	}

	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return err
	}
	c.conn = conn
	c.connected = true

	if !conf.Config.Silence {
		c.goReceiveData()
	}

	return nil
}


func (c *TCPClient) IsConnected() bool {
	return c.connected
}

func (c *TCPClient) Close() error {
	if !c.IsConnected() {
		return nil
	}
	err := c.conn.Close()
	if err != nil {
		return err
	}
	if c.CloseListen != nil {
		c.connected = false
		c.CloseListen()
	}
	return err
}

func (c *TCPClient) goReceiveData() {
	go func() {
		for {
			raw, data, err := c.GetReader().Read(c)
			if err != nil {
				c.Close()
				break
			}
			if len(raw) == 0 {
				c.Close()
				break
			}
			c.ReceivedCount++
			//fmt.Printf("%v%d\ns", data, c.ReceivedCount)
			c.GetParser().Unmarshal(data)
		}
	}()
}


func (c *TCPClient) WriteData(v interface{}) ([]byte, error) {
	data, err := c.parser.Marshal(v)
	if err != nil {
		return nil, err
	}
	return c.Write(data)
}

func (c *TCPClient) Write(data []byte) ([]byte, error) {
	dataLen := len(data)
	if conf.Config.MessageLengthSize == 1 {
		data = util.AddUint8ToBytePrefix(data, uint8(dataLen))
	} else if conf.Config.MessageLengthSize == 2 {
		data = util.AddUint16ToBytePrefix(data, uint16(dataLen), conf.Config.LittleEndian)
	} else if conf.Config.MessageLengthSize == 4 {
		data = util.AddUint32ToBytePrefix(data, uint32(dataLen), conf.Config.LittleEndian)
	}
	//for _, b := range data {
	//	fmt.Printf("%d", b)
	//}
	//fmt.Println()
	n, err := c.conn.Write(data)
	if err != nil {
		c.Close()
		return nil, err
	} else if n != len(data) {
		c.Close()
		return nil, errors.New(fmt.Sprintf("sent failed: expected length %d, actual length %d \n", len(data), n))
	}

	if !conf.Config.Silence && conf.Config.SyncTime > 0 {
		time.Sleep(time.Duration(conf.Config.SyncTime) * time.Second)
	}

	return data, nil
}

func (c *TCPClient) SetReader(reader Reader) {
	c.reader = reader
}

func (c *TCPClient) GetReader() Reader {
	if c.reader == nil {
		if conf.Config.Reader == "line" {
			c.reader = new(LineReader)
		} else if conf.Config.Reader == "protobuf" {
			c.reader = new(ProtobufReader)
		}
	}
	return c.reader
}

func (c *TCPClient) SetParser(parser Parser) {
	c.parser = parser
}

func (c *TCPClient) GetParser() Parser {
	if c.parser == nil {
		if conf.Config.Parser == "string" {
			c.parser = new(StringParser)
		} else if conf.Config.Parser == "protobuf" {
			c.parser = new(ProtobufParser)
		}
	}
	return c.parser
}
