package main

import "github.com/carsonsx/tcptester/net"
import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/carsonsx/tcptester/cmd"
	"github.com/carsonsx/tcptester/conf"
)

var tcpClient = new (net.TCPClient)

func main() {

	flag.BoolVar(&conf.Config.Silence, "silence", conf.Config.Silence, "enable or disable receive data")
	flag.IntVar(&conf.Config.MessageLengthSize, "message_length_size", conf.Config.MessageLengthSize, "set message length size, if 0 then add nothing for message head")
	flag.IntVar(&conf.Config.ProtoBufferIdSize, "proto_buffer_id_size", conf.Config.ProtoBufferIdSize,"set protobuf id size")
	flag.BoolVar(&conf.Config.LittleEndian, "little_endian", conf.Config.LittleEndian, "use little endian for integer encode")
	flag.IntVar(&conf.Config.SyncTime, "sync_time", 1,"time for receiving data one by one")
	flag.StringVar(&conf.Config.Reader, "reader", conf.Config.Reader, "set reader with line|protobuf")
	flag.StringVar(&conf.Config.Parser, "parser", conf.Config.Parser, "set parser with string|protobuf")
	flag.StringVar(&conf.Config.InputMode, "input_mode", conf.Config.Parser, "set input mode string|hex")
	flag.Parse()

	conf.CheckAndSetConfig()

	fmt.Printf("Welcome to use %s.\n", conf.TCP_TESTER_NAME)
	fmt.Println("You may get source code from 'https://github.com/carsonsx/tcptester'")
	new(cmd.ConfigCommand).Run(tcpClient)
	fmt.Println("Please type '?' for command list or '? $cmd' for specific command usage")

	tcpClient.CloseListen = func() {
		fmt.Println("\ndisconnected from " + tcpClient.Addr)
		setDisconnectData()
		printPrompt()
	}

	tcpClient.GetParser().SetHandler(func(v interface{}) {
		fmt.Println(v)
	})

	if flag.NArg() > 0 {
		err := new(cmd.ConnectCommand).Run(tcpClient, flag.Arg(flag.NArg() - 1))
		if err != nil {
			os.Exit(1)
			return
		}
	} else {
		setDisconnectData()
	}
	for {
		input()
	}
	tcpClient.Close()
}

func setDisconnectData()  {
	conf.RunData.NetworkStatus = conf.TCP_TESTER_NOT_CONNECTED
	conf.RunData.ConsolePrompt = conf.TCP_TESTER_NAME
}

func printPrompt()  {
	fmt.Print(conf.RunData.ConsolePrompt + "> ")
}

func input() {
	printPrompt()
	reader := bufio.NewReader(os.Stdin)
	line, _, _ := reader.ReadLine()
	text := string(line)
	if text == "" {
		new(cmd.HelpCommand).Run(tcpClient)
		return
	}
	fields := strings.Fields(text)
	name := fields[0]
	command := cmd.GetInstance(name)
	if command == nil {
		fmt.Println(`(error) ERR unknown command '` + name + `'`)
		return
	}
	if command.NeedConnected() && !tcpClient.IsConnected() {
		conf.RunData.NetworkStatus = conf.TCP_TESTER_NOT_CONNECTED
		fmt.Println(conf.RunData.NetworkStatus)
		return
	}
	args := fields[1:]
	argLen := len(fields) - 1
	if argLen > command.MaxArgCount() || argLen < command.MinArgCount() || !command.Validate(tcpClient, args...) {
		fmt.Println("(error) ERR wrong arguments")
		fmt.Println(command.Usage())
		return
	}
	err := command.Run(tcpClient, args...)
	if err != nil {
		fmt.Println(err)
	}
}
