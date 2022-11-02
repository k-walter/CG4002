package eval

import (
	"bufio"
	"cg4002/eComm/common"
	pb "cg4002/protos"
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

type Client struct {
	// OPTIMIZATION reserve mem for reader, data
	conn     net.Conn
	key      string
	chEngine chan *pb.State
	mu       sync.Mutex
}

func Make(args *common.Arg) *Client {
	e := Client{
		conn:     nil,
		key:      args.EvalKey,
		chEngine: make(chan *pb.State, common.ChSz),
		mu:       sync.Mutex{},
	}

	// Connect to eval server
	var err error
	e.conn, err = net.Dial("tcp", args.EvalAddr)
	if err != nil {
		log.Fatal(err)
	}

	return &e
}

func (c *Client) Run() {
}

func (c *Client) Close() {
	_ = c.conn.Close()
}

func (c *Client) BlockingSend(s *pb.State) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get json
	msg := common.PbToJson(s.ProtoReflect())
	log.Println("eval|Send", string(msg))

	// Encrypt json
	iv := common.MakeIv()
	encMsg := common.Aes256(msg, c.key, iv, aes.BlockSize)
	data := base64.StdEncoding.EncodeToString(append(iv, encMsg...))

	// Send length
	_, err := c.conn.Write([]byte(fmt.Sprintf("%v_", len(data))))
	if err != nil {
		log.Fatal(err)
	}

	// Send data
	_, err = c.conn.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) BlockingRecv() *pb.State {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get length
	r := bufio.NewReader(c.conn)
	lenStr, err := r.ReadString('_')
	if err != nil {
		log.Fatal(err)
	}

	// Parse length
	dataSz, err := strconv.Atoi(lenStr[:len(lenStr)-1])
	if err != nil || dataSz <= 0 {
		log.Fatal(dataSz, err)
	}

	// Get data
	data := make([]byte, dataSz)
	if _, err := io.ReadFull(r, data); err != nil {
		log.Fatal(err)
	}
	log.Println("eval|Received", string(data))

	// Json to pb
	s := pb.State{}
	err = protojson.Unmarshal(data, &s)
	if err != nil {
		log.Fatal(err)
	}

	return &s
}
