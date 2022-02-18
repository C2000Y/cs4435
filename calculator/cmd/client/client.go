package main

import (
    "context"
    "log"
    "time"
    "io/ioutil"
    "fmt"
    "os"
    "bufio"
    "strings"

    "google.golang.org/grpc"

    pb "calculator/proto"
)

const (
    address     = "localhost"
    defaultName = "world"
    port = "../../bin/port.txt"
    input = "../../bin/input.txt"
    path string = "../../bin/output.txt"
     // use for go build file
    // input string = "input.txt"
    // port string = "port.txt"
    // path string = "output.txt"
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

func writeFile(action string) error{
    file, err := os.Create(path)
    if err != nil {
		log.Fatal("file path error")
	}
	defer file.Close()
    w := bufio.NewWriter(file)
    operations := strings.Split(action, "/n")
    for _, v := range operations {
        fmt.Fprintln(w, v)
    }
	return w.Flush()
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
    c := pb.NewCalculatorClient(conn)

    // Contact the server and print out its response.
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := c.Calculate(ctx, &pb.CalcRequest{Name: input})
    if err != nil {
        log.Fatalf("error message: %v", err)
    }
    wirteErr := writeFile(r.GetMessage())
    if wirteErr != nil {
        log.Fatalf("error message: %v", wirteErr)
    }
    log.Printf("File successfully output")
}


