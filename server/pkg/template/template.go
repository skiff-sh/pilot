package template

import (
	"strings"

	"github.com/flosch/pongo2/v6"
	"github.com/goccy/go-json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type Data map[string]any

func (d Data) ToProto() *structpb.Struct {
	v, _ := structpb.NewStruct(d)
	return v
}

func (d Data) Get(key string) (any, bool) {
	top, key := d.traverse(strings.Split(key, "."), false)
	if top == nil {
		return nil, false
	}

	val, ok := top[key]
	return val, ok
}

func (d Data) Set(key string, val any) {
	d.SetKeys(strings.Split(key, "."), val)
}

func (d Data) SetKeys(dots []string, val any) {
	top, key := d.traverse(dots, true)
	top[key] = val
}

func (d Data) toPongo() pongo2.Context {
	return pongo2.Context(d)
}

func (d Data) traverse(dots []string, create bool) (top Data, key string) {
	switch len(dots) {
	case 0:
		return
	case 1:
		return d, dots[0]
	}

	field := dots[0]

	var temp Data
	q, ok := d[field]
	if !ok {
		if !create {
			return
		}
		temp = Data{}
		d[field] = temp
	} else {
		temp = q.(Data)
	}
	return temp.traverse(dots[1:], create)
}

// Merge merges all Data into a singular where the last
// Data has highest precedence.
func Merge(d ...Data) Data {
	out := Data{}
	for _, v := range d {
		for k, val := range v {
			out[k] = val
		}
	}
	return out
}

// Unmarshal is a generic func that will unmarshal any type into Data.
// If o is a proto.Message, the ProtoMarshaller is used. If o is a
// Data or the underlying type, it is a noop. All others are subject
// to the standard json.Unmarshal func.
func Unmarshal(o any) (Data, error) {
	switch typ := o.(type) {
	case proto.Message:
		return unmarshalProtoMessage(typ)
	case map[string]any:
		return typ, nil
	case Data:
		return typ, nil
	default:
		return unmarshalJSON(o)
	}
}

func IsTruthy(a any) bool {
	if a == nil {
		return false
	}

	switch t := a.(type) {
	case string:
		t = strings.ToLower(strings.TrimSpace(t))
		return t != "" && t != "false" && t != "0" && t != "nil" && t != "null"
	case bool:
		return t
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float64, float32:
		return t != 0
	}
	return false
}

var (
	ProtoMarshaller = &protojson.MarshalOptions{
		AllowPartial: true,
	}

	ProtoUnmarshaller = &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
)

func unmarshalProtoMessage(msg proto.Message) (Data, error) {
	b, err := ProtoMarshaller.Marshal(msg)
	if err != nil {
		return nil, err
	}

	m := map[string]any{}
	return m, json.Unmarshal(b, &m)
}

func unmarshalJSON(msg any) (Data, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	m := map[string]any{}
	return m, json.Unmarshal(b, &m)
}
