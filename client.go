package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type Config struct {
	DialAddress            string `json:"dial_address"`
	WriteString            string `json:"write_string"`
	GoroutineNum           int    `json:"goroutine_num"`
	RequestNumPerGoroutine int    `json:"request_num_per_goroutine"`
	TimeInterval           int    `json:"time_interval"`
	TestTime               int    `json:"test_time"`
}

type SuccessCount struct {
	SuccessNum int
	mu         sync.Mutex
}

func (sc *SuccessCount) SuccessPlus() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.SuccessNum++
}

type FailCount struct {
	FailNum int
	mu      sync.Mutex
}

func (fc *FailCount) FailNumPlus() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.FailNum++
}

var config Config
var successCount = SuccessCount{}
var failCount = FailCount{}

var writeData []byte

// InitConfigFromJson 读取 config.json 并初始化信息为 GlobalConfig 对象
func InitConfigFromJson() {
	file, _ := os.Open("config.json")
	// the file needs to be closed
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("读取配置文件失败，请检查文件路径及文件内容！")
		}
	}(file)
	// create a decoder to decode json object by using the file descriptor
	decoder := json.NewDecoder(file)

	// read and construct the config struct
	err := decoder.Decode(&config)
	if err != nil {
		log.Println("初始化配置失败")
	}
}

// InitWriteData initialize writeData from the config
func InitWriteData() {
	writeData = []byte(config.WriteString)
}

func SendDataToServer(conn net.Conn, index int) {
	for {
		log.Println("goroutine ", index, " is sending data")
		_, err := conn.Write(writeData)
		if err != nil {
			// new goroutine to get the FailNumPlus function
			go failCount.FailNumPlus()
			log.Println("goroutine ", index, " sending data failed")
			time.Sleep(time.Millisecond * time.Duration(config.TimeInterval))
			continue
		}
		readData := make([]byte, 50)
		byteLength, err := conn.Read(readData)
		if err != nil {
			// new goroutine to get the FailNumPlus function
			log.Println("goroutine ", index, " sending data failed")
			go failCount.FailNumPlus()
			time.Sleep(time.Millisecond * time.Duration(config.TimeInterval))
			continue
		}
		log.Printf("收到了服务端数据:%s ", string(readData[:byteLength]))
		// new goroutine to get the SuccessPlus function
		go successCount.SuccessPlus()
		time.Sleep(time.Millisecond * time.Duration(config.TimeInterval))
	}
}

func init() {
	InitConfigFromJson()
	InitWriteData()
}

func main() {
	// start time
	startTime := time.Now()

	for i := 0; i < config.GoroutineNum; i++ {
		// initialize a tcp connection
		conn, err := net.Dial("tcp", config.DialAddress)
		if err != nil {
			log.Fatalln("TCP连接失败.....")
		}
		log.Println("连接至 ", config.DialAddress, "....")

		go SendDataToServer(conn, i)
	}

	time.Sleep(time.Minute * time.Duration(config.TestTime))
	log.Println("time cost:", time.Now().Sub(startTime).Minutes())
}
