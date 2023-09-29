package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

// FlushDB removes all data in current database
func FlushDB(cluster *ClusterDatabase, c resp.Connection, cmdAndArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(c, cmdAndArgs)
	var errReply reply.ErrorReply
	for _, r := range replies {
		if reply.IsErrorReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
	}
	if errReply == nil {
		return &reply.OkReply{}
	}
	return reply.MakeErrReply("error occurs: " + errReply.Error())
}
