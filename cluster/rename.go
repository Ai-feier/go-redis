package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

// Rename renames a key, the origin and the destination must within the same node
func Rename(cluster *ClusterDatabase, c resp.Connection, cmdAndArgs [][]byte) resp.Reply {
	if len(cmdAndArgs) != 3 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'rename' command")
	}
	src := string(cmdAndArgs[1])
	dest := string(cmdAndArgs[2])

	srcPeer := cluster.peerPicker.PickNode(src)
	destPeer := cluster.peerPicker.PickNode(dest)

	if srcPeer != destPeer {
		return reply.MakeErrReply("ERR rename must within one slot in cluster mode")
	}
	return cluster.relay(srcPeer, c, cmdAndArgs)
}
