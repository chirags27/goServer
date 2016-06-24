package main

import(
"fmt"
"net"
)

var socketFilePath = "socketFileForComm.sock"

func main(){

	conn, err := net.Dial("unix", socketFilePath)
	defer conn.Close() 
	if err!=nil{
		fmt.Printf("Error while connecting to the server.")
		return
	}
	return
}