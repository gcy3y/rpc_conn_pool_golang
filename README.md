
# A Handy, Hign-Performace, Extensible Golang's Rpc Conn Pool 

* Based on github.com/fatih/pool Conn Pool Currently
* Suppprt JsonRpc and GobRpc(golang default rpc) Currently

## Install and Usage

Install the package with:

```bash
go get github.com/gcy3y/rpc_conn_pool_golang
```

Please vendor the package with one of the releases: https://github.com/gcy3y/rpc_conn_pool_golang/releases.
`master` branch is **development** branch and will contain always the latest changes.


## Example

>> create a factory() to be used with rpc pool
```
factory := func() (net.Conn, error) {
   return net.DialTimeout("tcp", 801, time.Second*time.Duration(5))
}
```


>> create a new json rpc pool with an initial capacity of 1000 and maximum
>> capacity of 10000. The factory will create 10000 initial rpc connections and put it
>> into the pool. rpc type can be "json" or "gob" or "", "" means gob rpc

```
rpc_pool, err := NewRpcConnPool("json", 1000, 10000, factory)
```

>> now you can get a rpc connection from the pool, if there is no connection
>> available it will create a new one via the factory function.
```
rpc_conn, err := p.Get()
```

>> now you can get use this rpc connection call rpc method
```
err = rpc_conn.Call("HelloWorld.Hello", struct{}{}, &response)
```

>> do something with conn and put it back to the pool by closing the rpc connection
>> (this doesn't close the underlying rpc connection instead it's putting it back
>> to the pool).

```
rpc_pool.Release(rpc_conn)
```

>> close the underlying rpc connection instead of returning it to pool
```
rpc_pool.CloseRpcConn(rpc_conn)
```

>> close pool any time you want, this closes all the rpc connections inside a pool
```
rpc_pool.Close()
```

>> currently available connections in the rpc pool
```
current_len := rpc_pool.Len()
```


## Credits
 * [gcy3y](https://github.com/gcy3y)

## License

The MIT License (MIT) - see LICENSE for more details
