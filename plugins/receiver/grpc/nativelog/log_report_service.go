// SiGG-Satellite-Network-SII  //

package nativelog

import (
	"io"
	"time"

	"github.com/apache/skywalking-satellite/plugins/server/grpc"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	logging "skywalking.apache.org/repo/goapi/collect/logging/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const eventName = "grpc-log-event"

type LogReportService struct {
	receiveChannel chan *v1.SniffData
	logging.UnimplementedLogReportServiceServer
}

func (s *LogReportService) Collect(stream logging.LogReportService_CollectServer) error {
	dataList := make([][]byte, 0)
	originalData := grpc.NewOriginalData(nil)
	for {
		err := stream.RecvMsg(originalData)
		if err == io.EOF {
			s.flushLogs(dataList)
			return stream.SendAndClose(&common.Commands{})
		}
		if err != nil {
			s.flushLogs(dataList)
			return err
		}
		dataList = append(dataList, originalData.Content)
	}
}

func (s *LogReportService) flushLogs(dataList [][]byte) {
	if len(dataList) == 0 {
		return
	}
	e := &v1.SniffData{
		Name:      eventName,
		Timestamp: time.Now().UnixNano() / 1e6,
		Meta:      nil,
		Type:      v1.SniffType_Logging,
		Remote:    true,
		Data: &v1.SniffData_LogList{
			LogList: &v1.BatchLogList{
				Logs: dataList,
			},
		},
	}
	s.receiveChannel <- e
}
