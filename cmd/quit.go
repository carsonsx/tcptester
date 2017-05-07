package cmd

import (
	"fmt"
	"os"
	"github.com/carsonsx/tcptester/net"
)

type QuitCommand struct {
}

func (c *QuitCommand) Name() string {
	return "quit"
}

func (c *QuitCommand) Usage() string {
	return "usage: quit"
}

func (c *QuitCommand) NeedConnected() bool {
	return false
}

func (c *QuitCommand) MaxArgCount() int {
	return 0
}

func (c *QuitCommand) MinArgCount() int {
	return 0
}

func (c *QuitCommand) Validate(tcpClient *net.TCPClient, args ...string) bool {
	return true
}

func (c *QuitCommand) Run(tcpClient *net.TCPClient, args ...string) error {
	fmt.Println("bye")
	os.Exit(0)
	return nil
}

func init()  {
	RegisterInstance(new(QuitCommand))
}

