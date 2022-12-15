// SiGG-Satellite-Network-SII  //

package grpc

import (
	"fmt"
	"reflect"

	"google.golang.org/grpc"

	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

type CommonGRPCReceiverFields struct {
	Server        *grpc.Server
	OutputChannel chan *v1.SniffData // The channel is to bridge the LogReportService and the Gatherer to delivery the data.
}

// InitCommonGRPCReceiverFields init the common fields for gRPC receivers.
func InitCommonGRPCReceiverFields(server interface{}) *CommonGRPCReceiverFields {
	s, ok := server.(*grpc.Server)
	if !ok {
		panic(fmt.Errorf("registerHandler does not support %s", reflect.TypeOf(server).String()))
	}
	return &CommonGRPCReceiverFields{
		Server:        s,
		OutputChannel: make(chan *v1.SniffData),
	}
}
