package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var port = ":4560"

type FileRequest struct {
	Username string
	Password string
	File     string
}

type FileResponse struct {
	Status       int
	FileContents string
}

var id string = "chshah"
var pass string = "abc"
var dir string = "./../secure/"

var accessInfo = make(map[string]int)
var lockingInterval time.Duration = 3

func main() {

	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Printf("USAGE: \n./server <ID_TO_SET> <PASSWORD_TO_SET> [DIRECTORY_FOR_FILES]\n")
		fmt.Printf("Default Direcotry is ./../secure\n")
		return
	}

	id = os.Args[1]
	pass = os.Args[2]

	if len(os.Args) == 4 {
		dir = os.Args[3]
	}

	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Error creating a UNIX domain socket.\n")
		panic(err)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = <-sigs
		fmt.Printf("\nExiting Server Gracefully.\n")
		os.Exit(0)
	}()

	fmt.Printf("Remote File Server -- Listening on port 4560\n")

	var m sync.Mutex

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error in accpeting the connection.\n")
			return
		}

		go serveClient(conn, &m)
	}

	return

}

func serveClient(conn net.Conn, m *sync.Mutex) {

	fmt.Printf("Started reading the Request\n")

	d := json.NewDecoder(conn)
	var receivedReq FileRequest
	err := d.Decode(&receivedReq)
	if err != nil {
		fmt.Printf("Error in decoding the received message.\n")
		return
	}

	sendChannel := make(chan int)
	var statusToSend int = 1
	var fileContentsString string = ""

	if receivedReq.Username == id && receivedReq.Password == pass {
		fmt.Printf("User Validated\n")
		go printRoutine(sendChannel)
		// check whether the file asked is being transfered.
		// If yes, wait for the request to complete
		for {
			m.Lock()
			if accessInfo[receivedReq.File] == 0 {
				accessInfo[receivedReq.File] = 1
				break
			} else {
				m.Unlock()
				time.Sleep(lockingInterval * time.Second)
			}
		}
		m.Unlock()

		fileContents, err := ioutil.ReadFile(dir + receivedReq.File)

		m.Lock()
		accessInfo[receivedReq.File] = 0
		m.Unlock()

		fileContentsString = string(fileContents)
		if err != nil {
			statusToSend = 0
			sendChannel <- 1
		}

	} else {
		fmt.Printf("Invalid User\n")
		statusToSend = 0
	}

	var reply FileResponse
	reply.Status = statusToSend
	reply.FileContents = fileContentsString

	replyToSend, _ := json.Marshal(&reply)
	conn.Write(replyToSend)

	sendChannel <- 0
}

func printRoutine(c chan int) {

	// print ... every time when processing
	fmt.Printf("Fetching the file\n")
	for {
		select {
		case recvStatus := <-c:
			if recvStatus == 1 {
				fmt.Printf("Requested File not found\n")
			} else {
				fmt.Printf("\nProcessing done!\n")
			}
			return
		default:
			fmt.Printf(".")
			time.Sleep(time.Millisecond * 50)
		}
	}
}
