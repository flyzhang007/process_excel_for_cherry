package main

import (
	"container/list"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

func main() {

	files := []string{}
	path, _ := os.Getwd()
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		name := f.Name()
		listDir := list.New()
		listDir.PushBack(name)
		files = append(files, name)
		return nil
	})

	_, direrr := os.Stat("处理后文件夹")
	if direrr != nil {
		fmt.Printf("文件夹不存在\n")
		os.Mkdir("./处理后文件夹", 0777)
	}

	for i := 0; i < len(files); i++ {
		//name := "江苏高价数据(1735-1795).xlsx"
		name := files[i]
		//查找xlsx文件
		if strings.Contains(name, "xlsx") == false {
			continue
		}
		xlFile, err := xlsx.OpenFile(name)
		if err != nil {
			continue
		}
		newFile := xlsx.NewFile()

		for _, sheet := range xlFile.Sheets {
			fmt.Printf("sheet name:%s\n", sheet.Name)
			if len(sheet.Rows) < 3 {
				fmt.Printf("%s 是空白表格\n", sheet.Name)
				continue
			}
			newSheet, _ := newFile.AddSheet(sheet.Name + "处理后")

			//前两行直接写固定
			//第一行
			newRow1 := newSheet.AddRow()
			newCell1 := newRow1.AddCell()
			newCell1.Value, _ = sheet.Rows[0].Cells[0].String()
			//第二行
			newRow2 := newSheet.AddRow()
			//第二行两列
			newCell2 := newRow2.AddCell()
			newCell2.Value = "年代月份"
			newCell3 := newRow2.AddCell()
			newCell3.Value = "数量"
			//for _, row := range sheet.Rows {
			//剩余的行
			for i := 2; i < len(sheet.Rows); i++ {
				//for _, cell := range sheet.Rows[i].Cells {
				for j := 0; j < len(sheet.Rows[i].Cells); j++ {
					//data, _ := sheet.Rows[i].Cells[j].String()
					if j > 0 {
						innerRow := newSheet.AddRow()
						innerCell1 := innerRow.AddCell()
						//j<=9 拼接成:年份+0+月份,eg 178903[1789年3月]
						if j < 10 {
							//innerCell1.Value, _ = sheet.Rows[i].Cells[0].String()
							//innerCell1.Value += "0" + strconv.Itoa(j)
							value, _ := sheet.Rows[i].Cells[0].String()
							value += "0" + strconv.Itoa(j)
							ivalue, _ := strconv.Atoi(value)
							innerCell1.SetInt(ivalue)
						} else { //j>9 直接拼接成:年份+月,eg 178911
							//innerCell1.Value, _ = sheet.Rows[i].Cells[0].String()
							//innerCell1.Value += strconv.Itoa(j)
							value, _ := sheet.Rows[i].Cells[0].String()
							value += strconv.Itoa(j)
							ivalue, _ := strconv.Atoi(value)
							innerCell1.SetInt(ivalue)
						}
						//具体数据
						innerCell2 := innerRow.AddCell()
						//innerCell2.Value, _ = sheet.Rows[i].Cells[j].String()
						dataV, _ := sheet.Rows[i].Cells[j].String()
						idataV, _ := strconv.Atoi(dataV)
						if idataV != 0 { //==0的空白不填
							innerCell2.SetInt(idataV)
						}
					}
				}
			}
		}
		newFile.Save("./处理后文件夹/处理后" + name)
	}
}
