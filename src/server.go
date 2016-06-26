package main

import(
"fmt"
"net"
"os"
"os/signal"
"syscall"
"encoding/json"
"time"
"io/ioutil"
)

var port = ":4560"


type  FileRequest struct{

	Username string
	Password string
	File 	 string
}

type  FileResponse struct{

	Status 			 int
	FileContents 	 string
}

var id string = "chshah"
var pass string = "abc"
var dir string = "./../secure/"


func main(){


	ln, err := net.Listen("tcp", port)
	if err!=nil{
		fmt.Printf("Error creating a UNIX domain socket.\n")
		panic(err)
		return
	}

	sigs := make(chan os.Signal, 1)	
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func(){
		_ = <-sigs
		fmt.Printf("\nExiting Server Gracefully.\n")
		os.Exit(0)
	}()
	
	for {
		conn, err := ln.Accept()
		if err!=nil{
			fmt.Printf("Error in accpeting the connection.\n")
		}

		go serveClient(conn)
	}

	return

}


func serveClient(conn net.Conn){

	fmt.Printf("Servicing ....\n")


	d := json.NewDecoder(conn)
	var receivedReq FileRequest
	err := d.Decode(&receivedReq)
	if err!=nil{
		fmt.Printf("Error in decoding the received message.\n")
		return
	}

	sendChannel := make(chan int)
	var statusToSend int = 1
	var fileContentsString string = ""

	if receivedReq.Username == id && receivedReq.Password == pass{
		fmt.Printf("User Validated\n")
		go printRoutine(sendChannel) 
		fileContents, err := ioutil.ReadFile(dir + receivedReq.File)
		fileContentsString = string(fileContents)
		if err!=nil{
			statusToSend = 0
		}
		
	}else{
		fmt.Printf("Invalid User\n")
		statusToSend = 0
	}

	var reply FileResponse
	reply.Status = statusToSend
	reply.FileContents = fileContentsString

	replyToSend,_ := json.Marshal(&reply)
	conn.Write(replyToSend)

	sendChannel <- 1
}

func printRoutine(c chan int){

	// print ... every time when processing
	fmt.Printf("Fetching the file\n")
	for{
		select{
		case <-c:
			fmt.Printf("Processing done!\n")
			return
		default:
			fmt.Printf("..")
			time.Sleep(time.Second*1)
		}
	}
}