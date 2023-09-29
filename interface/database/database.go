package database

import "go-redis/interface/resp"

/*
	redis 业务核心
 */

type CmdLine = [][]byte 

type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply  // 执行指令
	Close()
	AfterClientClose(c resp.Connection)
}

type DataEntity struct {
	Data interface{}
}

