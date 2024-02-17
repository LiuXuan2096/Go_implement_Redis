package reply

/**
 * Redis向客户端返回的一些固定回复写在这个文件里
 */

// PongReply is +Pong, 向客户端发送回复：pong
type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

// ToBytes
func (r *PongReply) ToBytes() []byte {
	return pongBytes
}

// OkReply is ok，向客户端发送回复：OK
type OkReply struct{}

var okBytes = []byte("+OK\r\n")

func (r *OkReply) ToBytes() []byte {
	return okBytes
}

var theOkReply = new(OkReply)

func MakeOkReply() *OkReply {
	return theOkReply
}

// 向客户端回复null
var nullBulkBytes = []byte("$-1\r\n")

type NullBulkReply struct{}

func (r *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

// 向客户端回复空字符串
var emptyMultiBulkBytes = []byte("*0\r\n")

type EmptyMultiBulkReply struct{}

func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

// NoReply 向客户端发送一个空回复，
type NoReply struct {
}

var noBytes = []byte("")

func (r *NoReply) ToBytes() []byte {
	return noBytes
}
