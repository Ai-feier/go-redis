package cluster

import (
	"context"
	"fmt"
	pool "github.com/jolestar/go-commons-pool/v2"
	"go-redis/config"
	"go-redis/database"
	databaseface "go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/lib/consistenthash"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
	"runtime/debug"
	"strings"
)

// CmdFunc represents the handler of a redis command
type CmdFunc func(cluster *ClusterDatabase, c resp.Connection, cmdAndArgs [][]byte) resp.Reply

var router = makeRouter()

// ClusterDatabase represents a node of go-redis cluster
// it holds part of data and coordinates other nodes to finish transactions
type ClusterDatabase struct {
	self string
	
	nodes []string
	peerPicker     *consistenthash.NodeMap  // 一致性哈希
	peerConnection map[string]*pool.ObjectPool  // 连接池
	db             databaseface.Database  // 当前节点下的单机redis
}

// MakeClusterDatabase creates and starts a node of cluster
func MakeClusterDatabase() *ClusterDatabase {
	cluster := &ClusterDatabase{
		self:           config.Properties.Self,

		peerPicker:     consistenthash.NewNodeMap(nil),
		peerConnection: make(map[string]*pool.ObjectPool),
		db:             database.NewStandaloneDatabase(),
	}

	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	nodes = append(nodes, cluster.self)
	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}
	cluster.nodes = nodes
	// 将节点映射到一致性哈希环
	cluster.peerPicker.AddNode(nodes...)

	ctx := context.Background()

	for _, peer := range config.Properties.Peers {
		cluster.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{
			Peer: peer,
		})
	}
	return cluster
}

func (cluster *ClusterDatabase) Close() {
	cluster.db.Close()
}

func (cluster *ClusterDatabase) Exec(client resp.Connection, args [][]byte)(res resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
			res = &reply.UnknownErrReply{}
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return reply.MakeErrReply("ERR unknown command '" + cmdName + "', or not supported in cluster mode")
	}
	res = cmdFunc(cluster, client, args)
	return 
}

func (cluster *ClusterDatabase) AfterClientClose(c resp.Connection) {
	cluster.db.AfterClientClose(c)
}
















