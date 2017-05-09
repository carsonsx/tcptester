package cmd

import (
	"fmt"
	"github.com/carsonsx/tcptester/conf"
	"github.com/carsonsx/tcptester/net"
)

type StatusCommand struct {
}

func (c *StatusCommand) Name() string {
	return "status"
}

func (c *StatusCommand) Usage() string {
	return "usage: status"
}

func (c *StatusCommand) NeedConnected() bool {
	return false
}

func (c *StatusCommand) MaxArgCount() int {
	return 0
}

func (c *StatusCommand) MinArgCount() int {
	return 0
}

func (c *StatusCommand) Validate(tcpClient *net.TCPClient, args ...string) bool {
	return true
}

func (c *StatusCommand) Run(tcpClient *net.TCPClient, args ...string) error {
	fmt.Println(conf.RunData.NetworkStatus)
	return nil
}

func init() {
	RegisterInstance(new(StatusCommand))
}
