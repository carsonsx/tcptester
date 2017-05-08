package cmd

import (
	"fmt"
	"github.com/carsonsx/tcptester/cmd/input"
	"github.com/carsonsx/tcptester/conf"
	"github.com/carsonsx/tcptester/net"
	"strconv"
	"sync"
	"time"
)

type LoadCommand struct {
	connectedClientCount int
	mutex                sync.Mutex
	stop                 bool
}

type LoadData struct {
	BeginData  []*TestData
	StressData []*TestData
	EndData    []*TestData
}

type TestData struct {
	GlobalData []byte
	ClientData [][]byte
	SleepTime int
}

type loadStat struct {
	id       int
	sent     int
	success  int
	fail     int
	received int
}

func (s *loadStat) print() {
	if s.id == 0 {
		fmt.Printf("total send %d, success %d, fail %d, received %d\n", s.sent, s.success, s.fail, s.received)
	} else {
		fmt.Printf("client[%d] send %d, success %d, fail %d, received %d\n", s.id, s.sent, s.success, s.fail, s.received)
	}
}

func (s *loadStat) add(another *loadStat) {
	s.sent += another.sent
	s.success += another.success
	s.fail += another.fail
	s.received += another.received
}

func (s *loadStat) count(err error) {
	s.sent++
	if err == nil {
		s.success++
	} else {
		s.fail++
	}
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

	inputTestData := args[0]
	clientCount, _ := strconv.Atoi(args[1])
	sendCount, _ := strconv.Atoi(args[2])
	sendInterval := 0
	if len(args) >= 4 {
		sendInterval, _ = strconv.Atoi(args[3])
	}
	data, err := input.GetInstance(conf.Config.InputMode).Convert(inputTestData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	testData := new(TestData)
	testData.GlobalData = data
	testData.SleepTime = sendInterval
	loadData := new(LoadData)
	loadData.StressData = []*TestData{testData}
	return c.Load(tcpClient, loadData, clientCount, sendCount)
}

func (c *LoadCommand) Load(tcpClient *net.TCPClient, loadData *LoadData, clientCount int, sendCount int) error {
	//for _, d := range loadData {
	//	fmt.Printf("%d ", d)
	//}
	//fmt.Println()

	done := make(chan *loadStat)

	//ignore syncTime
	syncTime := conf.Config.SyncTime
	conf.Config.SyncTime = 0

	for i := 1; i <= clientCount; i++ {
		go c.load(tcpClient, i, loadData, sendCount, done)
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

func (c *LoadCommand) load(tcpClient *net.TCPClient, clientId int, loadData *LoadData, sendCount int, done chan *loadStat) {

	// create new tcp client
	newTcpClient := new(net.TCPClient)
	defer newTcpClient.Close()
	newTcpClient.Addr = tcpClient.Addr
	newTcpClient.SetParser(tcpClient.GetParser())
	newTcpClient.SetReader(tcpClient.GetReader())
	err := newTcpClient.Connect()

	stat := new(loadStat)
	stat.id = clientId

	if err != nil {
		fmt.Printf("client[%d] failed to connect %s: %v\n", clientId, newTcpClient.Addr, err)
	} else {
		c.mutex.Lock()
		c.connectedClientCount++
		c.mutex.Unlock()
		fmt.Printf("client[%d] established\n", clientId)

		c.SendTestDataGroup(newTcpClient, clientId, loadData.BeginData, stat)

		for j := 1; !c.stop && j <= sendCount; j++ {
			c.SendTestDataGroup(newTcpClient, clientId, loadData.StressData, stat)
		}

		c.SendTestDataGroup(newTcpClient, clientId, loadData.EndData, stat)

		//sleep 5 seconds for loadData receiving
		time.Sleep(5 * time.Second)
		stat.received = newTcpClient.ReceivedCount

		stat.print()

	}

	done <- stat
}


func (c *LoadCommand) SendTestDataGroup(tcpClient *net.TCPClient, clientId int, testDataGroup []*TestData, stat *loadStat)  {

	if testDataGroup == nil || len(testDataGroup) == 0 {
		return
	}

	var err error
	for _, testData := range testDataGroup {
		err = c.SendTestData(tcpClient, clientId, testData)
		stat.count(err)
	}
}

func (c *LoadCommand) SendTestData(tcpClient *net.TCPClient, clientId int, testData *TestData) error {
	var err error
	if testData.ClientData != nil {
		_, err = tcpClient.Write(testData.ClientData[clientId-1])
	} else {
		_, err = tcpClient.Write(testData.GlobalData)
	}
	if err != nil {
		fmt.Printf("client[%d] send data error: %v", clientId, err)
	}
	if testData.SleepTime > 0 {
		time.Sleep(time.Millisecond * time.Duration(testData.SleepTime))
	}
	return err
}

func init() {
	RegisterInstance(new(LoadCommand))
}
