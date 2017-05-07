package cmd

import (
	"github.com/carsonsx/tcptester/net"
	"strings"
	"github.com/carsonsx/tcptester/cmd/input"
	"github.com/carsonsx/tcptester/conf"
)

type SendCommand struct {
}

func (c *SendCommand) Name() string {
	return "send"
}

func (c *SendCommand) Usage() string {
	return "usage: send string|0x1111"
}

func (c *SendCommand) NeedConnected() bool {
	return true
}

func (c *SendCommand) MaxArgCount() int {
	return 1
}

func (c *SendCommand) MinArgCount() int {
	return 1
}

func (c *SendCommand) Validate(tcpClient *net.TCPClient, args ...string) bool {
	return true
}

func (c *SendCommand) Run(tcpClient *net.TCPClient, args ...string) error {
	inputString := strings.Join(args, "")
	data, err := input.GetInstance(conf.Config.InputMode).Convert(inputString)
	if err == nil {
		_, err = tcpClient.Write(data)
	}
	return err
}

func init()  {
	RegisterInstance(new(SendCommand))
}

