package Action

import (
	"fmt"
	"log"
	"time"

	"ARPSpoofing/Socket"

	"github.com/andlabs/ui"
)

type modelHandler struct {
	Data         []Socket.HostItem
	GatewayIndex int
	TimeTicker   <-chan time.Time
}

func NewModelHandler() *modelHandler {
	m := new(modelHandler)
	m.Data = []Socket.HostItem{}
	m.GatewayIndex = -1
	log.Println("启动定时器")
	m.TimeTicker = time.Tick(1 * time.Second)
	return m
}

func (this *modelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
	}
}

//定义table的行数
func (this *modelHandler) NumRows(m *ui.TableModel) int {
	return len(this.Data)
}

//从0开始计数
func (this *modelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	switch column {
	case 0:
		return ui.TableString(this.Data[row].IP.String())
	case 1:
		return ui.TableString(this.Data[row].MAC.String())
	case 2:
		return ui.TableString(this.Data[row].MACInfo)
	case 3:
		if row == this.GatewayIndex {
			return ui.TableString("Yes")
		} else {
			return ui.TableString("No")
		}
	case 4:
		return ui.TableString(this.Data[row].Spooling)
	case 5:
		return ui.TableString(this.Data[row].PacketType)
	case 6:
		if this.Data[row].IsCutOff {
			return ui.TableString("NoCut")
		} else {
			return ui.TableString("CutOff")
		}

	default:
		return nil
	}

}
func (mh *modelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
	fmt.Println(row, column, value)
	switch column {
	case 3:

		lastGatewayIndex := mh.GatewayIndex
		if lastGatewayIndex == -1 {
			mh.GatewayIndex = row
			m.RowChanged(row)
		} else {
			if lastGatewayIndex != row {
				mh.GatewayIndex = row
				m.RowChanged(lastGatewayIndex)
				m.RowChanged(row)
			}
		}

	case 4:
		if mh.Data[row].Spooling == "Host" {
			mh.Data[row].Spooling = "Gateway"
		} else {
			mh.Data[row].Spooling = "Host"
		}
		m.RowChanged(row)
	case 5:
		if mh.Data[row].PacketType == "Reply" {
			mh.Data[row].PacketType = "Request"
		} else {
			mh.Data[row].PacketType = "Reply"
		}
		m.RowChanged(row)
	case 6:
		fmt.Println("Clicked the button cutoff")
		fmt.Printf("%+v\n", mh.Data[row])
		if mh.GatewayIndex == -1 {
			log.Println("请先选择网关")
			return
		}
		mh.Data[row].IsCutOff = !mh.Data[row].IsCutOff
		m.RowChanged(row)
	}
}
