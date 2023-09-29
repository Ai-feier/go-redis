package reply

type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func (p *PongReply) ToBytes() []byte {
	return pongBytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

type OkReply struct {
}

var okBytes = []byte("+OK\r\n")

var theOkReply = new(OkReply)

func (o *OkReply) ToBytes() []byte {
	return okBytes
}

func MakeOkReply() *OkReply {
	return theOkReply
}

type NullBulkReply struct {
}

var nullBulkReply = []byte("$-1\r\n")

var theNullBulkReply = new(NullBulkReply)

func (n *NullBulkReply) ToBytes() []byte {
	return nullBulkReply
}

func MakeNullBulkReply() *NullBulkReply {
	return theNullBulkReply
}

type EmptyBulkReply struct {
}

var emptyBulkBytes = []byte("$0\r\n")

var theEmptyBulkReply = new(EmptyBulkReply)

func (n *EmptyBulkReply) ToBytes() []byte {
	return emptyBulkBytes
}

func MakeEmptyBulkReply() *EmptyBulkReply {
	return theEmptyBulkReply
}

type EmptyMultiBulkReply struct {
}

var (
	emptyMultiBulkBytes = []byte("*0\r\n")
	theEmptyMultiBulkReply = new(EmptyMultiBulkReply)
)

func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return theEmptyMultiBulkReply
}

type NoReply struct {
}

var noReplyBytes = []byte("")

func (n *NoReply) ToBytes() []byte {
	return noReplyBytes
}






