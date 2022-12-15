// SiGG-Satellite-Network-SII  //

package grpc

import (
	"fmt"

	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
)

func init() {
	encoding.RegisterCodec(codec{})
}

// OriginalData is keep binary Content
type OriginalData struct {
	Content []byte
}

func NewOriginalData(data []byte) *OriginalData {
	return &OriginalData{Content: data}
}

// codec is overwritten the original "proto" codec, and support using OriginalData to skip data en/decoding.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		original, ok := v.(*OriginalData)
		if !ok {
			return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message or grpc.OriginalData", v)
		}
		return original.Content, nil
	}
	return proto.Marshal(vv)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	vv, ok := v.(proto.Message)
	if !ok {
		original, ok := v.(*OriginalData)
		if !ok {
			return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message or grpc.OriginalData", v)
		}
		original.Content = data
		return nil
	}
	return proto.Unmarshal(data, vv)
}

func (codec) Name() string {
	return "proto"
}
