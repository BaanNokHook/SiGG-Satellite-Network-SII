// SiGG-Satellite-Network-SII  //

package nativemeter

import (
	"io"
	"time"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	meter "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const eventName = "grpc-nativemeter-event"

type MeterService struct {
	receiveChannel chan *v1.SniffData
	meter.UnimplementedMeterReportServiceServer
}

func (m *MeterService) Collect(stream meter.MeterReportService_CollectServer) error {
	dataList := make([]*meter.MeterData, 0)
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			m.flushMeter(dataList)
			return stream.SendAndClose(&common.Commands{})
		}
		if err != nil {
			m.flushMeter(dataList)
			return err
		}
		dataList = append(dataList, item)
	}
}

func (m *MeterService) flushMeter(dataList []*meter.MeterData) {
	if len(dataList) == 0 {
		return
	}
	d := &v1.SniffData{
		Name:      eventName,
		Timestamp: time.Now().UnixNano() / 1e6,
		Meta:      nil,
		Type:      v1.SniffType_MeterType,
		Remote:    true,
		Data: &v1.SniffData_MeterCollection{
			MeterCollection: &meter.MeterDataCollection{MeterData: dataList},
		},
	}
	m.receiveChannel <- d
}

func (m *MeterService) CollectBatch(batch meter.MeterReportService_CollectBatchServer) error {
	for {
		item, err := batch.Recv()
		if err == io.EOF {
			return batch.SendAndClose(&common.Commands{})
		}
		if err != nil {
			return err
		}

		d := &v1.SniffData{
			Name:      eventName,
			Timestamp: time.Now().UnixNano() / 1e6,
			Meta:      nil,
			Type:      v1.SniffType_MeterType,
			Remote:    true,
			Data: &v1.SniffData_MeterCollection{
				MeterCollection: item,
			},
		}
		m.receiveChannel <- d
	}
}
