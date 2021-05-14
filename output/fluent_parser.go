package output

import (
	"reflect"

	"github.com/tinylib/msgp/msgp"
	"github.com/valyala/fastjson"
)

type LogMarshaler struct {
	line *fastjson.Object
	meta map[string]interface{}
}

func (l *LogMarshaler) MarshalMsg(out []byte) ([]byte, error) {
	size := l.line.Len()
	if l.meta != nil {
		size++
	}
	out = msgp.AppendMapHeader(out, uint32(size))
	out, err := parse(out, l.line)
	if err != nil {
		return nil, err
	}
	if l.meta != nil {
		out = msgp.AppendString(out, "meta")
		out = msgp.AppendMapHeader(out, uint32(len(l.meta)))
		for k, v := range l.meta {
			out = msgp.AppendString(out, k)
			t := reflect.TypeOf(v)
			if t == nil {
				out = msgp.AppendNil(out)
			} else {
				switch t.Kind() {
				case reflect.String:
					vv, _ := v.(string)
					out = msgp.AppendString(out, vv)
				default: // FIXME
					out = msgp.AppendNil(out)
				}
			}
		}
	}
	return out, nil
}

func parse(out []byte, o *fastjson.Object) ([]byte, error) {
	var err error
	o.Visit(func(k []byte, v *fastjson.Value) {
		out = msgp.AppendStringFromBytes(out, k)
		switch v.Type() {
		case fastjson.TypeString:
			s, err := v.StringBytes()
			if err != nil {
				return
			}
			out = msgp.AppendString(out, string(s))
		case fastjson.TypeNumber:
			var i int64
			i, err = v.Int64()
			if err != nil {
				return
			}
			out = msgp.AppendInt64(out, i)
		case fastjson.TypeTrue:
			out = msgp.AppendBool(out, true)
		case fastjson.TypeFalse:
			out = msgp.AppendBool(out, false)
		case fastjson.TypeNull:
			out = msgp.AppendNil(out)
		case fastjson.TypeObject:
			var o *fastjson.Object
			o, err = v.Object()
			if err != nil {
				return
			}
			out = msgp.AppendMapHeader(out, uint32(o.Len()))
			out, err = parse(out, o)
		default:
			panic(v.Type())
		}
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
