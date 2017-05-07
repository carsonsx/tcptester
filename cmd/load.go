package cmd

import (
	"fmt"
	"github.com/carsonsx/tcptester/net"
	"strconv"
	"sync"
	"time"
	"github.com/carsonsx/tcptester/conf"
	"github.com/carsonsx/tcptester/cmd/input"
)

type LoadCommand struct {
	connectedClientCount int
	mutex sync.Mutex
	stop bool
}

type loadStat struct {
	id int
	sent int
	success int
	fail int
	received int
}

func (s *loadStat) print()  {
	if s.id == 0 {
		fmt.Printf("total send %d, success %d, fail %d, received %d\n", s.sent, s.success, s.fail, s.received)
	} else {
		fmt.Printf("client[%d] send %d, success %d, fail %d, received %d\n", s.id, s.sent, s.success, s.fail, s.received)
	}
}


func (s *loadStat) add(another loadStat)  {
	s.sent += another.sent
	s.success += another.success
	s.fail += another.fail
	s.received += another.received
}

func (c *LoadCommand) Name() string {
	return "load"
}

func (c *LoadCommand) Usage() string {
	return "usage: load test_data client_count send_count send_interval"
}

func (c *LoadCommand) NeedConnected() bool {
	return true
}

func (c *LoadCommand) MaxArgCount() int {
	return 4
}

func (c *LoadCommand) MinArgCount() int {
	return 1
}

func (c *LoadCommand) Validate(tcpClient *net.TCPClient, args ...string) bool {
	return true
}

func (c *LoadCommand) Run(tcpClient *net.TCPClient, args ...string) error {

	testData := args[0]
	clientCount, _ := strconv.Atoi(args[1])
	sendCount, _ := strconv.Atoi(args[2])
	sendInterval := 0
	if len(args) >= 4 {
		sendInterval, _ = strconv.Atoi(args[3])
	}
	data, err := input.GetInstance(conf.Config.InputMode).Convert(testData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//for _, d := range data {
	//	fmt.Printf("%d ", d)
	//}
	//fmt.Println()

	done := make(chan loadStat)

	//ignore syncTime
	syncTime := conf.Config.SyncTime
	conf.Config.SyncTime = 0

	for i := 1; i <= clientCount; i++ {
		go c.start(tcpClient, i, data, sendCount, sendInterval, done)
	}

	doneClientCount := 0
	var total loadStat
	for {
		total.add(<-done)
		doneClientCount++
		if doneClientCount >= clientCount {
			fmt.Println("all clients completed")
			break
		}
	}

	total.print()

	conf.Config.SyncTime = syncTime

	return nil
}


func (c *LoadCommand) start(tcpClient *net.TCPClient, clientId int, data []byte, sendCount int, sendInterval int, done chan loadStat) {
	newTcpClient := new(net.TCPClient)
	defer newTcpClient.Close()
	newTcpClient.Addr = tcpClient.Addr
	newTcpClient.SetParser(tcpClient.GetParser())
	newTcpClient.SetReader(tcpClient.GetReader())
	err := newTcpClient.Connect()

	var stat loadStat
	stat.id = clientId
	stat.sent = sendCount

	if err != nil {
		fmt.Printf("client[%d] failed to connect %s: %v\n", clientId, newTcpClient.Addr, err)
	} else {
		c.mutex.Lock()
		c.connectedClientCount++
		c.mutex.Unlock()
		fmt.Printf("client[%d] established\n", clientId)
		for j := 1; !c.stop && j <= sendCount; j++ {
			_, err = newTcpClient.Write(data)
			if err != nil {
				fmt.Printf("client[%d] write data error: %v", clientId, err)
				stat.fail++
			} else {
				//fmt.Printf("client[%d] write data done\n", clientId)
				stat.success++
				if sendInterval > 0 {
					//fmt.Printf("client[%d] sleep %d milliseconds\n", clientId, sendInterval)
					time.Sleep(time.Millisecond * time.Duration(sendInterval))
				}
			}
		}

		//sleep 5 seconds for data receiving
		time.Sleep(5 * time.Second)
		stat.received = newTcpClient.ReceivedCount

		stat.print()

	}


	done <- stat
}

func init() {
	RegisterInstance(new(LoadCommand))
}
