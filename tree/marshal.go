//
// marshal.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package tree

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func Marshal(v interface{}) ([]byte, error) {
	out := new(bytes.Buffer)

	err := marshalValue(out, reflect.ValueOf(v))
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func marshalValue(out io.Writer, value reflect.Value) error {
	var buf [8]byte

	if !value.IsValid() {
		fmt.Printf("Invalid value %v\n", value)
		return nil
	}

	switch value.Type().Kind() {

	case reflect.Uint8:
		buf[0] = uint8(value.Uint())
		_, err := out.Write(buf[:1])
		return err

	case reflect.Int:
		binary.BigEndian.PutUint32(buf[:4], uint32(value.Int()))
		_, err := out.Write(buf[:4])
		return err

	case reflect.Uint32:
		binary.BigEndian.PutUint32(buf[:4], uint32(value.Uint()))
		_, err := out.Write(buf[:4])
		return err

	case reflect.Uint64:
		binary.BigEndian.PutUint64(buf[:8], value.Uint())
		_, err := out.Write(buf[:8])
		return err

	case reflect.Int64:
		binary.BigEndian.PutUint64(buf[:8], uint64(value.Int()))
		_, err := out.Write(buf[:8])
		return err

	case reflect.Slice:
		binary.BigEndian.PutUint32(buf[:4], uint32(value.Len()))
		_, err := out.Write(buf[:4])
		if err != nil {
			return err
		}
		if value.Type().Elem().Kind() == reflect.Uint8 {
			_, err = out.Write(value.Bytes())
			return err
		} else {
			for i := 0; i < value.Len(); i++ {
				if err := marshalValue(out, value.Index(i)); err != nil {
					return err
				}
			}
			return nil
		}

	case reflect.String:
		data := []byte(value.String())
		binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
		_, err := out.Write(buf[:4])
		if err != nil {
			return err
		}
		_, err = out.Write(data)
		return err

	case reflect.Ptr:
		return marshalValue(out, reflect.Indirect(value))

	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			err := marshalValue(out, value.Field(i))
			if err != nil {
				return err
			}
		}
		return nil

	default:
		return fmt.Errorf("Unsupported type: %s", value.Type().Kind().String())
	}
}
