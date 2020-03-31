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
    "net/http"
    "encoding/json"
    "time"

    "k8s.io/klog"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "sync"
    
)

var upgrader = websocket.Upgrader{
    // 这里有很多个参数
    //ReadBufferSize:   1024,
    //WriteBufferSize:  1024,
    //HandshakeTimeout: 5 * time.Second,
    //EnableCompression: true,

    CheckOrigin: func (r *http.Request) bool {
        return true
    },
}


type WSClient struct {
    Mutex        sync.Mutex  // gorilla 不能读或写操作不同同时进行，需要互斥
    Conn         *websocket.Conn
    ReceiverChan chan interface{}
    Url          string
    Id           string
}

var wonce sync.Once

func InitWSServer(id string, conn *websocket.Conn) (cli* WSClient){
    cli = &WSClient {
        Id: id,
        Conn: conn,
        ReceiverChan: make(chan interface{}, 128),
    }
    
    go cli.HandleSend()
    go cli.HandleRecieve()
    
    // 仅查看在线客户端，调试用
    wonce.Do(func() {
        go showList()
    })

    return cli
}

// 不需要
/*
func InitWSClient(url, id string) (cli* WSClient){
    cli = &WSClient {
        Conn: nil,
        ReceiverChan: make(chan interface{}, 128),
        Url: url, //Config.ServerUrl,
        Id:  id, //Config.Key,
    }
    
    return cli
}

func (cli *WSClient) Connect() {
    url := fmt.Sprintf("ws://%s/device/ws/%s", cli.Url, cli.Id)
    klog.Info("connecting to ", url)
    if cli.Conn != nil {
        cli.Conn.Close()
        cli.Conn = nil
    }

    for {
        c, _, err := websocket.DefaultDialer.Dial(url, nil)
        if err == nil {
            klog.Info("websocket connect ok")
            cli.Conn = c
            break
        }
        klog.Info("websocket connect failed: ", err)
        cli.Conn = nil // ??? need
        time.Sleep(time.Duration(Config.HeartBeatTime)*time.Second)
    }
}
*/

func (cli *WSClient) HandleSend() {
    for {
        select {
        case msg, ok := <- cli.ReceiverChan:
            if !ok {
                klog.Info("msg <- error")
                return
            }
            if cli.Conn == nil {
                klog.Info("ws is nil, will not send")
                continue
            }
            // 发送消息
            cli.Mutex.Lock()
            err := cli.Conn.WriteMessage(websocket.TextMessage, msg.([]byte))
            cli.Mutex.Unlock()
            if err != nil {
                klog.Info("WriteMessage:", err)
                delListConn(cli.Conn)
                //cli.Connect()
            }
        }       
    }
}

func (cli *WSClient) HandleRecieve() {
    for {
        if cli.Conn == nil {
            klog.Info("ws is nil, will not read")
            continue
        }
        cli.Mutex.Lock()
        _, message, err := cli.Conn.ReadMessage()
        cli.Mutex.Unlock()
        if err != nil {
            klog.Info("read error:", err)
            //cli.Conn.Close()
            delListConn(cli.Conn)
            //cli.Connect()
            //continue
            return
        }

        var recvMsg map[string]interface{}
        err = json.Unmarshal(message, &recvMsg)
        if err != nil {
            klog.Info("Unmarshal msg failed.")
            continue
        }
        
        if recvMsg["op"] == "devstatus" {
            klog.Info("recv devstatus and send back.....")
            cli.Send(message) // 直接返回
        }

        if recvMsg["op"] == "heartbeat" {
            //klog.Info("recv heartbeat.....")
            sendHeartBeat(recvMsg["id"].(string)) // 模拟外部函数调用，已知ID情况
        }
        klog.Info("server recv: ", string(message))
    }
}

func (cli *WSClient) Send(content interface{}) {
    //klog.Print("cli: ", cli)
    if cli != nil {
        cli.ReceiverChan <- content
    }
}

func showList() {
    cnt := 0
    for {
        for idx, v := range gclientList {
            klog.Println("list: ", idx, " ", v)
        }
        
        if cnt % 2 == 0 {
            gclientList["111"].Send([]byte("111111"))
        } else if cnt % 2 == 1 {
            gclientList["222"].Send([]byte("22222"))
        }
        
        cnt += 1
        
        time.Sleep(time.Second * 10)
        
    }
    
}

