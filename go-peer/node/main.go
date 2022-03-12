package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	pb "go-peer/proto"

	"github.com/hashicorp/consul/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type node struct {
	// Self information
	Name       string
	Addr       string
	NodeNumber int
	NodeMap    map[string]int64
	// input operation
	OperationsFile string
	TimeStamp      int64
	operations     []string
	// Consul related variables
	SDAddress string
	SDKV      api.KV
	// used to make requests
	Clients map[string]pb.HelloServiceClient
}

type JsonData struct {
	Name      string
	AccountID int
	Balance   json.Number
}

var (
	AccountDetails []JsonData
	Operations     []string
	OperationIndex int = 0
)

const (
	accountFile    string = "bin/accounts.json"
	filePathPrefix string = "bin/"
	outPutPath     string = "bin/accountsUpdated"
)

// SayHello implements helloworld.GreeterServer
func (n *node) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	n.NodeMap[in.Name] = in.Timestamp
	maxTime := math.Max(float64(n.TimeStamp), float64(in.Timestamp)) + 1
	n.TimeStamp = int64(maxTime)
	if strings.Compare(in.Type, "newNode") == 0 {
		timestamp := strconv.Itoa(int(in.Timestamp))
		fmt.Println("Get from: " + in.Name + ", operation: " + in.Operation + ", timeStamp: " + timestamp)
		err := n.TimeOrder(in.Name, in.Operation)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := n.processAccount()
		if err != nil {
			return nil, err
		}
	}
	return &pb.HelloReply{Message: strconv.Itoa(int(n.TimeStamp))}, nil
}

func (n *node) TimeOrder(name string, operation string) error {
	err := errors.New("error")
	n.TimeStamp++
	n.NodeMap[n.Name] = n.TimeStamp
	// 排序插入，循环，假如大于前者且小于后者，插入到输出数组
	n.operations, err = insertSort(n.operations, operation)
	if err != nil {
		return err
	}
	return nil
}

func insertSort(originalArr []string, insertItem string) ([]string, error) {
	if len(originalArr) == 0 {
		finalArr := append(originalArr, insertItem)
		return finalArr, nil
	} else {
		for index, value := range originalArr {
			num := strings.Split(value, ":")[0]
			curr := strings.Split(num, ":")[0]
			currNum, err := strconv.ParseFloat(curr, 64)
			if err != nil {
				return nil, errors.New(curr + " --- parse float fialed")
			}

			input := strings.Split(insertItem, ":")[0]
			inputNum, err := strconv.ParseFloat(input, 64)
			if err != nil {
				return nil, errors.New(input + " --- parse float fialed")
			}
			if inputNum < currNum {
				finalArr := append(originalArr, "a")
				copy(finalArr[index+1:], finalArr[index:])
				finalArr[index] = insertItem
				return finalArr, nil
			}
		}
		finalArr := append(originalArr, insertItem)
		return finalArr, nil
	}
}

func (n *node) processAccount() (string, error) {
	operation := strings.Split(n.operations[0], ":")[1]
	n.operations = n.operations[1:]
	act := strings.Fields(operation)
	for index, acc := range AccountDetails {
		bal := string(AccountDetails[index].Balance)
		balance, err := strconv.ParseFloat(bal, 64)
		if err != nil {
			return operation, errors.New("account balance converting error")
		}
		amount, err := strconv.ParseFloat(act[2], 64)
		if err != nil {
			return operation, errors.New("inout amount converting to float error")
		}
		if strings.Compare(act[1], strconv.Itoa(acc.AccountID)) == 0 {
			switch act[0] {
			case "deposit":
				AccountDetails[index].Balance = json.Number(fmt.Sprintf("%0.2f", math.Round((balance+amount)*100)/100))
			case "withdraw":
				AccountDetails[index].Balance = json.Number(fmt.Sprintf("%0.2f", math.Round((balance-amount)*100)/100))
			case "interest":
				AccountDetails[index].Balance = json.Number(fmt.Sprintf("%0.2f", math.Round((balance+(balance*amount/100))*100)/100))
			}
		}
	}
	return operation, nil
}

