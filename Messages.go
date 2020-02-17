package multi_threading

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
)

//type bytevector []byte

type messageArgType int

const (
	_ messageArgType = iota
	int16Type
	int64Type
	stringType
	vectorIntType
	eofType
)

func (mtype messageArgType) String() string {
	names := [...]string{
		"Unknown",
		"int16Type",
		"int64Type",
		"stringType",
		"vectorIntType",
		"eofType",
	}
	if mtype < int16Type || mtype > eofType {
		return "Unknown"
	}
	return names[mtype]
}

type messageArg struct {
	//enum {
	//        IntType='A', Int64Type, StringType, VectorIntType, EOFType
	//    };
	MessageType messageArgType //look up name
	Body        bytes.Buffer   //bytevector
}

func makeMessageArg(q interface{}) (newMA messageArg) {
	order := binary.BigEndian
	body := &newMA.Body
	switch pq := q.(type) {
	default:
		fmt.Errorf("makeMessageArg: unknown arg type %T", pq)
	case int16:
		newMA.MessageType = int16Type
		err := binary.Write(body, order, pq)
		if err != nil {
			fmt.Errorf("makeMessageArg: binary.Write failed: %v", err)
		}
	case int64:
		newMA.MessageType = int64Type
		err := binary.Write(body, order, pq)
		if err != nil {
			fmt.Errorf("makeMessageArg: binary.Write failed: %v", err)
		}
	case int:
		tmpq := int64(pq)
		newMA.MessageType = int64Type
		err := binary.Write(body, order, tmpq)
		if err != nil {
			fmt.Errorf("makeMessageArg: binary.Write failed: %v", err)
		}
	case string:
		tmpq := make([]int8, len(pq))
		for i, v := range pq {
			tmpq[i] = int8(v)
		}
		newMA.MessageType = stringType
		err := binary.Write(body, order, tmpq)
		if err != nil {
			fmt.Errorf("makeMessageArg: binary.Write failed: %v", err)
		}
	case []int64:
		newMA.MessageType = vectorIntType
		err := binary.Write(body, order, pq)
		if err != nil {
			fmt.Errorf("makeMessageArg: binary.Write failed: %v", err)
		}
	case []int:
		tmpq := make([]int64, len(pq))
		for i, v := range pq {
			tmpq[i] = int64(v)
		}
		newMA.MessageType = vectorIntType
		err := binary.Write(body, order, tmpq)
		if err != nil {
			fmt.Errorf("makeMessageArg: binary.Write failed: %v", err)
		}
	}
	return newMA
}

type Message struct {
	SendTime, DeliveryTime int64
	From, To, Ptr          int          // from = -1, to = -1, ptr = 0
	Body                   bytes.Buffer //bytevector
}

func MakeMessageFromArg(mArgs ...messageArg) Message          {}
func MakeMessage(from int, to int, body bytes.Buffer) Message {}
func (Message) GetString() string                             { return "" }
func (Message) GetInt() int                                   { return 0 }
func (Message) GetInt64() int64                               { return 0 }
func (Message) More(oth Message) bool                         { return false }
func (Message) append(a messageArg)                           {}

type priority_queue smth
type MessageQueue struct {
	//TODO(low) implement priority_queue throught heap see more https://golang.org/pkg/container/heap/
	queue priority_queue // priority_queue queue <Message, vector<Message>, greater<Message> >
	//TODO(hard): implement recursuveness
	_mutex sync.Mutex //recursive_mute
}
