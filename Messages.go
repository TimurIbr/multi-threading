package multi_threading

import (
	"container/heap"
	"fmt"
	"sync"
)

type bytevector []byte

type messageArgType int

const (
	intType messageArgType = iota
	int64Type
	stringType
	vectorIntType
	eofType
)

func (mtype messageArgType) String() string {
	names := [...]string{
		"intType",
		"int64Type",
		"stringType",
		"vectorIntType",
		"eofType",
	}
	if mtype < intType || mtype > eofType {
		return "Unknown"
	}
	return names[mtype]
}

type messageArg struct {
	//enum {
	//        IntType='A', Int64Type, StringType, VectorIntType, EOFType
	//    };
	MessageType messageArgType //look up name
	Body        bytevector
}

type Message struct {
	SendTime, DeliveryTime int64
	From, To, Ptr          int // from = -1, to = -1, ptr = 0
	Body                   bytevector
}

func MakeMessageFromArg(mArgs ...messageArg) Message        {}
func MakeMessage(from int, to int, body bytevector) Message {}
func (Message) GetString() string                           { return "" }
func (Message) GetInt() int                                 { return 0 }
func (Message) GetInt64() int64                             { return 0 }
func (Message) More(oth Message) bool                       { return false }
func (Message) append(a messageArg)                         {}

type priority_queue smth
type MessageQueue struct {
	//TODO(low) implement priority_queue throught heap see more https://golang.org/pkg/container/heap/
	queue priority_queue // priority_queue queue <Message, vector<Message>, greater<Message> >
	//TODO(hard): implement recursuveness
	_mutex sync.Mutex //recursive_mute
}
