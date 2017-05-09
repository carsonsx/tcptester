package cmd

import "github.com/carsonsx/tcptester/net"

var commands = make(map[string]Command)

type Command interface {
	Name() string
	Usage() string
	NeedConnected() bool
	MaxArgCount() int
	MinArgCount() int
	Validate(tcpClient *net.TCPClient, args ...string) bool
	Run(tcpClient *net.TCPClient, args ...string) error
}

func GetInstance(name string) Command {
	return commands[name]
}

func RegisterInstance(cmd Command) {
	commands[cmd.Name()] = cmd
}
