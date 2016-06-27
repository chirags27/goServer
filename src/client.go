package main

import(
"fmt"
"net"
"os"
"encoding/json"
// "time"
)

var tcpAddr string = "127.0.0.1:4560"

var targetDirectory string = "downloaded"


type  FileRequest struct{

	Username string
	Password string
	File 	 string
}

type  FileResponse struct{

	Status 			 int
	FileContents 	 string
}




func main(){

	// Get the credentails and files requested
	
	if len(os.Args) !=4 && len(os.Args) !=5{
		fmt.Printf("USAGE: ./client <ID> <PASSWORD> <FILE> <TARGET_DIRECTORY>\n")
		fmt.Printf("Default Target Direcory: downloaded/ \n")
		return
	}
	if len(os.Args) == 5{
		targetDirectory = os.Args[4]
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

	// fmt.Printf("Sent: %s -- Waiting for Reply\n", string(jsonReq))

	d := json.NewDecoder(conn)
	var resp FileResponse
	err = d.Decode(&resp)
	if err!=nil{
		fmt.Printf("Error in decoding the received message.\n")
		return
	}

	if resp.Status != 1{
		fmt.Printf("Invalid ID, Password or Filename. Try again.\n")
		return
	}

	printChannel := make(chan int)
	go printRoutine(printChannel)

	printChannel<-1
	// Create a file in the local directory and write the contents to it
	err = os.Mkdir(targetDirectory, 0777)
	f,err := os.Create(targetDirectory + "/" + reqToSend.File)
	err = os.Chmod(targetDirectory + "/" + reqToSend.File, 0666)

	if err!=nil{
		panic(err)
	}

	_, err = f.Write([]byte(resp.FileContents))

	if err!=nil{
		panic(err)
	}
	printChannel<-1
	<-printChannel
}


func printRoutine(printChannel chan int){
	<-printChannel
	fmt.Printf("Started writing the received contents!\n")
	<-printChannel
	fmt.Printf("Done writing the received contents!\n")
	printChannel<-1

}