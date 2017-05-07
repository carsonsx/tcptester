package cmd

import (
	"fmt"
	"github.com/carsonsx/tcptester/net"
	"github.com/carsonsx/tcptester/conf"
	"encoding/json"
	"strconv"
)

type ConfigCommand struct {
}

func (c *ConfigCommand) Name() string {
	return "config"
}

func (c *ConfigCommand) Usage() string {
	return "usage: config"
}

func (c *ConfigCommand) NeedConnected() bool {
	return true
}

func (c *ConfigCommand) MaxArgCount() int {
	return 2
}

func (c *ConfigCommand) MinArgCount() int {
	return 0
}

func (c *ConfigCommand) Validate(tcpClient *net.TCPClient, args ...string) bool {
	return true
}

func (c *ConfigCommand) Run(tcpClient *net.TCPClient, args ...string) error {

	argLen := len(args)

	if argLen == 0 {
		bytes, _ := json.Marshal(conf.Config)
		fmt.Println("Config -> " + string(bytes))
	} else if argLen == 1 {
		key := args[0]
		switch key {
		case "silence":
			fmt.Printf("%s=%v\n", key, conf.Config.Silence)
		case "message_length_size":
			fmt.Printf("%s=%d\n", key, conf.Config.MessageLengthSize)
		case "proto_buffer_id_size":
			fmt.Printf("%s=%d\n", key, conf.Config.ProtoBufferIdSize)
		case "little_endian":
			fmt.Printf("%s=%v\n", key, conf.Config.LittleEndian)
		case "sync_time":
			fmt.Printf("%s=%d\n", key, conf.Config.SyncTime)
		case "reader":
			fmt.Printf("%s=%s\n", key, conf.Config.Reader)
		case "parser":
			fmt.Printf("%s=%s\n", key, conf.Config.Parser)
		case "input_mode":
			fmt.Printf("%s=%s\n", key, conf.Config.InputMode)
		}
	} else if argLen == 2 {
		key := args[0]
		value := args[1]
		switch key {
		case "silence":
			if silence, err := strconv.ParseBool(value); err == nil {
				conf.Config.Silence = silence
			} else {
				fmt.Println(err)
			}
		case "message_length_size":
			if message_length_size, err := strconv.ParseInt(value, 10, 0); err == nil {
				conf.CheckAndSetMessageLengthSize(int(message_length_size), conf.Config.MessageLengthSize)
			} else {
				fmt.Println(err)
			}
		case "proto_buffer_id_size":
			if proto_buffer_id_size, err := strconv.ParseInt(value, 10, 0); err == nil {
				conf.CheckAndSetProtoBufferIdSize(int(proto_buffer_id_size), conf.Config.ProtoBufferIdSize)
			} else {
				fmt.Println(err)
			}
		case "little_endian":
			if little_endian, err := strconv.ParseBool(value); err == nil {
				conf.Config.LittleEndian = little_endian
			} else {
				fmt.Println(err)
			}
		case "sync_time":
			if sync_time, err := strconv.ParseInt(value, 10, 0); err == nil {
				conf.Config.SyncTime = int(sync_time)
			} else {
				fmt.Println(err)
			}
		case "reader":
			conf.Config.Reader = value
		case "parser":
			conf.Config.Parser = value
		case "input_mode":
			conf.Config.InputMode = value
		}
	}

	return nil
}

func init()  {
	RegisterInstance(new(ConfigCommand))
}

