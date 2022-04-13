package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"net"
	"os"
	"strings"
	"strconv"
	"time"
	"math"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/google/uuid"
	chat "ChatRoom4435/proto"
)

//客户端管理
type ClientManager struct {
	//客户端 map 储存并管理所有的长连接client，在线的为true，不在的为false
	clients map[*Client]bool
	//web端发送来的的message我们用broadcast来接收，并最后分发给所有的client
	broadcast chan []byte
	//新创建的长连接client
	register chan *Client
	//新注销的长连接client
	unregister chan *Client
}

type node struct {
	// Self information
	Name       string
	Addr       string
	NodeName   string
	NodeNumber int
	NodeMap   map[string]int64
	// input operation
	OperationsFile string
	TimeStamp      int64
	operations     []string
	// Consul related variables
	SDAddress string
	SDKV      api.KV
	// used to make requests
	Clients map[string]chat.ChatClient
}

var (
	nd node
)

func (n *node) SendMessage(ctx context.Context, in *chat.Message) (*chat.MessageReply, error) {
	log.Printf("Received: %v, %v, %v", in.Name, in.Event, in.Content)
	n.NodeMap[in.Name] = in.Timestamp
	maxTime := math.Max(float64(n.TimeStamp), float64(in.Timestamp)) + 1
	n.TimeStamp = int64(maxTime)
	n.NodeMap[n.Name] = n.TimeStamp
	var (
		cnt string
		e string
	)
	switch(in.Event){
	case "hello":
		cnt = "received"
		e = "hello back"
		break
	case "msg":
		cnt = "received"
		e = "msg process"
		jsonMessage, _ := json.Marshal(&Message{Sender: in.Cid, Content: string(in.Content)})
		fmt.Println("get: ",string(in.Content))
		// TODO: process data
		manager.broadcast <- jsonMessage
		break
	}
	return &chat.MessageReply{Name: n.Name, Content: cnt, Event: e, Timestamp: n.TimeStamp}, nil
}

//客户端 Client
type Client struct {
	//用户id
	id string
	//连接的socket
	socket *websocket.Conn
	//发送的消息
	send chan []byte
}

//会把Message格式化成json
type Message struct {
	//消息struct
	Sender    string `json:"sender,omitempty"`    //发送者
	Recipient string `json:"recipient,omitempty"` //接收者
	Content   string `json:"content,omitempty"`   //内容
}

//创建客户端管理者
var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (manager *ClientManager) start() {
	for {
		select {
		//如果有新的连接接入,就通过channel把连接传递给conn
		case conn := <-manager.register:
			//把客户端的连接设置为true
			manager.clients[conn] = true
			//把返回连接成功的消息json格式化
			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			//调用客户端的send方法，发送消息
			manager.send(jsonMessage, conn)
			//如果连接断开了
		case conn := <-manager.unregister:
			//判断连接的状态，如果是true,就关闭send，删除连接client的值
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				manager.send(jsonMessage, conn)
			}
			//广播
		case message := <-manager.broadcast:
			//遍历已经连接的客户端，把消息发送给他们
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

//定义客户端管理的send方法
func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		//不给屏蔽的连接发送消息
		if conn != ignore {
			conn.send <- message
		}
	}
}

