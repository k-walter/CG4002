package eval

import (
	"bufio"
	"cg4002/eComm/common"
	pb "cg4002/protos"
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
)

type EEvalResp struct {
	*pb.State
	Time time.Time
}

type Client struct {
	// OPTIMIZATION reserve mem for reader, data
	conn     net.Conn
	key      string
	chEngine chan *pb.State
}

func Make(args *common.Arg) *Client {
	e := Client{
		key:      args.EvalKey,
		chEngine: make(chan *pb.State, common.ChSz),
	}

	// Connect to eval server
	var err error
	e.conn, err = net.Dial("tcp", args.EvalAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to state updates
	common.SubOld(common.State2Eval, func(i interface{}) {
		go func(i *pb.State) { e.chEngine <- i }(i.(*pb.State))
	})

	return &e
}

func (c *Client) Close() {
	_ = c.conn.Close()
}

func (c *Client) Run() {
	for curState := range c.chEngine {
		c.send(curState)
		trueState := c.receive()
		// OPTIMIZE add RTT/2 to shield time left?
		common.PubOld(common.State2Eng, trueState)
	}
}

func (c *Client) send(s *pb.State) {
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

func (c *Client) receive() *pb.State {
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
