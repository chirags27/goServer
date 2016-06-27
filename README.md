The src/ directory contains the implementation of a server and client written in the 
Go-Programming language.
The server protects a directory of files while the client requests to fetch a file
from a server. Only if the credentials sent by the client are valid, the server 
processes the request and replies with the file (if it exists)

To build the executable:


```bash
mkdir build
go build ./../src/server.go
go build ./../src/client.go
```
To run either directly type the following commands or generate executables as 
mentioned above and run them:

```bash
go run ./../src/server.go <SET_ID> <SET_PASSWORD> [DIRECTORY_TO_SEARCH]
go run ./../src/client.go <ID> <PASSWORD> <FILENAME> [DIRECTORY_TO_SAVE]
```

The default directory to search is ./../secure/
The default directory to save is ./downloaded
