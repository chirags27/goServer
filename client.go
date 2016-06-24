package main

import(
"fmt"
"net"
"os"
"encoding/json"
// "time"
)

var tcpAddr = "localhost:4560"


type  FileRequest struct{

	Username string
	Password string
	File 	 string
}


func main(){

	// Get the credentails and files requested
	
	if len(os.Args) !=4 {
		fmt.Printf("USAGE: ./client ID PASSWORD FILE\n")
		return
	}




	conn, err := net.Dial("tcp", tcpAddr)

	if err!=nil{
		fmt.Printf("Error while connecting to the server.\n")
		return
	}

	defer conn.Close()

	// send the request to the client and wait and read the files

	var reqToSend FileRequest
	reqToSend.Username = os.Args[1]
	reqToSend.Password = os.Args[2]
	reqToSend.File = os.Args[3]
	fmt.Printf("File Requested: %s\n", reqToSend.File)

	jsonReq, err := json.Marshal(&reqToSend)
	if err!=nil{
		fmt.Printf("Error in encoding struct to json")
		return
	}



	_,err = conn.Write(jsonReq)

	if err!=nil{
		fmt.Printf("Error in sending json message to server")
		return
	}

	fmt.Printf("Sent: %s\n", string(jsonReq))
	
	return
}