package rpc_conn_pool_golang

import (
	"errors"
	"github.com/fatih/pool"
	"io"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"
)

const (
	DEFAULT_RPC_TYPE = "gob"
	JSON_RPC_TYPE    = "json"
	GOB_RPC_TYPE     = "gob"
)

type RpcConnPoolI interface {
	Get() (*rpc.Client, error)
	Release(*rpc.Client)
	CloseRpcConn(*rpc.Client)
	Close()
	Len() int
}

type Rpc_driver func(conn io.ReadWriteCloser) *rpc.Client

var (
	driver_lock     sync.RWMutex
	pool_driver_map = make(map[string]Rpc_driver)
)

func init() {
	pool_driver_map[DEFAULT_RPC_TYPE] = rpc.NewClient
	pool_driver_map[GOB_RPC_TYPE] = rpc.NewClient
	pool_driver_map[JSON_RPC_TYPE] = jsonrpc.NewClient
}

func getDriver(rpc_type string) Rpc_driver {
	driver_lock.RLock()
	v, ok := pool_driver_map[rpc_type]
	if ok {
		return v
	}
	driver_lock.RUnlock()
	return nil
}

func Register(rpc_type string, driver Rpc_driver) {
	driver_lock.RLock()
	_, ok := pool_driver_map[rpc_type]
	if ok {
		return
	}
	driver_lock.RUnlock()
	driver_lock.Lock()
	pool_driver_map[rpc_type] = driver
	driver_lock.Unlock()
}

type RpcConnPool struct {
	conn_pool pool.Pool
	driver    Rpc_driver // init once

	conn_rpc_map      map[*rpc.Client]net.Conn
	conn_rpc_map_lock *sync.RWMutex
}

func (t *RpcConnPool) Get() (*rpc.Client, error) {
	conn, err := t.conn_pool.Get()
	if err != nil {
		return nil, err
	}
	rpc_client := t.driver(conn)
	t.conn_rpc_map_lock.Lock()
	t.conn_rpc_map[rpc_client] = conn
	t.conn_rpc_map_lock.Unlock()
	return rpc_client, nil
}

func (t *RpcConnPool) Release(rpc_conn *rpc.Client) {
	t.conn_rpc_map_lock.RLock()
	conn, ok := t.conn_rpc_map[rpc_conn]
	t.conn_rpc_map_lock.RUnlock()
	if ok {
		if pc, ok := conn.(*pool.PoolConn); ok {
			pc.Close()
		}
	}
}

func (t *RpcConnPool) CloseRpcConn(rpc_conn *rpc.Client) {
	t.conn_rpc_map_lock.RLock()
	conn, ok := t.conn_rpc_map[rpc_conn]
	t.conn_rpc_map_lock.RUnlock()
	if ok {
		if pc, ok := conn.(*pool.PoolConn); ok {
			pc.MarkUnusable()
			pc.Close()
		}
	}
}

func (t *RpcConnPool) Close() {
	t.conn_pool.Close()
	t.driver = nil
	t.conn_pool = nil
}

func (t *RpcConnPool) Len() int {
	return t.conn_pool.Len()
}

func NewRpcConnPool(rpc_type string, initialCap, maxCap int, factory pool.Factory) (*RpcConnPool, error) {
	if factory == nil {
		factory = func() (net.Conn, error) { return net.DialTimeout("tcp", "127.0.0.1:4000", 5*time.Second) }
	}
	if rpc_type == "" {
		rpc_type = GOB_RPC_TYPE
	}
	dr := getDriver(rpc_type)
	if dr == nil {
		return nil, errors.New("Unsupported Rpc Type")
	}
	conn_pool, err := pool.NewChannelPool(initialCap, maxCap, factory)
	if err != nil {
		return nil, err
	}
	rpc_pool := RpcConnPool{
		conn_pool:         conn_pool,
		driver:            dr,
		conn_rpc_map_lock: &sync.RWMutex{},
		conn_rpc_map:      make(map[*rpc.Client]net.Conn),
	}
	return &rpc_pool, nil
}
