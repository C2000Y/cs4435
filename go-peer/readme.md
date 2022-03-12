to run the program, follow command format:   
**go run file [name] [listening address] [consul address] [operation file]**  
for example:  
go run node/main.go "Gnode 1" :2333 localhost:8500 input1  
go run node/main.go "Gnode 2" :2334 localhost:8500 input2  
go run node/main.go "Gnode 3" :2335 localhost:8500 input2  
  
Or you can use the exexuctable file in directory bin  
./bin/main "Gnode 1" :2001 localhost:8500 input1  
./bin/main "Gnode 2" :2002 localhost:8500 input2  
./bin/main "Gnode 3" :2003 localhost:8500 input3  
  
Use the following command before make:  
export PATH="$PATH:$(go env GOPATH)/bin"  
  
If the node has trouble starting, try:  
consul  kv delete -recurse Gnode  
