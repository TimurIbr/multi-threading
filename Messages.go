package multi_threading

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"sync"
)

//type bytevector []byte

type messageArgType int

func (mAT messageArgType) Read(p []byte) (n int, err error) {
	p[0] = byte(mAT)
	return 1, io.EOF ////////////////////
}

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
	order := binary.LittleEndian
	//body := &newMA.Body
	body := new(bytes.Buffer)
	switch pq := q.(type) {
	default:
		fmt.Errorf("makeMessageArg: unknown arg type %T", pq)
	case int16:
		tmpq := pq
		newMA.MessageType = int16Type
		err := binary.Write(body, order, tmpq)
		if err != nil {
			fmt.Errorf("makeMessageArg: binary.Write failed: %v", err)
		}
	case int64:
		tmpq := pq
		newMA.MessageType = int64Type
		err := binary.Write(body, order, tmpq)
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
		tmpq := pq
		newMA.MessageType = vectorIntType
		err := binary.Write(body, order, tmpq)
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
	if _, err := newMA.Body.ReadFrom(newMA.MessageType); err != nil {
		fmt.Errorf("makeMessageArg: newMA.Body.ReadFrom failed: %v", err)
	}
	if err := binary.Write(&newMA.Body, order, body.Bytes()); err != nil {
		fmt.Errorf("makeMessageArg: binary.Write failed: %v", err)
	}
	return newMA
}

type Message struct {
	SendTime, DeliveryTime int64
	From, To, Ptr          int          // from = -1, to = -1, ptr = 0
	Body                   bytes.Buffer //bytevector
}

func MakeMessageFromArg(mArgs ...messageArg) Message {
	ms := Message{}
	for _, arg := range mArgs {
		ms.append(arg)
	}
	ms.From = -1
	ms.To = -1
	ms.Ptr = 0
	ms.DeliveryTime = 0
	return ms
}
func MakeMessage(from int, to int, body bytes.Buffer) Message {
	return Message{From: from, To: to, Body: body}
}
func (ms *Message) GetString() (ret string) {
	if (ms.Ptr < ms.Body.Len()) && (messageArgType(ms.Body.Bytes()[ms.Ptr]) == stringType) {
		ms.Ptr += ms.Body.Len()
		ms.Body.ReadByte()
		ret = ms.Body.String()
		ms.Body.UnreadByte()
		return ret
	}
	return ""
}
func (ms *Message) GetInt16() (ret int16) {
	fmt.Println(messageArgType(ms.Body.Bytes()[ms.Ptr]))
	if (ms.Ptr+2 < ms.Body.Len()) && (messageArgType(ms.Body.Bytes()[ms.Ptr]) == int16Type) {
		for i, val := range ms.Body.Bytes()[1:] {
			ret |= int16(val) << uint16(i*8)
		}
		ms.Ptr += ms.Body.Len()
		return ret
	}

	return 0
}
func (ms *Message) GetInt64() (ret int64) {
	fmt.Println(messageArgType(ms.Body.Bytes()[ms.Ptr]))
	if (ms.Ptr+8 < ms.Body.Len()) && (messageArgType(ms.Body.Bytes()[ms.Ptr]) == int64Type) {
		for i, val := range ms.Body.Bytes()[1:] {
			ret |= int64(val) << uint64(i*8)
		}
		ms.Ptr += ms.Body.Len()
		return ret
	}

	return 0
}
func (ms *Message) GetInt() (ret int) {
	return int(ms.GetInt64())
}
func (ms Message) More(oth Message) bool {
	return ms.DeliveryTime > oth.DeliveryTime
}
func (ms *Message) append(a messageArg) {
	if _, err := ms.Body.ReadFrom(&a.Body); err != nil {
		fmt.Errorf("message.append: failed to read from &a.Body %v", err)
	}
}

type priority_queue smth
type MessageQueue struct {
	//TODO(low) implement priority_queue throught heap see more https://golang.org/pkg/container/heap/
	queue priority_queue // priority_queue queue <Message, vector<Message>, greater<Message> >
	//TODO(hard): implement recursuveness
	_mutex sync.Mutex //recursive_mute
}
