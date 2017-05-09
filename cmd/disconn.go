package cmd

import (
	"fmt"
	"github.com/carsonsx/tcptester/conf"
	"github.com/carsonsx/tcptester/net"
)

type DisconnectCommand struct {
}

func (c *DisconnectCommand) Name() string {
	return "disconnect"
}

func (c *DisconnectCommand) Usage() string {
	return "usage: disconn"
}

func (c *DisconnectCommand) NeedConnected() bool {
	return true
}

func (c *DisconnectCommand) MaxArgCount() int {
	return 0
}

func (c *DisconnectCommand) MinArgCount() int {
	return 0
}

func (c *DisconnectCommand) Validate(tcpClient *net.TCPClient, args ...string) bool {
	return true
}

func (c *DisconnectCommand) Run(tcpClient *net.TCPClient, args ...string) error {
	fmt.Printf("Thank you for using %s\n", conf.TCP_TESTER_NAME)
	return tcpClient.Close()
}

func init() {
	RegisterInstance(new(DisconnectCommand))
}
