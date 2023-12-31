package cluster

import (
	"context"
	"errors"
	"go-redis/interface/resp"
	"go-redis/lib/utils"
	"go-redis/resp/client"
	"go-redis/resp/reply"
	"strconv"
)

func (cluster *ClusterDatabase) getPeerClient(peer string) (*client.Client, error) {
	factory, ok := cluster.peerConnection[peer]  // 拿到连接池连接
	if !ok {
		return nil, errors.New("connection factory not found")
	}
	raw, err := factory.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}
	conn, ok := raw.(*client.Client)
	if !ok {
		return nil, errors.New("connection factory make wrong type")
	}
	return conn, nil
}

// relay relays command to peer
// select db by c.GetDBIndex()
// cannot call Prepare, Commit, execRollback of self node
func (cluster *ClusterDatabase) relay(peer string, c resp.Connection, args [][]byte) resp.Reply {
	if peer == cluster.self {
		return cluster.db.Exec(c, args)
	}
	peerClient, err := cluster.getPeerClient(peer)
	if err != nil {
		return reply.MakeErrReply(err.Error())
	}
	defer func() {
		_ = cluster.returnPeerClient(peer, peerClient)
	}()
	// select db
	peerClient.Send(utils.ToCmdLine("SELECT", strconv.Itoa(c.GetDBIndex())))
	return peerClient.Send(args)
}

func (cluster *ClusterDatabase) returnPeerClient(peer string, peerClient *client.Client) error {
	factory, ok := cluster.peerConnection[peer]
	if !ok {
		return errors.New("connection factory not found")
	}
	return factory.ReturnObject(context.Background(), peerClient)
}

// broadcast broadcasts command to all node in cluster
func (cluster *ClusterDatabase) broadcast(c resp.Connection, args [][]byte) map[string]resp.Reply {
	res := make(map[string]resp.Reply)
	for _, node := range cluster.nodes {
		r := cluster.relay(node, c, args)
		res[node] = r 
	}
	return res
}
