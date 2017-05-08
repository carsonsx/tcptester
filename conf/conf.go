package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/carsonsx/tcptester/cmd/input"
)

const (
	TCP_TESTER_NAME        = "tcptester"
	TCP_TESTER_CONFIG_FILE = TCP_TESTER_NAME + ".conf"
)

type config struct {
	Silence           bool   `json:"silence"`
	MessageLengthSize int    `json:"message_length_size"`
	ProtoBufferIdSize int    `json:"proto_buffer_id_size"`
	LittleEndian      bool   `json:"little_endian"`
	SyncTime          int    `json:"sync_time"`
	Reader            string `json:"reader"`
	Parser            string `json:"parser"`
	InputMode         string  `json:"input_mode"`
}

var Config config

func init() {
	initDefaultConfig()
	mergeFileConfig()
}

func initDefaultConfig() {
	Config.Silence = false
	Config.MessageLengthSize = 0
	Config.ProtoBufferIdSize = 2
	Config.LittleEndian = false
	Config.SyncTime = 0
	Config.Reader = "line"
	Config.Parser = "string"
	Config.InputMode = new(input.StringConverter).Name()
}

func mergeFileConfig() {

	data, err := ioutil.ReadFile(TCP_TESTER_CONFIG_FILE)
	if err != nil {
		//fmt.Printf("Read %s failed, use default config\n", TCP_TESTER_CONFIG_FILE)
		return
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		fmt.Printf("Unmarshal %s failed, use default config", TCP_TESTER_CONFIG_FILE)
		return
	}

	CheckAndSetConfig()
}

func CheckAndSetConfig() {
	CheckAndSetMessageLengthSize(Config.MessageLengthSize, 0)
	CheckAndSetProtoBufferIdSize(Config.ProtoBufferIdSize, 2)
}

func CheckAndSetMessageLengthSize(expectedSize, oldSize int)  {
	if expectedSize != 0 &&
		expectedSize != 1 &&
		expectedSize != 2 &&
		expectedSize != 4 {
		output := "Invalid message length size %d, must be 0, 1, 2 or 4, keep value %d\n"
		fmt.Printf(output, expectedSize, oldSize)
		Config.MessageLengthSize = oldSize
	} else {
		Config.MessageLengthSize = expectedSize
	}
}


func CheckAndSetProtoBufferIdSize(expectedSize, oldSize int)  {
	if expectedSize != 1 &&
		expectedSize != 2 &&
		expectedSize != 4 {
		output := "Invalid protobuf id size %d, must be 1, 2 or 4, keep value %d\n"
		fmt.Printf(output, expectedSize, oldSize)
		Config.ProtoBufferIdSize = oldSize
	} else {
		Config.MessageLengthSize = expectedSize
	}
}