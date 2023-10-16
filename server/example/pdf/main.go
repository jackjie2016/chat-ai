package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

func main() {
	// 打开 XLSX 文件
	xlFile, err := xlsx.OpenFile("./example/pdf/www.xlsx")
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}

	// 遍历每个工作表
	for _, sheet := range xlFile.Sheets {
		// 遍历每行数据
		for _, row := range sheet.Rows {
			// 遍历每个单元格
			Cells := row.Cells
			m := make(map[string]string, 0)
			m["id"] = Cells[0].String()    //序号
			m["image"] = Cells[1].String() //序号
			m["info"] = Cells[2].String()  //序号
			fmt.Println(m)
			//for kk, cell := range row.Cells {
			//	// 获取单元格的值
			//	value := cell.String()
			//
			//	fmt.Println(index, kk, value)
			//}
		}
	}
}
