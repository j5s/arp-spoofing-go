package Action

import "github.com/andlabs/ui"

type Table struct {
	Table  *ui.Table
	Model  *ui.TableModel
	HModel *modelHandler
}

func MakeHostTable() *Table {
	hModel := NewModelHandler()
	hostModel := ui.NewTableModel(hModel)
	hostTable := ui.NewTable(&ui.TableParams{
		Model:                         hostModel,
		RowBackgroundColorModelColumn: 3, //-1 默认颜色
	})
	hostTable.AppendTextColumn("IP", 0, 0, nil)
	hostTable.AppendTextColumn("MAC", 1, 0, nil)
	hostTable.AppendTextColumn("MACInfo", 2, 0, nil)
	hostTable.AppendButtonColumn("IsGateway", 3, ui.TableModelColumnAlwaysEditable)
	hostTable.AppendButtonColumn("Spooling", 4, ui.TableModelColumnAlwaysEditable)
	hostTable.AppendButtonColumn("PacketType", 5, ui.TableModelColumnAlwaysEditable)
	hostTable.AppendButtonColumn("Cut Off", 6, ui.TableModelColumnAlwaysEditable)

	return &Table{
		Table:  hostTable,
		Model:  hostModel,
		HModel: hModel,
	}
}
