package main

import (
    "context"
    "log"
    "net"
    "io/ioutil"
    "fmt"
    "strings"
    "strconv"
    "errors"
    "encoding/json"
    "math"

    "google.golang.org/grpc"

    pb "bank/proto"
)

const (
    path string = "../../bin/accountsUpdated.json"
    port string = "../../bin/port.txt"
    inputJson string = "../../bin/accounts.json"
    // use for go build file
    // path string = "accountsUpdated.txt"
    // port string = "port.txt"
    // inputJson string = "accounts.json"
)

type JsonData struct {
	Name      string
	AccountID int
	Balance   json.Number
}

// server is used to implement helloworld.GreeterServer.
type server struct {
    pb.UnimplementedBankServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) BankUpdate(ctx context.Context, input *pb.BankUpdateRequest) (*pb.BankUpdateReply, error) {
    log.Printf("Received message! processing...")
    action := strings.Split(input.GetName(), "\n")
    err := writeFile(action)
    return &pb.BankUpdateReply{Message: "done outputing file"}, err
}

func writeFile(action []string) error{
    jsContent, err := readFile(inputJson)
    if err != nil {
        return err
    }
    data, err := decodeJson(jsContent)
    if err != nil {
        return err
    }
    result, err := changeBalance(data, action)
    output(result, path)
    return nil
}

// changeBalance
func changeBalance(data []JsonData, action []string) ([]JsonData,error) {
    s := map[int]bool{}
	for index, acc := range data {
        for j, actions := range action {
            _, ok := s[j]
            // skip used action
            if(ok != true){
                act := strings.Fields(actions)
                bal := string(data[index].Balance)
                balance, err := strconv.ParseFloat(bal, 64)
                if err != nil {
                    return data, errors.New("account balance converting error")
                }
                amount, err := strconv.ParseFloat(act[2], 64)
                if err != nil {
                    return data, errors.New("inout amount converting to float error")
                }
                if strings.Compare(act[1],strconv.Itoa(acc.AccountID)) == 0 {
                    switch act[0] {
                    case "deposit":
                        data[index].Balance = json.Number(fmt.Sprintf("%0.2f", math.Round((balance + amount)*100)/100))
                    case "withdraw":
                        data[index].Balance = json.Number(fmt.Sprintf("%0.2f", math.Round((balance - amount)*100)/100))
                    case "interest":
                        data[index].Balance = json.Number(fmt.Sprintf("%0.2f", math.Round((balance + (balance * amount / 100))*100)/100))
                    }
                    // flag used actions
                    s[j] = true
                }
            }
        }
	}
	return data,nil
}

// Read file by fileName, return file content
func readFile(fileName string) (string, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", errors.New("input json file not found")
	}
	str := string(file) // convert content to a 'string'
	// fmt.Println(str)    // print the content as a 'string'
	return str, nil
}

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

// input file content and convert it to listed JsonData struct
func decodeJson(content string) ([]JsonData, error){
	jsonAsBytes := []byte(content)
	data := make([]JsonData, 0)
	err := json.Unmarshal(jsonAsBytes, &data)
	// fmt.Printf("%#v\n", data)
	if err != nil {
		return data, errors.New("error json format from input file")
	}
	return data, nil
}

func output(data []JsonData, fileName string) error{
	// json, _ := json.Marshal(data)
	json, _ := json.MarshalIndent(data, "", "  ")
	err := ioutil.WriteFile(fileName, json, 0644)
	if err != nil {
		return errors.New("output file error")
	}
	fmt.Println("Write Success! accountsUpdated.json Created")
    return nil
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
