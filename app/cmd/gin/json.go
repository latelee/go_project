/*
本文返回json格式
*/

package gin

import (
    // "fmt"
    "strconv"

    //"k8s.io/klog"

    "github.com/gin-gonic/gin"
    //"net/http"
)

// 用户信息
type Person struct {
    Id   int
    Name string
    City string
    Age int
}

// http://127.0.0.1:4000/api/v1/userinfo/
func FetchAllUsers(c *gin.Context) {

/*
   c.JSON(http.StatusOK, gin.H{
            "message": "hello",
            "location": "hell",
            "time": "long long ago",
        })
*/

    person := make([]Person, 3)
    
    for i := 0; i < len(person); i++ {
        person[i].Age = 20+i+1
        person[i].Id = i+1
        person[i].Name = "Late Lee"
        person[i].City = "Shenzhen"
    }
    person[0].Name = "Jim Kent"
    person[1].Name = "Kevin Clark"
    
    c.JSON(200, gin.H {
        "result": person,
        "count":  1,
        })
    
}

var setHeartBeat = 1

// http://127.0.0.1:4000/api/v1/userinfo/1
// http://127.0.0.1:4000/api/v1/userinfo/250
func FetchSingleUser(c *gin.Context) {
    id := c.Param("id")

    var (
        person Person
        result gin.H
    )
    nid, _ := strconv.Atoi(id)

    // 临时测试
    if nid == 0 {
        setHeartBeat = 0
    } else if nid == 1 {
        setHeartBeat = 1
    }
    
    // 理论上是查询，此处从简，直接赋值，根据ID区别
    person.Name = "Late Lee"
    person.Age = 33
    person.Id = nid
    person.City = "Beijing"
    if nid == 250 {
        result = gin.H {
        "result": nil,
        "count":  0,
        }
    } else {
    result = gin.H {
        "result": person,
        "count":  1,
        }
    }

    c.JSON(200, result)
}