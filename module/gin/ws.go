/*
本文返回json格式

发送思路：
接到请求，保存到全局的map中，做通道，单独函数发送，但要指定是哪个ID，查看是否有连接，有则发送。

TODO：对方断开事件通知？
TODO：每个客户端单独一个对象？这样可单独保存其ID
*/

package gin

import (
	// "fmt"
	//"strconv"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"webdemo/pkg/klog"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var setHeartBeat = 1

var upgrader = websocket.Upgrader{
	// 这里有很多个参数
	//ReadBufferSize:   1024,
	//WriteBufferSize:  1024,
	//HandshakeTimeout: 5 * time.Second,
	//EnableCompression: true,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSClient struct {
	Mutex        sync.Mutex // gorilla 不能读或写操作不同同时进行，需要互斥
	Conn         *websocket.Conn
	ReceiverChan chan interface{}
	Url          string
	Id           string
}

var wonce sync.Once

var gWSMsg = make(chan *BaseMessage, 128)

func InitWSServer(id string, conn *websocket.Conn) (cli *WSClient) {
	cli = &WSClient{
		Id:   id,
		Conn: conn,
		//ReceiverChan: make(chan interface{}, 128),
	}

	//go cli.HandleSend()
	//go cli.HandleRecieve()

	//defer cli.Conn.Close()

	return cli
}

func (cli *WSClient) Recv() (v BaseMessage, err error) {
	if cli != nil {
		klog.Infoln("waiting for !!!!!!!")
		klog.Printf("ws: %d %s\n", &cli.Conn, cli.Id)
		// 读取ws中的数据
		_, message, err := cli.Conn.ReadMessage()
		if err != nil {
			klog.Infof("ReadMessage failed from %s err: %s", cli.Conn.LocalAddr(), err)
			return BaseMessage{}, errors.New("websocket ReadMessage failed")
		}
		// 解析包
		err = json.Unmarshal(message, &v)
		if err != nil {
			klog.Infof("Unmarshal msg failed %s , msg: %s", err, string(message))
			return BaseMessage{}, errors.New("BaseMessage Unmarshal failed")
		} else {
			return v, nil
		}
	}
	return BaseMessage{}, errors.New("WSClient is nil")
}

func (cli *WSClient) Send(v *BaseMessage) (err error) {
	if cli != nil {
		if v.Id == cli.Id {
			klog.Infof("will send to %s op: %s", v.Id, v.Op)
			allJson, err := json.Marshal(v)
			if err != nil {
				klog.Println("Marshal err: ", err)
				return errors.New("BaseMessage Marshal failed")
			}
			err = cli.Conn.WriteMessage(websocket.TextMessage, allJson)
			if err != nil {
				klog.Println("WriteMessage err: ", err)
				return errors.New("websocket WriteMessage failed")
			}
		} else {
			klog.Infof("ws client %s(op: %s) not found", v.Id, v.Op)
		}
	}
	return nil
}

func genUserInfo() interface{} {

	return ""
}

// 注：如果没有连接，则不向通道发数据
func WSSend(content *BaseMessage) {
	gWSMsg <- content
}

func WSHandler(c *gin.Context) {
	id := c.Param("id")

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		klog.Println("Upgrade error: ", err)
		return
	}
	klog.Printf("local addr: %v remote addr: %v, id: %v\n", ws.LocalAddr(), ws.RemoteAddr(), id)

	defer ws.Close()
	klog.Printf("ws: %d\n", &ws)
	// 保存
	//InitWSServer(id, ws)
	client := InitWSServer(id, ws) //&Client{Id: id, Conn: ws}
	//gclientList[id] = client

	ws.SetPongHandler(func(string) error { klog.Println("pong handle"); return nil })

	// 客户端主动发关闭事件，在此函数中响应，如果没发，则不会响应
	// for _, v := range gclientList { if v.Conn == ws { delete(gclientList, v.Id) } };
	ws.SetCloseHandler(func(code int, text string) error {
		klog.Println("close handle from ", ws.RemoteAddr())
		message := websocket.FormatCloseMessage(code, "")
		ws.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	// 单独协程，处理发送的命令
	go func() {
		for {
			select {
			case v := <-gWSMsg:
				err := client.Send(v)
				if err != nil {
					return
				}
				//case <-time.After(5 * time.Second):
				//    klog.Infoln("timeout!!!!!!")
			}
		}
	}()

	// 如果在本函数中循环处理，则在客户端断开连接后，gin才返回，如果是协程则马上返回
	for {
		//klog.Info("waiting for read....")
		recvMsg, err := client.Recv()
		if err != nil {
			return
		}
		klog.Printf("op: %v id: %v timestamp: %v\n", recvMsg.Op, recvMsg.Id, recvMsg.Timestamp)
		var allData interface{}
		// 3. 分析包、处理
		// 判断条件，发送。。。
		if recvMsg.Op == "heartbeat" {
			if setHeartBeat == 1 { // 临时测试用
				allData = ""
			} else {
				klog.Println("not send heartbeat...")
				continue
			}
		} else if recvMsg.Op == "getuser" {
			allData = genUserInfo()
		} else {
			allData = ""
		}

		sendMsg := &BaseMessage{
			Op:        recvMsg.Op,
			Id:        recvMsg.Id,
			Timestamp: time.Now().UnixNano() / 1e6,
			Data:      allData,
		}
		// 4. 回写ws数据
		err = client.Send(sendMsg)
		if err != nil {
			return
		}
	}

}

func WSHandler2(c *gin.Context) {
	id := c.Param("id")

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		klog.Println("Upgrade error: ", err)
		return
	}
	klog.Printf("local addr: %v remote addr: %v, id: %v\n", ws.LocalAddr(), ws.RemoteAddr(), id)

	defer ws.Close()
	// 保存
	//InitWSServer(id, ws)
	//client := InitWSServer(id, ws) //&Client{Id: id, Conn: ws}
	//gclientList[id] = client

	ws.SetPongHandler(func(string) error { klog.Println("pong handle"); return nil })

	// 客户端主动发关闭事件，在此函数中响应，如果没发，则不会响应
	// for _, v := range gclientList { if v.Conn == ws { delete(gclientList, v.Id) } };
	ws.SetCloseHandler(func(code int, text string) error {
		klog.Println("close handle from ", ws.RemoteAddr())
		message := websocket.FormatCloseMessage(code, "")
		ws.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	// 单独协程，处理发送的命令
	go func() {
		for {
			select {
			case v := <-gWSMsg:
				if v.Id != id {
					klog.Infof("ws client %s(op: %s) not found", v.Id, v.Op)
					continue
				}
				klog.Infof("will send to %s op: %s", v.Id, v.Op)
				allJson, err := json.Marshal(v)
				err = ws.WriteMessage(websocket.TextMessage, allJson)
				if err != nil {
					klog.Println("WriteMessage err: ", err)
					//return
				}
				//case <-time.After(5 * time.Second):
				//    klog.Infoln("timeout!!!!!!")
			}
		}
	}()

	// 如果在本函数中循环处理，则在客户端断开连接后，gin才返回，如果是协程则马上返回
	for {
		// 1. 读取ws中的数据
		_, message, err := ws.ReadMessage()
		if err != nil {
			klog.Infof("ReadMessage failed from %s err: %s", ws.LocalAddr(), err)
			return
		}

		klog.Printf("server recv msg: %s", message)
		// 2. 解析包
		var recvMsg BaseMessage
		err = json.Unmarshal(message, &recvMsg)
		if err != nil {
			klog.Println("Unmarshal msg failed.")
			continue
		}

		klog.Printf("op: %v id: %v timestamp: %v\n", recvMsg.Op, recvMsg.Id, recvMsg.Timestamp)
		var allData interface{}
		// 3. 分析包、处理
		// 判断条件，发送。。。
		if recvMsg.Op == "heartbeat" {
			if setHeartBeat == 1 { // 临时测试用
				allData = ""
			} else {
				klog.Println("not send heartbeat...")
				continue
			}
		} else if recvMsg.Op == "getuser" {

			allData = genUserInfo()
		}

		allMsg := make(map[string]interface{})
		allMsg["op"] = "get_time"
		allMsg["id"] = recvMsg.Id
		allMsg["timestamp"] = time.Now().UnixNano() / 1e6

		allMsg["data"] = allData
		allJson, err := json.Marshal(allMsg)
		if err != nil {
			klog.Println("Marshal err: ", err)
			continue
		}
		// 4. 回写ws数据
		err = ws.WriteMessage(websocket.TextMessage, allJson)
		if err != nil {
			klog.Println("WriteMessage err: ", err)
			return
		}
	}

}
