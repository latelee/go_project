package ws

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
    "errors"

	"github.com/gorilla/websocket"
    
    "faces/src/baoqi/library/common"
)

type BaseMessage struct {
    Id string `json:"id"`
    Op string `json:"op"`
    Timestamp int64 `json:"timestamp"`
    Data interface{} `json:"data"`
}

type WSClient struct {
    Mutex        sync.Mutex  // gorilla 不能读或写操作不同同时进行，需要互斥
    Conn         *websocket.Conn
    ReceiverChan chan interface{}
    Url          string
    Id           string
    Timeout      int
}

func InitWSServer(id string, conn *websocket.Conn) (cli* WSClient){
    cli = &WSClient {
        Id:   id,
        Conn: conn,
        //ReceiverChan: make(chan interface{}, 128),
    }

    return cli
}

func InitWSClient(url string, id string, timeout int) (cli* WSClient){
    cli = &WSClient {
        Conn: nil,
        ReceiverChan: make(chan interface{}, 128),
        Url: url,
        Id:  id,
        Timeout: timeout,
    }
    
    return cli
}

func (cli *WSClient) Connect() {
    url := fmt.Sprintf("ws://%s/device/ws/%s", cli.Url, cli.Id)
    common.Logger.Info("connecting to ", url)
    if cli.Conn != nil {
        cli.Conn.Close()
        cli.Conn = nil
    }

    for {
        c, _, err := websocket.DefaultDialer.Dial(url, nil)
        if err == nil {
            common.Logger.Info("websocket connect ok")
            cli.Conn = c
            break
        }
        common.Logger.Info("websocket connect failed: ", err)
        cli.Conn = nil // ??? need
        time.Sleep(time.Duration(cli.Timeout)*time.Second)
        //time.Sleep(10*time.Second)
    }
}

func (cli *WSClient) HandleSend() {
    for {
        select {
        case msg, ok := <- cli.ReceiverChan:
            if !ok {
                common.Logger.Info("msg <- error")
                return
            }
            if cli.Conn == nil {
                common.Logger.Info("ws is nil, will not send")
                continue
            }
            common.Logger.Info("..........send msg...... ", string(msg.([]byte)))
            // 发送消息
            cli.Mutex.Lock()
            err := cli.Conn.WriteMessage(websocket.TextMessage, msg.([]byte))
            cli.Mutex.Unlock()
            if err != nil {
                common.Logger.Info("WriteMessage:", err)
                cli.Connect()
            }
        }       
    }
}

func (cli *WSClient) HandleRecieve() {
    for {
        if cli.Conn == nil {
            common.Logger.Info("ws is nil, will not read")
            continue
        }
        cli.Mutex.Lock()
        _, message, err := cli.Conn.ReadMessage()
        cli.Mutex.Unlock()
        if err != nil {
            common.Logger.Info("read error:", err)
            cli.Conn.Close()
            cli.Connect()
            continue
        }
        
        var recvMsg map[string]interface{}
        err = json.Unmarshal(message, &recvMsg)
        if err != nil {
            common.Logger.Info("Unmarshal msg failed.")
            continue
        }
    
        if recvMsg["op"] == "heartbeat" {
            common.Logger.Info("recv heartbeat.....")
            //setHeart()
        }
        common.Logger.Info("client recv: ", string(message))
    }
}

func (cli *WSClient) Send111(content interface{}) {
    cli.ReceiverChan <- content
}

// 单纯的封装

func (cli *WSClient) Recv() (v BaseMessage, err error) {
    if cli != nil {
        // 读取ws中的数据
        cli.Mutex.Lock()
        _, message, err := cli.Conn.ReadMessage()
        cli.Mutex.Unlock()
        if err != nil {
            common.Logger.InfoFormat("ReadMessage failed from %s err: %s", cli.Conn.LocalAddr(), err)
            return BaseMessage{}, errors.New("websocket ReadMessage failed")
        }
        // 解析包
        err = json.Unmarshal(message, &v)
        if err != nil {
            common.Logger.InfoFormat("Unmarshal msg failed %s , msg: %s", err, string(message))
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
            common.Logger.InfoFormat("will send to %s op: %s", v.Id, v.Op)
            allJson, err := json.Marshal(v)
            if err != nil {
                common.Logger.Error("Marshal err: ", err)
                return errors.New("BaseMessage Marshal failed")
            }
            cli.Mutex.Lock()
            err = cli.Conn.WriteMessage(websocket.TextMessage, allJson)
            cli.Mutex.Unlock()
            if err != nil {
                common.Logger.Error("WriteMessage err: ", err)
                return errors.New("websocket WriteMessage failed")
            }
        } else {
            common.Logger.InfoFormat("ws client %s(op: %s) not found", v.Id, v.Op)
        }
    }
    return nil
}

/*
func (cli *WSClient) SimpleSend(content interface{}) (error) {
    if cli.Conn == nil {
        common.Logger.Info("ws is nil, will not send")
        return errors.New("ws is nil")
    }
    cli.Mutex.Lock()
    err := cli.Conn.WriteMessage(websocket.TextMessage, content.([]byte))
    cli.Mutex.Unlock()

    return err
}

func (cli *WSClient) SimpleRecv() ([]byte, error){
    if cli.Conn == nil {
        common.Logger.Info("ws is nil, will not read")
        return nil, errors.New("ws is nil")
    }
    cli.Mutex.Lock()
    _, message, err := cli.Conn.ReadMessage()
    cli.Mutex.Unlock()
    
    return message, err
}
*/
