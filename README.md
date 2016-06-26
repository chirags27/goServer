The src/ directory contains the implementation of a server and client written in the 
Go-Programming language.
The server protects a directory of files while the client requests to fetch a file
from a server. Only if the credentials sent by the client are valid, the server 
processes the request and replies with the file (if it exists)

To build the executable:

```bash
go build server.go
go build client.go
```