package rpc_conn_pool_golang

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"testing"
	"time"
)

const (
	IP_PORT_LISTEN_GOB   = "0.0.0.0:6666"
	IP_PORT_CONNECT_GOB  = "127.0.0.1:6666"
	IP_PORT_LISTEN_JSON  = "0.0.0.0:6667"
	IP_PORT_CONNECT_JSON = "127.0.0.1:6667"
)

type HelloWorld int

func (t *HelloWorld) Hello(req struct{}, resp *string) error {
	*resp = "HelloWorld"
	return nil
}

func init() {
	log.Println("Gob Server Starting")
	log.Println("Json Server Starting")
	go CreateGobServer()
	go CreateJsonServer()
	time.Sleep(time.Second * time.Duration(5))
	log.Println("Gob Server Started")
	log.Println("Json Server Started")
}

func CreateGobServer() {
	SERVER_GOB, err := net.Listen("tcp", IP_PORT_LISTEN_GOB)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	server := rpc.NewServer()
	server.Register(new(HelloWorld))
	for {
		conn, err := SERVER_GOB.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go server.ServeConn(conn)
	}
}

func CreateJsonServer() {
	SERVER_JSON, err := net.Listen("tcp", IP_PORT_LISTEN_JSON)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	server := rpc.NewServer()
	server.Register(new(HelloWorld))
	for {
		conn, err := SERVER_JSON.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func Test_GobRpcConnPool(t *testing.T) {

	factory := func() (net.Conn, error) {
		return net.DialTimeout("tcp", IP_PORT_CONNECT_GOB, time.Second*time.Duration(5))
	}
	log.Println("Start Create Pool")
	rpc_pool, err := NewRpcConnPool("gob", 10000, 100000, factory)
	log.Println("End Create Pool")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer func() {
		rpc_pool.Close()
	}()

	response := ""
	log.Println("Start Get Conn")
	rpc_conn, err := rpc_pool.Get()
	log.Println("End Get Conn")
	if err != nil {
		t.Fatal(err)
		rpc_pool.CloseRpcConn(rpc_conn)
		return
	}
	log.Println("Start Call")
	err = rpc_conn.Call("HelloWorld.Hello", struct{}{}, &response)
	log.Println("End Call")
	if err != nil {
		log.Println("err:", err.Error())
		rpc_pool.CloseRpcConn(rpc_conn)
		return
	}
	log.Println("response:", response)
	rpc_pool.Release(rpc_conn)
}

func Benchmark_GobRpcConnPool(b *testing.B) {
	b.StopTimer()
	log.Println("Sleep 60s Benchmark_GobRpcConnPool")
	time.Sleep(time.Duration(60) * time.Second)

	factory := func() (net.Conn, error) {
		return net.DialTimeout("tcp", IP_PORT_CONNECT_GOB, time.Second*time.Duration(5))
	}
	rpc_pool, err := NewRpcConnPool("gob", 10000, 100000, factory)
	if err != nil {
		b.Fatal(err)
		return
	}
	defer func() {
		rpc_pool.Close()
	}()
	response := ""

	log.Println("Start Benchmark_GobRpcConnPool")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rpc_conn, err := rpc_pool.Get()
		if err != nil {
			rpc_pool.CloseRpcConn(rpc_conn)
			b.Log(err)
			continue
		}
		rpc_conn.Call("HelloWorld.Hello", struct{}{}, &response)
		rpc_pool.Release(rpc_conn)
	}
}

func Benchmark_GobRpcConnPoolWithoutCall(b *testing.B) {
	b.StopTimer()
	log.Println("Sleep 60s Benchmark_GobRpcConnPoolWithoutCall")
	time.Sleep(time.Duration(60) * time.Second)

	factory := func() (net.Conn, error) {
		return net.DialTimeout("tcp", IP_PORT_CONNECT_GOB, time.Second*time.Duration(5))
	}
	rpc_pool, err := NewRpcConnPool("gob", 10000, 100000, factory)
	if err != nil {
		b.Fatal(err)
		return
	}
	defer func() {
		rpc_pool.Close()
	}()

	log.Println("Start Benchmark_GobRpcConnPoolWithoutCall")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rpc_conn, err := rpc_pool.Get()
		if err != nil {
			rpc_pool.CloseRpcConn(rpc_conn)
			continue
		}
		rpc_pool.Release(rpc_conn)
	}
}

func Test_JsonRpcConnPool(t *testing.T) {

	factory := func() (net.Conn, error) {
		return net.DialTimeout("tcp", IP_PORT_CONNECT_JSON, time.Second*time.Duration(5))
	}
	log.Println("Start Create Pool")
	rpc_pool, err := NewRpcConnPool("json", 10000, 100000, factory)
	log.Println("End Create Pool")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer func() {
		rpc_pool.Close()
	}()

	response := ""
	log.Println("Start Get Conn")
	rpc_conn, err := rpc_pool.Get()
	log.Println("End Get Conn")
	if err != nil {
		t.Fatal(err)
		rpc_pool.CloseRpcConn(rpc_conn)
		return
	}
	log.Println("Start Call")
	err = rpc_conn.Call("HelloWorld.Hello", struct{}{}, &response)
	log.Println("End Call")
	if err != nil {
		log.Println("err:", err.Error())
		rpc_pool.CloseRpcConn(rpc_conn)
		return
	}
	log.Println("response:", response)
	rpc_pool.Release(rpc_conn)
}

func Benchmark_JsonRpcConnPool(b *testing.B) {
	b.StopTimer()
	log.Println("Sleep 60s Benchmark_JsonRpcConnPool")
	time.Sleep(time.Duration(60) * time.Second)

	factory := func() (net.Conn, error) {
		return net.DialTimeout("tcp", IP_PORT_CONNECT_JSON, time.Second*time.Duration(5))
	}
	rpc_pool, err := NewRpcConnPool("json", 10000, 100000, factory)
	if err != nil {
		b.Fatal(err)
		return
	}
	defer func() {
		rpc_pool.Close()
	}()
	response := ""

	log.Println("Start Benchmark_JsonRpcConnPool")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rpc_conn, err := rpc_pool.Get()
		if err != nil {
			rpc_pool.CloseRpcConn(rpc_conn)
			b.Log(err)
			continue
		}
		rpc_conn.Call("HelloWorld.Hello", struct{}{}, &response)
		rpc_pool.Release(rpc_conn)
	}
}

func Benchmark_JsonRpcConnPoolWithouCall(b *testing.B) {
	b.StopTimer()
	log.Println("Sleep 60s Benchmark_JsonRpcConnPoolWithouCall")
	time.Sleep(time.Duration(60) * time.Second)

	factory := func() (net.Conn, error) {
		return net.DialTimeout("tcp", IP_PORT_CONNECT_JSON, time.Second*time.Duration(5))
	}
	rpc_pool, err := NewRpcConnPool("json", 10000, 100000, factory)
	if err != nil {
		b.Fatal(err)
		return
	}
	defer func() {
		rpc_pool.Close()
	}()

	log.Println("Start Benchmark_JsonRpcConnPoolWithouCall")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rpc_conn, err := rpc_pool.Get()
		if err != nil {
			rpc_pool.CloseRpcConn(rpc_conn)
			b.Log(err)
			continue
		}
		rpc_pool.Release(rpc_conn)
	}
}
