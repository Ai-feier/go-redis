package cluster

import "go-redis/interface/resp"

func ping(cluster *ClusterDatabase, c resp.Connection, cmdAndArgs [][]byte) resp.Reply {
	// 本地db
	return cluster.db.Exec(c, cmdAndArgs)
}