func sendHeartBeat(id string) {
    allMsg := make(map[string]interface{})
    allMsg["op"] = "heartbeat"
    allMsg["id"] = id
    allMsg["timestamp"] = time.Now().UnixNano() / 1e6

    allMsg["data"] = ""
    allJson, err := json.Marshal(allMsg)
    if err != nil {
        klog.Info("Marshal err: ", err)
        return
    }
    //klog.Info("send heartbeat ...")
    gclientList[id].Send(allJson)
}

type BaseMessag struct {
    Id string `json:"id"`
    Op string `json:"op"`
    Timestamp int64 `json:"timestamp"`
    Data interface{} `json:"data"`
}

var gclientList = make(map[string]*WSClient)

func delListConn(conn *websocket.Conn) {
    for _, v := range gclientList {
        if v.Conn == conn {
            v.Conn.Close() // 是否在此关闭？
            klog.Printf("delete ws conn: %v \n", v.Id)
            delete(gclientList, v.Id)
        }
    }
}

func delListID(id string) {
    for _, v := range gclientList {
        if v.Id == id {
            v.Conn.Close()
            klog.Printf("delete ws id: %v \n", v.Id)
            delete(gclientList, v.Id)
        }
    }
}

func delListAll() {
    for _, v := range gclientList {
        v.Conn.Close() // 是否直接关闭？
        klog.Printf("delete ws client: %v \n", v.Id)
        delete(gclientList, v.Id)
    }
}

func WSHandler(c *gin.Context) {
    id := c.Param("id")

    ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        klog.Println("Upgrade error: ", err)
        return
    }
    klog.Printf("local addr: %v remote addr: %v, id: %v\n", ws.LocalAddr(), ws.RemoteAddr(), id)

    // 保存
    
    client := InitWSServer(id, ws) //&Client{Id: id, Conn: ws}
    gclientList[id] = client

    ws.SetPongHandler(func(string) error {klog.Println("pong handle"); return nil })
    
    // 客户端主动发关闭事件，在此函数中响应，如果没发，则不会响应
    ws.SetCloseHandler(func(code int, text string) error { klog.Println("close handle from ", ws.RemoteAddr()); for _, v := range gclientList { if v.Conn == ws { delete(gclientList, v.Id) } }; message := websocket.FormatCloseMessage(code, "");ws.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second));return nil })

// 如果在本函数中循环处理，则在客户端断开连接后，gin才返回，如果是协程则马上返回
/*
    for {
        //读取ws中的数据
        _, message, err := ws.ReadMessage()
        if err != nil {
            klog.Info("ReadMessage failed")
            for _, v := range gclientList {
                if v.Conn == ws {
                    ws.Close()
                    klog.Printf("delete for: %v \n", v.Id)
                    delete(gclientList, v.Id)
                }
            }
            //return
            break
        }

        klog.Printf("server recv msg: %s", message)

        var recvMsg BaseMessag
        // 如果用map，则直接取元素，但每个是interface{}类型，需要用.(string)之类转换
        //var recvMsg map[string]interface{}
        //klog.Printf("op: %v id: %v timestamp: %v\n", recvMsg["op"], recvMsg["id"], recvMsg["timestamp"])
        err = json.Unmarshal(message, &recvMsg)
        if err != nil {
            klog.Println("Unmarshal msg failed.")
            continue
        }
        
        klog.Printf("op: %v id: %v timestamp: %v\n", recvMsg.Op, recvMsg.Id, recvMsg.Timestamp)
    
        // 判断条件，发送。。。
        if recvMsg.Op == "heartbeat" {
            if setHeartBeat == 1 { // 临时测试用
                sendHeartBeat(recvMsg.Id)
            } else {
                klog.Println("not send heartbeat...")
            }
        }
        
        // 处理，组装
        person := make([]Person, 3)

        for i := 0; i < len(person); i++ {
            person[i].Age = 20+i+1
            person[i].Id = i+1
            person[i].Name = "Late Lee"
            person[i].City = "Shenzhen"
        }
        person[0].Name = "Jim Kent"
        person[1].Name = "Kevin Clark"

        allMsg := make(map[string]interface{})
        allMsg["op"] = "get_time"
        allMsg["id"] = recvMsg.Id
        allMsg["timestamp"] = time.Now().UnixNano() / 1e6
        
        allMsg["data"] = person
        
        allJson, err := json.Marshal(allMsg)
        if err != nil {
            klog.Println("Marshal err: ", err)
            return
        }

        //写入ws数据
        err = ws.WriteMessage(websocket.TextMessage, allJson)
        if err != nil {
            klog.Println("WriteMessage err: ", err)
            break
        }

    }
    */

}