package main

import(
"fmt"
"net"
"os"
"os/signal"
"syscall"
"io"
)

var port = ":4560"


type  FileRequest struct{

	Username string
	Password string
	File 	 string
}


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
		fmt.Printf("Exiting Server Gracefully.\n")
		os.Exit(0)
	}()

	numServed :=0
	sendChannel := make(chan int)
	
	go printRoutine(sendChannel)

	for {
		conn, err := ln.Accept()
		if err!=nil{
			fmt.Printf("Error in accpeting the connection.\n")
		}
		numServed = numServed + 1
		go serveClient(conn, numServed, sendChannel)
	}

	return

}


func serveClient(conn net.Conn, numServed int, c chan int){

	fmt.Printf("Servicing ....\n")
	var receivedReq []byte
	_,err := conn.Read(receivedReq)

	if err!=nil && err!=io.EOF{
		fmt.Printf("Error occured while reading. Couldn't serve. Returning.\n")
		panic(err)
		return
	}

	fmt.Printf("Received: %s\n", string(receivedReq))
	c <- numServed
}

func printRoutine(c chan int){

	fmt.Printf("Waiting to serve clients.\n")
	for {
	num := <-c
	fmt.Printf("Total Clients Served: %v\n", num)
}


}