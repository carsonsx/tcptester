package cmd

import (
	"fmt"
	"github.com/carsonsx/tcptester/conf"
	"github.com/carsonsx/tcptester/net"
)

type ConnectCommand struct {
}

func (c *ConnectCommand) Name() string {
	return "connect"
}

func (c *ConnectCommand) Usage() string {
	return "usage: connect 127.0.0.1:8888"
}

func (c *ConnectCommand) NeedConnected() bool {
	return false
}

func (c *ConnectCommand) MaxArgCount() int {
	return 1
}

func (c *ConnectCommand) MinArgCount() int {
	return 1
}

func (c *ConnectCommand) Validate(tcpClient *net.TCPClient, args ...string) bool {
	return true
}

func (c *ConnectCommand) Run(tcpClient *net.TCPClient, args ...string) error {
	tcpClient.Addr = args[0]
	fmt.Println("Trying " + tcpClient.Addr)
	err := tcpClient.Connect()
	if err != nil {
		return err
	}
	conf.RunData.NetworkStatus = "Connected to " + tcpClient.Addr
	fmt.Println(conf.RunData.NetworkStatus)
	conf.RunData.ConsolePrompt = tcpClient.Addr
	return nil
}

func init() {
	RegisterInstance(new(ConnectCommand))
}
