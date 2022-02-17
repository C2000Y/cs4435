package main

import (
    "context"
    "log"
    "net"
    "io/ioutil"
    // "fmt"
    "strings"
    // "strconv"
    // "os"
    // "bufio"
    // "errors"

    "google.golang.org/grpc"

    pb "bank/proto"
)

const (
    path string = "../../bin/output.txt"
    port string = "../../bin/port.txt"
    // use for go build file
    // path string = "output.txt"
    // port string = "port.txt"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
    pb.UnimplementedBankServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) BankUpdate(ctx context.Context, input *pb.BankUpdateRequest) (*pb.BankUpdateReply, error) {
    log.Printf("Received message! processing...")
    action := strings.Split(input.GetName(), "\n")
    // err := writeFile(action)
    log.Println(action)
    return &pb.BankUpdateReply{Message: "done outputing file"}, nil
}

// func writeFile(action []string) error{
//     file, err := os.Create(path)
//     if err != nil {
// 		return errors.New("output path error")
// 	}
// 	defer file.Close()
//     w := bufio.NewWriter(file)
//     for _, v := range action {
//         action := strings.Fields(v)
//         result,err := getResult(action)
//         if err!=nil{
//             return err
//         }
//         fmt.Fprintln(w, result)
//     }
// 	return w.Flush()
// }

// calculation
// func getResult(action []string) (string,error){
//     a, err := strconv.Atoi(action[1])
//     if(err != nil){
//         return "error", errors.New("input not integer")
//     }
//     b, err := strconv.Atoi(action[2])
//     if(err != nil){
//         return "error", errors.New("input not integer")
//     }
//     // fmt.Println(a,action[0],b)
//     switch action[0] {
//     case "add": // input add
//         return strconv.Itoa(a + b),nil
// 	case "sub": // input sub
// 		return strconv.Itoa(a - b),nil
// 	case "mult": // input mult
//         return strconv.Itoa(a * b),nil
// 	case "div": // input div
// 		if b == 0 {
// 			return "error", errors.New("cannot divide by zero!")
// 		}
// 		return fmt.Sprintf("%v", float64(a) / float64(b)),nil
//     }
//     log.Fatal("error action!")
//     return "error", errors.New("unknown error")
// }

// read Port file to get port number
func readPortFile(fileName string) string {
    file, err := ioutil.ReadFile(fileName)
    if err != nil {
		log.Fatal("port file not Found")
	}
	str := string(file) // convert content to a 'string'
	// fmt.Println(str)    // print the content as a 'string'
	return str
}

func main() {
    //getPort
    port := readPortFile(port)
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    log.Println("listen to port ",port)
    s := grpc.NewServer()
    pb.RegisterBankServer(s, &server{})
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
