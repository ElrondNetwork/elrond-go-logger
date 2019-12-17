package mock

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/golang/protobuf/proto"
)

var TestingMarshalizers = map[string]logger.Marshalizer{
	"capnp": &CapnpMarshalizer{},
	"json":  &JsonMarshalizer{},
	"proto": &ProtobufMarshalizer{},
}

// CapnpHelper is an interface that defines methods needed for
// serializing and deserializing Capnp structures into Go structures and viceversa
type CapnpHelper interface {
	// Save saves the serialized data of the implementer type into a stream through Capnp protocol
	Save(w io.Writer) error
	// Load loads the data from the stream into a go structure through Capnp protocol
	Load(r io.Reader) error
}

//-------- capnp

type CapnpMarshalizer struct{}

func (x *CapnpMarshalizer) Marshal(obj interface{}) ([]byte, error) {
	out := bytes.NewBuffer(nil)

	o := obj.(CapnpHelper)
	// set the members to capnp struct
	err := o.Save(out)

	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (x *CapnpMarshalizer) Unmarshal(obj interface{}, buff []byte) error {
	out := bytes.NewBuffer(buff)

	o := obj.(CapnpHelper)
	// set the members to capnp struct
	err := o.Load(out)

	return err
}

func (x *CapnpMarshalizer) IsInterfaceNil() bool {
	return x == nil
}

//-------- Json

type JsonMarshalizer struct{}

func (j JsonMarshalizer) Marshal(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nil, errors.New("NIL object to serilize from!")
	}

	return json.Marshal(obj)
}

func (j JsonMarshalizer) Unmarshal(obj interface{}, buff []byte) error {
	if obj == nil {
		return errors.New("nil object to serilize to")
	}
	if buff == nil {
		return errors.New("nil byte buffer to deserialize from")
	}
	if len(buff) == 0 {
		return errors.New("empty byte buffer to deserialize from")
	}

	return json.Unmarshal(buff, obj)
}

func (j *JsonMarshalizer) IsInterfaceNil() bool {
	return j == nil
}

//------- protobuf

type ProtobufMarshalizer struct{}

func (x *ProtobufMarshalizer) Marshal(obj interface{}) ([]byte, error) {
	if msg, ok := obj.(proto.Message); ok {
		enc, err := proto.Marshal(msg)
		if err != nil {
			return nil, err
		}
		return enc, nil
	}
	return nil, errors.New("can not serialize the object")
}

func (x *ProtobufMarshalizer) Unmarshal(obj interface{}, buff []byte) error {
	if msg, ok := obj.(proto.Message); ok {
		return proto.Unmarshal(buff, msg)
	}
	return errors.New("obj does not implement proto.Message")
}

func (x *ProtobufMarshalizer) IsInterfaceNil() bool {
	return x == nil
}

//------- stub

type MarshalizerStub struct {
	MarshalCalled   func(obj interface{}) ([]byte, error)
	UnmarshalCalled func(obj interface{}, buff []byte) error
}

func (ms *MarshalizerStub) Marshal(obj interface{}) ([]byte, error) {
	return ms.MarshalCalled(obj)
}

func (ms *MarshalizerStub) Unmarshal(obj interface{}, buff []byte) error {
	return ms.UnmarshalCalled(obj, buff)
}

func (ms *MarshalizerStub) IsInterfaceNil() bool {
	return ms == nil
}