// Start listening/service.
func (n *node) StartListening() {
	lis, err := net.Listen("tcp", n.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	_n := grpc.NewServer() // n is for serving purpose

	pb.RegisterHelloServiceServer(_n, n)
	// Register reflection service on gRPC server.
	reflection.Register(_n)

	// start listening
	if err := _n.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Register self with the service discovery module.
// This implementation simply uses the key-value store. One major drawback is that when nodes crash. nothing is updated on the key-value store. Services are a better fit and should be used eventually.
func (n *node) registerService() {
	config := api.DefaultConfig()
	config.Address = n.SDAddress
	consul, err := api.NewClient(config)
	if err != nil {
		log.Panicln("Unable to contact Service Discovery.")
	}

	kv := consul.KV()
	p := &api.KVPair{Key: n.Name, Value: []byte(n.Addr)}
	_, err = kv.Put(p, nil)
	if err != nil {
		log.Panicln("Unable to register with Service Discovery.")
	}

	// store the kv for future use
	n.SDKV = *kv

	log.Println("Successfully registered with Consul.")
}

// Start the node.
// This starts listening at the configured address. It also sets up clients for it's peers.
func (n *node) Start() {
	// init required variables
	n.Clients = make(map[string]pb.HelloServiceClient)
	n.NodeMap = make(map[string]int64)

	// start service / listening
	go n.StartListening()

	// register with the service discovery unit
	n.registerService()

	var err error
	// read accounts
	AccountDetails, err = initReadAccount(accountFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(AccountDetails)

	// read input Operations
	Operations, err = readInputOperations(n.OperationsFile)
	if err != nil {
		log.Fatal(err)
	}

	// start the main loop here
	// in our case, simply time out for 1 minute and greet all
	// wait for other nodes to come up
	// for {
	time.Sleep(5 * time.Second)
	n.GreetAll()
	// }
}

// Read file by fileName and decode josn, return file content
func initReadAccount(fileName string) ([]JsonData, error) {
	data := make([]JsonData, 0)
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return data, errors.New("input json file not found")
	}
	str := string(file) // convert content to a 'string'
	// fmt.Println(str)    // print the content as a 'string'
	jsonAsBytes := []byte(str)
	err2 := json.Unmarshal(jsonAsBytes, &data)
	// fmt.Printf("%#v\n", data)
	if err2 != nil {
		return data, errors.New("error json format from input file")
	}
	return data, nil
}

// read input operation file
func readInputOperations(fileName string) ([]string, error) {
	file, err := os.Open(filePathPrefix + fileName)
	if err != nil {
		return nil, errors.New("operation input file not Found, path:" + filePathPrefix + fileName)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	return text, nil
}

// Setup a new grpc client for contacting the server at addr.
func (n *node) SetupClient(node string, addr string, timestamp int64, operation string, operationType string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	n.Clients[node] = pb.NewHelloServiceClient(conn)
	r, err := n.Clients[node].SayHello(context.Background(), &pb.HelloRequest{Name: n.Name, Timestamp: timestamp, Operation: operation, Type: operationType})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting from the other node: %s, timeStamp: %s", node, r.Message)
	// TODO:  记录收到的node的timestamp
	stampstamp, err := strconv.Atoi(r.Message)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	if strings.Compare(operationType, "newNode") == 0 {
		n.NodeMap[node] = int64(stampstamp)
		maxTime := math.Max(float64(n.TimeStamp), float64(stampstamp)) + 1
		n.TimeStamp = int64(maxTime)
		n.NodeMap[n.Name] = int64(stampstamp)
	} else {
		n.NodeMap[n.Name] = int64(stampstamp)
	}
}

// Busy Work module, greet every new member you find
func (n *node) GreetAll() {
	kvpairs, _, err := n.SDKV.List("Gnode", nil)
	if err != nil {
		log.Panicln(err)
		return
	}
	fmt.Println(n.Name + "starts")
	n.NodeMap[n.Name] = 0
	for strings.Compare(Operations[OperationIndex], "done") != 0 {
		n.TimeStamp++
		operation := strconv.Itoa(int(n.TimeStamp)) + "." + strconv.Itoa(n.NodeNumber) + ":" + Operations[OperationIndex]
		timetmp := n.TimeStamp
		n.operations, err = insertSort(n.operations, operation)
		if err != nil {
			log.Fatalf("%v", err)
		}
		for _, kventry := range kvpairs {
			if strings.Compare(kventry.Key, n.Name) != 0 {
				fmt.Println("send to: ", kventry.Key, ", time:", int(timetmp), ", operation: ", Operations[OperationIndex])
				n.SetupClient(kventry.Key, string(kventry.Value), timetmp, operation, "newNode")
			}
		}
		//判断当前node的操作是否在队列第一位，并确保所node的timestamp大于此操作的timestamp
		if len(n.operations) > 0 {
			n.DoOperation(timetmp)
		}
		time.Sleep(time.Second / 2)
		OperationIndex++
	}
	for len(n.operations) > 0 {
		n.DoOperation(n.TimeStamp)
		time.Sleep(time.Second / 2)
	}
	// for _, num := range n.operations {
	// 	fmt.Print(num + " | ")
	// }
	fmt.Println("waiting...")
	time.Sleep(5 * time.Second)
	fmt.Println(AccountDetails)
	err = n.Output(outPutPath)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func (n *node) DoOperation(timetmp int64) {
	kvpairs, _, err := n.SDKV.List("Gnode", nil)
	if err != nil {
		log.Panicln(err)
		return
	}
	fmt.Println(n.operations[0])
	nodeNum := strings.Split(strings.Split(n.operations[0], ":")[0], ".")
	if strings.Compare(nodeNum[1], strconv.Itoa(n.NodeNumber)) == 0 {
		timestp, err := strconv.Atoi(nodeNum[0])
		if err != nil {
			log.Fatalf("%v", err)
		}
		var do bool = true
		for _, kventry := range kvpairs {
			if n.NodeMap[kventry.Key] <= int64(timestp) {
				do = false
			}
		}
		if do {
			operation, err := n.processAccount()
			if err != nil {
				log.Fatalf("%v", err)
			}
			for _, kventry := range kvpairs {
				if strings.Compare(kventry.Key, n.Name) != 0 {
					fmt.Println("release: " + operation + ", release to: " + kventry.Key + ", time:" + strconv.Itoa(int(timetmp)))
					n.SetupClient(kventry.Key, string(kventry.Value), n.TimeStamp, "", "release")
				}
			}
		}
	}
}

func (n *node) Output(fileName string) error {
	nodeName := strings.Split(n.Name, " ")
	fileName += nodeName[0] + nodeName[1] + ".json"
	json, _ := json.MarshalIndent(AccountDetails, "", "  ")
	err := ioutil.WriteFile(fileName, json, 0644)
	if err != nil {
		return errors.New("output file error")
	}
	fmt.Println("Write Success! " + fileName + " Created")
	return nil
}

func main() {
	// pass the port as an argument and also the port of the other node
	args := os.Args[1:]
	// example: go run node/main.go "Node 1" :5000 localhost:8500 input
	if len(args) < 4 {
		fmt.Println("Arguments required: <name> <listening address> <consul address> <operation file>")
		os.Exit(1)
	}

	// args in order
	name := args[0]
	listenaddr := args[1]
	sdaddress := args[2]
	operation := args[3]

	nodeNumber, err := strconv.Atoi(strings.Split(name, " ")[1])
	if err != nil {
		log.Fatal("input argument <name> should have string with an integer, for example: Gnode 1")
	}
	noden := node{Name: name, Addr: listenaddr, SDAddress: sdaddress, OperationsFile: operation, Clients: nil, TimeStamp: 0, NodeNumber: nodeNumber, NodeMap: nil} // noden is for opeartional purposes

	// start the node
	noden.Start()
}
