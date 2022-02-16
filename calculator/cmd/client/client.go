package main

import (
    "context"
    "log"
    // "os"
    "time"
    "io/ioutil"
    "fmt"
    // "strings"

    "google.golang.org/grpc"

    pb "cs4435/calculator/proto"
)

const (
    address     = "localhost"
    defaultName = "world"
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
    port := readPortFile("../../bin/port.txt")
    input := readInputFile("../../bin/input.txt")
    fmt.Println(input)
    // Set up a connection to the server.
    conn, err := grpc.Dial(address+port, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewGreeterClient(conn)

    // Contact the server and print out its response.
    // name := defaultName
    // if len(os.Args) > 1 {
    //     name = os.Args[1]
    // }
    // fmt.Println(os.Args[1])
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := c.Calculate(ctx, &pb.CalcRequest{Name: input})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("Greeting: %s", r.GetMessage())
}


