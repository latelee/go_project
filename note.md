```
一般性的post请求处理
curl http://127.0.0.1:9000/testing -X POST -H "Content-Type:application/json" -d  '{"id":"test_001", "op":"etc", "timestamp":12342134341234, "data":{"name":"foo", "addr":"bar", "code":450481, "age":100}}'

// 发送空 或单纯字符
curl http://127.0.0.1:9000/testing -X POST -d "abc"

// 发送文件
curl http://127.0.0.1:9000/testing -X POST -F "file=@foo.json"
```