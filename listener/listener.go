package listener

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/listener/notify"
	"log"
	"os"
)

var (
	ErrPayloadLength = errors.New("Header中len长度与实际读取长度不一致")
)
var (
	handled_events = map[string]int{ "PROCESS_STATE_EXITED": 1, "PROCESS_STATE_FATAL": 1, "PROCESS_LOG_STDERR": 1, "PROCESS_LOG_STDOUT": 1}
)

func Start() {
	defer func() {
		if err := recover(); err != nil {
			log.Print("panic", err)
		}
	}()
	log.Print("Event Listener Start!!!")
	listen()
}

// 监听事件, 从标准输入获取事件内容
func listen() {
	reader := bufio.NewReader(os.Stdin)
	for {
		ready()
		header, err := readHeader(reader)
		if err != nil {
			failure(err)
			continue
		}
		payload, err := readPayload(reader, header.Len)
		if err != nil {
			failure(err)
			continue
		}
		log.Print("Got Event!", header)
		if _,ok := handled_events[header.EventName]; ok {
			notify.Push(header, payload)
		}
		/*
		// 只处理进程异常退出事件
		if header.EventName == "PROCESS_STATE_EXITED" {
			notify.Push(header, payload)
		}
		*/
		success()
	}
}

// 读取header
func readHeader(reader *bufio.Reader) (*event.Header, error) {
	// 读取Header
	data, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	//log.Print("HEADER:", data)
	// 解析Header
	header, err := event.ParseHeader(data)
	if err != nil {
		return nil, err
	}

	return header, nil
}

// 读取payload
func readPayload(reader *bufio.Reader, payloadLen int) (*event.Payload, error) {
	// 读取payload
	buf := make([]byte, payloadLen)
	length, err := reader.Read(buf)
	if err != nil {
		return nil, err
	}
	if payloadLen != length {
		return nil, ErrPayloadLength
	}
	//log.Print("PAYLOAD:", string(buf))
	// 解析payload
	payload, err := event.ParsePayload(string(buf))
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func ready() {
	fmt.Fprint(os.Stdout, "READY\n")
}

func success() {
	fmt.Fprint(os.Stdout, "RESULT 2\nOK")
}

func failure(err error) {
	fmt.Fprintln(os.Stderr, err)
	fmt.Fprint(os.Stdout, "Result 2\nFAIL")
}
