package main

import (
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
	Name string
	Addr string

	// input operation
	OperationsFile string
	TimeStamp      int64
	operations     string

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
	AccountDetails []string
	Operations     []string
	OperationIndex int = 0
)

const (
	accountFile    string = "bin/accounts.json"
	filePathPrefix string = "bin/"
)

// SayHello implements helloworld.GreeterServer
func (n *node) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	timestamp := strconv.Itoa(int(in.Timestamp))
	fmt.Println("From " + n.Name + " Input name: " + in.Name + ", operation: " + in.Operation + "timeStamp: " + timestamp)
	err := n.TimeOrder(in.Name, in.Operation, in.Timestamp)
	if err != nil {
		return nil, err
	}
	return &pb.HelloReply{Message: "Hello from " + n.Name + ", timestamp:" + strconv.Itoa(int(n.TimeStamp)) + ", operation:" + n.operations}, nil
}

func (n *node) TimeOrder(name string, operation string, timestamp int64) error {
	maxTime := math.Max(float64(n.TimeStamp), float64(timestamp))
	n.TimeStamp = int64(maxTime)
	return nil
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

	// start service / listening
	go n.StartListening()

	// register with the service discovery unit
	n.registerService()

	// read accounts
	json, err := initReadAccount(accountFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(json)

	// read input Operations
	Operations, err = readInputOperations(n.OperationsFile)
	if err != nil {
		log.Fatal(err)
	}

	// start the main loop here
	// in our case, simply time out for 1 minute and greet all
	// wait for other nodes to come up
	for {
		time.Sleep(2 * time.Second)
		n.GreetAll()
	}
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
	file, err := ioutil.ReadFile(filePathPrefix + fileName)
	if err != nil {
		return nil, errors.New("operation input file not Found, path:" + filePathPrefix + fileName)
	}
	str := string(file)
	result := strings.Split(str, "\n")
	return result, nil
}

// Setup a new grpc client for contacting the server at addr.
func (n *node) SetupClient(name string, addr string) {
	// setup connection with other node
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	n.TimeStamp = n.TimeStamp + 1
	// n.operations = append(n.operations, strconv.Itoa(int(n.TimeStamp))+":"+strconv.Itoa(OperationIndex)+";")

	defer conn.Close()
	n.Clients[name] = pb.NewHelloServiceClient(conn)
	// fmt.Println("Name:" + n.Name + ",Timestamp:" + strconv.Itoa(n.TimeStamp) + ", Operation:" + n.operations)
	// fmt.Println(&pb.HelloRequest{Name: n.Name, Timestamp: strconv.Itoa(n.TimeStamp), Operation: n.operations})
	r, err := n.Clients[name].SayHello(context.Background(), &pb.HelloRequest{Name: n.Name, Timestamp: n.TimeStamp, Operation: n.operations})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting from the other node: %s", r.Message)

}

// Busy Work module, greet every new member you find
func (n *node) GreetAll() {
	// get all nodes -- inefficient, but this is just an example
	kvpairs, _, err := n.SDKV.List("Gnode", nil)
	if err != nil {
		log.Panicln(err)
		return
	}

	// fmt.Print("Found nodes: ")
	for _, kventry := range kvpairs {
		fmt.Println(kventry.Key)
		if strings.Compare(kventry.Key, n.Name) == 0 {
			// ourself
			// fmt.Println(kventry.Key)
			continue
		}
		if n.Clients[kventry.Key] == nil {
			fmt.Println("New member: ", kventry.Key)
			// connection not established previously
			n.SetupClient(kventry.Key, string(kventry.Value))
		}
	}
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

	noden := node{Name: name, Addr: listenaddr, SDAddress: sdaddress, OperationsFile: operation, Clients: nil, TimeStamp: 1} // noden is for opeartional purposes

	// start the node
	noden.Start()
}
