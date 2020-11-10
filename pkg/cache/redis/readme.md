软件包redis是Redis数据库的客户端。

[Redigo常见问题解答](https://github.com/garyburd/redigo/wiki/FAQ) 包含更多内容有关此软件包的文档。

## 连接

连接是Redis主要界面，应用程序通过调用 `Dial`, `DialWithTimeout` or `NewConn`用于创建功能分片和其他类型的连接。

当完成连接时候，应用程序必须在应用程序调用连接时的Close方法，以确保完成连接。

## 执行命令

### Do

> Conn接口具有执行Redis命令的通用方法：
```go
func Do(commandName string, args ...interface{}) (reply interface{}, err error)
```

> [Redis命令参考](http://redis.io/commands) 列出了可用的命令。 使用Redis APPEND命令的示例是：
```go
n, err := conn.Do("APPEND", "key", "value")
```

> Do方法将命令参数转换为二进制字符串以进行传输发送到服务器，如下所示：

| Go Type         | Conversion                          |
|-----------------|-------------------------------------|
| []byte          | Sent as is                          |
| string          | Sent as is                          |
| int, int64      | strconv.FormatInt(v)                |
| float64         | strconv.FormatFloat(v, 'g', -1, 64) |
| bool            | true -> "1", false -> "0"           |
| nil             | ""                                  |
| all other types | fmt.Print(v)                        |


> Redis命令的回复类型使用以下Go类型表示：

| Redis type    | Go type                                    |
|---------------|--------------------------------------------|
| error         | redis.Error                                |
| integer       | int64                                      |
| simple string | string                                     |
| bulk string   | []byte or nil if value not present.        |
| array         | []interface{} or nil if value not present. |


使用类型断言或回复帮助器函数进行转换`interface{}`到命令结果的特定Go类型。

### Pipelining 管道

连接支持使用`Send`，`Flush`和`Receive`方法进行管道传输。
```go
// Send 将命令写入连接的输出缓冲区
func Send(commandName string, args ...interface{}) error
// Flush 将连接的输出缓冲区刷新到服务器
func Flush() error
// Receive 从服务器读取单个答复
func Receive() (reply interface{}, err error)
```

> 以下示例显示了一个简单的管道。
```go
c.Send("SET", "foo", "bar")
c.Send("GET", "foo")
c.Flush()
c.Receive() // reply from SET
v, err = c.Receive() // reply from GET
```

> 使用Send和Do方法来实现流水线事务。

Do方法结合了`Send`，`Flush`和`Receive`的功能方法。 Do方法通过编写命令并刷新输出开始缓冲。 

接下来，Do方法接收所有待处理的答复，包括答复对于Do刚发送的命令。 

如果收到的任何答复是错误的，然后Do返回错误。

如果没有错误，则Do返回最后一个回复。 

如果Do方法的命令参数为""，则Do方法将刷新输出缓冲区并接收未处理的回复，而无需发送命令。

```go
c.Send("MULTI")
c.Send("INCR", "foo")
c.Send("INCR", "bar")
r, err := c.Do("EXEC")
fmt.Println(r) // prints [1, 1]
```

### Concurrency 并发

连接不支持并发调用`write`方法 (Send, Flush)并发调用`read`方法 (Receive).
连接允许并发的`reader`和`writer`。

由于`Do`方法结合了`Send`, `Flush`和`Receive`功能，`Do`方法不能与其他方法同时调用。

要完全并发访问Redis，请使用`thread-safe Pool`获取和从`goroutine`内释放连接。


> Publish and Subscribe (发布和订阅)


使用`Send`，`Flush`和`Receive`方法来实现`Publish` and `Subscribe`。

```go
c.Send("SUBSCRIBE", "example")
c.Flush()
for {
    reply, err := c.Receive()
    if err != nil {
        return err
    }
    // 处理推送消息
}
```

> `PubSubConn`类型用方便的方法包装`Conn`来实现`Subscribe`。

`Subscribe`, `PSubscribe`, `Unsubscribe`, `PUnsubscribe`方法发送并刷新`subscription`管理命令。 

接收方法将推送的消息转换为方便的类型。

```go
psc := redis.PubSubConn{c}
psc.Subscribe("example")
for {
    switch v := psc.Receive().(type) {
    case redis.Message:
        fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
    case redis.Subscription:
        fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
    case error:
        return v
    }
}
```

### Reply Helpers 回复助手

> `Bool`，`Int`，`Bytes`，`String`，`Strings`和`Values`函数转换答复到特定类型的值。

为了方便将呼叫包装到连接`Do`和`Receive`方法，函数采用第二个参数输入错误。
如果错误为非nil，则辅助函数将返回错误。 
如果错误为nil，该函数将回复转换为指定的类型：

```go
exists, err := redis.Bool(c.Do("EXISTS", "foo"))
if err != nil {
    // handle error return from c.Do or type conversion error.
}
```

> 扫描功能将数组回复的元素转换为Go类型：
```go
var value1 int
var value2 string
reply, err := redis.Values(c.Do("MGET", "key1", "key2"))
if err != nil {
    // handle error
}
 if _, err := redis.Scan(reply, &value1, &value2); err != nil {
    // handle error
}
```