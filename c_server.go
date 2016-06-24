package main

import(
"fmt"
"net"
"os"
"os/signal"
"syscall"
)

var socketFilePath = "socketFileForComm.sock"

func main(){
	ln, err := net.Listen("unix", socketFilePath)
	if err!=nil{
		fmt.Printf("Error creating a UNIX domain socket.\n")
		panic(err)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func(){
		_ = <-sigs
		os.Remove(socketFilePath)
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
		go serveClient(&conn, numServed, sendChannel)
	}

	return

}


func serveClient(conn *net.Conn, numServed int, c chan int){

	fmt.Printf("Servicing ....\n")
	c <- numServed
}

func printRoutine(c chan int){

	fmt.Printf("Waiting to serve clients.\n")
	for {
	num := <-c
	fmt.Printf("Total Clients Served: %v\n", num)
}


}