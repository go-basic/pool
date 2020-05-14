## 安装
```
go get github.com/go-basic/pool
```
## 实现链接接口
```
type DemoCloser struct {
	Conn     net.Conn
	activeAt time.Time
}

func (p *DemoCloser) Close() error {
	return p.Conn.Close()
}

func (p *DemoCloser) GetActiveTime() time.Time {
	return p.activeAt
}

```
## 使用
```
	p, err := pool.NewGenericPool(2, 5, 30*time.Second, func() (poolable pool.Poolable, e error) {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		return &DemoCloser{Conn: conn, activeAt: time.Now()}, nil
	})
	if err != nil {
		fmt.Println("err=", err)
	}

	//从连接池中取得一个连接
	v, err := p.Get()

	//do something
	//conn=v.(net.Conn)

	//将连接放回连接池中
	p.Put(v)

	//释放连接池中的所有连接
	//p.Release()

	//查看当前连接中的数量
	current := p.Len()
	fmt.Println("len=", current)
```
## 更多见example
https://github.com/go-basic/pool/blob/master/example/main.go