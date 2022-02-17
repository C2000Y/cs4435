package main

import (
    "context"
    "log"
    "time"
    "io/ioutil"

    "google.golang.org/grpc"

    pb "cs4435/calculator/proto"
)

const (
    address     = "localhost"
    defaultName = "world"
    port = "../../bin/port.txt"
    input = "../../bin/input.txt"
     // use for go build file
    // input string = "input.txt"
    // port string = "port.txt"
)

func readPortFile(fileName string) string {
    file, err := ioutil.ReadFile(fileName)
    if err != nil {
		log.Fatal("port file not Found")
	}
	str := string(file) // convert content to a 'string'
	// fmt.Println(str)    // print the content as a 'string'
	return str
}

func readInputFile(fileName string) string {
    file, err := ioutil.ReadFile(fileName)
    if err != nil {
		log.Fatal("port file not Found")
	}
	str := string(file) // convert content to a 'string'
    // temp := strings.Split(str,"\n")
	// fmt.Println(str)    // print the content as a 'string'
	return str
}

func main() {
    //getPort
    port := readPortFile(port)
    input := readInputFile(input)
    // fmt.Println(input)
    // Set up a connection to the server.
    conn, err := grpc.Dial(address+port, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewGreeterClient(conn)

    // Contact the server and print out its response.
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := c.Calculate(ctx, &pb.CalcRequest{Name: input})
    if err != nil {
        log.Fatalf("error message: %v", err)
    }
    log.Printf("message: %s", r.GetMessage())
}