//定义客户端结构体的read方法
func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		//读取消息
		_, message, err := c.socket.ReadMessage()
		//如果有错误信息，就注销这个连接然后关闭
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}
		//如果没有错误信息就把信息放入broadcast
		nd.GreetAll(c.id, string(message))
		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
		fmt.Println("get: ",message)
		// time.Sleep(1 * time.Second)
		// // TODO: process data
		manager.broadcast <- jsonMessage
		fmt.Println("id:", c.id, ", ct",string(message))
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		//从send里读消息
		case message, ok := <-c.send:
			//如果没有消息
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//有消息就写入，发送给web端
			fmt.Println("send: ",message)
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// Start listening/service.
func (n *node) StartListening() {
	lis, err := net.Listen("tcp", n.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	_n := grpc.NewServer() // n is for serving purpose

	chat.RegisterChatServer(_n, n)
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

func wsHandler(res http.ResponseWriter, req *http.Request) {
	//将http协议升级成websocket协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	//每一次连接都会新开一个client，client.id通过uuid生成保证每次都是不同的
	
	client := &Client{id: uuid.Must(uuid.NewRandom()).String(), socket: conn, send: make(chan []byte)}
	//注册一个新的链接
	manager.register <- client

	//启动协程收web端传过来的消息
	go client.read()
	//启动协程把消息返回给web端
	go client.write()
}

// Setup a new grpc client for contacting the server at addr.
func (n *node) SetupClient(node string, addr string, timestamp int64, content string, event string, cid string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	n.Clients[node] = chat.NewChatClient(conn)
	r, err := n.Clients[node].SendMessage(context.Background(), &chat.Message{Name: n.Name, Timestamp: timestamp, Content: content, Event: event})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting from the other node: %s, content: %s, timeStamp: %d, event: %s", node, r.Content, r.Timestamp, r.Event)
	n.NodeMap[node] = r.Timestamp
	maxTime := math.Max(float64(n.TimeStamp), float64(r.Timestamp)) + 1
	n.TimeStamp = int64(maxTime)
	n.NodeMap[n.Name] = n.TimeStamp
	fmt.Println(n.NodeMap)
}


// Busy Work module, greet every new member you find
func (n *node) StartGreet() {
	kvpairs, _, err := n.SDKV.List(n.NodeName, nil)
	if err != nil {
		log.Panicln(err)
		return
	}
	fmt.Println(n.Name + " starts")
	timetmp := n.TimeStamp
	for _, kventry := range kvpairs {
		if strings.Compare(kventry.Key, n.Name) != 0 {
			fmt.Println("send to: ", kventry.Key)
			n.SetupClient(kventry.Key, string(kventry.Value), timetmp, "", "hello", "")
		}
	}
}

// Busy Work module, greet every new member you find
func (n *node) GreetAll(cid string, cnt string) {
	kvpairs, _, err := n.SDKV.List(n.NodeName, nil)
	if err != nil {
		log.Panicln(err)
		return
	}
	timetmp := n.TimeStamp
	for _, kventry := range kvpairs {
		if strings.Compare(kventry.Key, n.Name) != 0 {
			fmt.Println("send to: ", kventry.Key)
			n.SetupClient(kventry.Key, string(kventry.Value), timetmp, cnt, "msg", cid)
		}
	}
	fmt.Println(n.NodeMap)
}


func main() {
	args := os.Args[1:]
	// example: go run server/server.go "Node 1" :5000 localhost:8500
	// example: go run server/server.go "Node 2" :5001 localhost:8500
	if len(args) < 3 {
		fmt.Println("Arguments required: <name> <listening address> <consul address>")
		os.Exit(1)
	}

	// args in order
	name := args[0]
	listenAddr := args[1]
	consulAddr := args[2]

	nameSplit := strings.Split(name, " ")
	nodeNumber, err := strconv.Atoi(nameSplit[1])
	nodeName := nameSplit[0]
	if err != nil {
		log.Fatal("input argument <name> should have string with an integer, for example: Gnode 1, or node 2, or xxx 3")
	}

	nd = node{Name: name, Addr: listenAddr, SDAddress: consulAddr, NodeName: nodeName, NodeNumber: nodeNumber} // noden is for opeartional purposes

	nd.Clients = make(map[string]chat.ChatClient)
	nd.NodeMap = make(map[string]int64)

	// initialize timeStamp
	nd.TimeStamp = 0
	nd.NodeMap[name] = int64(0)

	fmt.Println("Initializing...")
	go nd.StartListening()
	nd.registerService()

	//开一个goroutine执行开始程序
	time.Sleep(5 * time.Second)
	nd.StartGreet()
	go manager.start()
	fmt.Println("Starting application...")
	//注册默认路由为 /ws ，并使用wsHandler这个方法
	http.HandleFunc("/ws", wsHandler)
	//监听本地的8011端口
	http.ListenAndServe(nd.Addr+"0", nil)
}