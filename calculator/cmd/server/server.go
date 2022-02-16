package main

import (
    "context"
    "log"
    "net"
    "io/ioutil"
    "fmt"
    "strings"
    "strconv"
    // "reflect"

    "google.golang.org/grpc"

    pb "cs4435/calculator/proto"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
    pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) Calculate(ctx context.Context, input *pb.CalcRequest) (*pb.CalcReply, error) {
    // log.Printf("Received: %v", input.GetName())
    action := strings.Split(input.GetName(), "\n")
    writeFile(action)
    return &pb.CalcReply{Message: "Hello "}, nil
}

func writeFile(action []string){
    for _, v := range action {
        action := strings.Fields(v)
        result := getResult(action)
        fmt.Println(result)
    }
}

func getResult(action []string) string{
    a, err := strconv.Atoi(action[1])
    numberErr(err)
    b, err := strconv.Atoi(action[2])
    numberErr(err)
    // fmt.Println(a,action[0],b)
    switch action[0] {
    case "add": // input add
        return strconv.Itoa(a + b)
	case "sub": // input sub
		return strconv.Itoa(a - b)
	case "mult": // input mult
        return strconv.Itoa(a * b)
	case "div": // input div
		if b == 0 {
			log.Fatal("cannot divide by 0")
			break
		}
		return fmt.Sprintf("%v", float64(a) / float64(b))
    }
    log.Fatal("error action!")
    return "error"
}

func readPortFile(fileName string) string {
    file, err := ioutil.ReadFile(fileName)
    if err != nil {
		log.Fatal("port file not Found")
	}
	str := string(file) // convert content to a 'string'
	fmt.Println(str)    // print the content as a 'string'
	return str
}

func numberErr(err error){
    if(err != nil){
        log.Fatal("number input not integer")
    }
}

func main() {
    //getPort
    port := readPortFile("../../bin/port.txt")
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterGreeterServer(s, &server{})
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
