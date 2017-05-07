package net

import (
	"github.com/golang/protobuf/proto"
	"log"
	"reflect"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/carsonsx/tcptester/conf"
	"github.com/carsonsx/tcptester/util"
)

type Parser interface {
	Register(args ...interface{}) error
	SetHandler(h func(interface{}), args ...interface{}) error
	Unmarshal(data []byte) (v interface{}, err error)
	Marshal(v interface{}) (data []byte, err error)
}

type StringParser struct {
	handler func(interface{})
}

func (p *StringParser) Register(args ...interface{}) error {
	//do nothing
	return nil
}

func (p *StringParser) SetHandler(h func(interface{}), args ...interface{}) error {
	p.handler = h
	return nil
}

func (p *StringParser) Unmarshal(data []byte) (v interface{}, err error) {
	v = string(data)
	//fmt.Print(v)
	if p.handler != nil {
		p.handler(v)
	}
	err = nil
	return
}

func (p *StringParser) Marshal(v interface{}) (data []byte, err error){
	data = []byte(v.(string))
	err = nil
	return
}

type ProtobufParser struct {
	id_type_map map[uint32]reflect.Type
	type_id_map map[reflect.Type]uint32
	handlers    map[reflect.Type]func(interface{})
	handler func(interface{})
}

func (p *ProtobufParser) Register(args ...interface{}) error {
	if p.id_type_map == nil {
		p.id_type_map = make(map[uint32]reflect.Type)
		p.type_id_map = make(map[reflect.Type]uint32)
	}
	_id := args[0].(uint32)
	_type := args[1].(reflect.Type)
	if _, ok := p.id_type_map[_id]; ok {
		return errors.New(fmt.Sprintf("id[%d] has beend registed", _id))
	}
	p.id_type_map[_id] = _type
	p.type_id_map[_type] = _id
	log.Printf("reistered id[%d] and type[%v]", _id, _type)
	return nil
}

func (p *ProtobufParser) SetHandler(h func(v interface{}), args ...interface{}) error {

	if len(args) == 0 {
		p.handler = h
	} else {
		_msg := args[0]
		_type := reflect.TypeOf(_msg)
		if p.handlers == nil {
			p.handlers = make(map[reflect.Type]func(interface{}))
		}
		if _, ok := p.handlers[_type]; ok {
			return errors.New(fmt.Sprintf("type[%v] has beend registered", _type))
		}
		p.handlers[_type] = h
	}
	return nil
}

func (p *ProtobufParser) Unmarshal(data []byte) (v interface{}, err error) {
	if p.handlers == nil {
		return
	}

	var id uint32
	if conf.Config.ProtoBufferIdSize == 1 {
		id = uint32(data[0])
	} else if conf.Config.ProtoBufferIdSize == 2 {
		if conf.Config.LittleEndian {
			id = uint32(binary.LittleEndian.Uint16(data))
		} else {
			id = uint32(binary.BigEndian.Uint16(data))
		}
	} else if conf.Config.ProtoBufferIdSize == 4 {
		if conf.Config.LittleEndian {
			id = binary.LittleEndian.Uint32(data)
		} else {
			id = binary.BigEndian.Uint32(data)
		}
	}

	if t, ok := p.id_type_map[id]; ok {
		v = reflect.New(t.Elem()).Interface()
		err = proto.UnmarshalMerge(data[conf.Config.ProtoBufferIdSize:], v.(proto.Message))
		if err == nil {
			log.Printf("rcvd data: %v", v)
			if p.handler != nil {
				p.handler(v)
			}
			if handler, ok := p.handlers[t]; ok {
				handler(v)
			}
		}
	}
	return
}

func (p *ProtobufParser) Marshal(v interface{}) (data []byte, err error) {
	if id, ok := p.type_id_map[reflect.TypeOf(v)]; ok {
		rawData, _ := proto.Marshal(v.(proto.Message))
		if conf.Config.ProtoBufferIdSize == 1 {
			data = util.AddUint8ToBytePrefix(rawData, uint8(id))
		} else if conf.Config.ProtoBufferIdSize == 2 {
			data = util.AddUint16ToBytePrefix(rawData, uint16(id), conf.Config.LittleEndian)
		} else if conf.Config.ProtoBufferIdSize == 4 {
			data = util.AddUint32ToBytePrefix(rawData, uint32(id), conf.Config.LittleEndian)
		}
	} else {
		err = errors.New(fmt.Sprintf("not found id of %v", reflect.TypeOf(v)))
	}
	return
}