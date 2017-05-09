package cmd

import (
	"fmt"
	"github.com/carsonsx/tcptester/net"
	"sort"
	"strings"
)

type HelpCommand struct {
}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Usage() string {
	return "usage: help\n       help cmd"
}

func (c *HelpCommand) NeedConnected() bool {
	return false
}

func (c *HelpCommand) MaxArgCount() int {
	return 1
}

func (c *HelpCommand) MinArgCount() int {
	return 0
}

func (c *HelpCommand) Validate(tcpClient *net.TCPClient, args ...string) bool {
	return true
}

func (c *HelpCommand) Run(tcpClient *net.TCPClient, args ...string) error {
	if len(args) == 0 {
		var names []string
		for name := range commands {
			names = append(names, name)
		}
		sort.Strings(names)
		fmt.Println(strings.Join(names, " "))
		return nil
	}
	name := args[0]
	_cmd := GetInstance(name)
	if _cmd != nil {
		fmt.Println(_cmd.Usage())
	}
	return nil
}

type Help2Command struct {
	*HelpCommand
}

func (c *Help2Command) Name() string {
	return "?"
}

func (c *Help2Command) Usage() string {
	return "usage: ?\n       ? cmd"
}

func init() {
	help := new(HelpCommand)
	RegisterInstance(help)
	help2 := new(Help2Command)
	help2.HelpCommand = help
	RegisterInstance(help2)
}
