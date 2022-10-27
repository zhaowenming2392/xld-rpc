package excel

import (
	//中文文档 https://xuri.me/excelize/zh-hans/
	excelize "github.com/xuri/excelize/v2"
)

// 将映射保存到新的表格文件
//
// 如果文件存在，会覆盖掉此文件
//
// data map[string][][]interface{} 表名称、每行或每列的值，支持多种类型值
//
// excel文件称作一个工作簿，内部的每个子表称为工作表
//
// filename 路径和名称，不带后缀名
//
// writeLine 是否按行写，按行map[string][行][列]interface{}，否则map[string][列][行]interface{}
//
// opt ...excelize.Options 写入选项
func DataToExcel(data map[string][][]interface{}, filename string, writeLine bool, opt ...excelize.Options) (err error) {
	isSheet1 := true
	//在创建新的工作簿时，将包含名为 Sheet1 的默认工作表。
	f := excelize.NewFile()
	defer func() {
		// Close the spreadsheet.
		err = f.Close()
	}()

	for sheet, sheetDatas := range data {
		if isSheet1 {
			//是第一张表，但不是默认的表名称，改名
			if sheet != "Sheet1" {
				f.SetSheetName("Sheet1", sheet)
			}
			isSheet1 = false
		} else {
			//创建新的表
			f.NewSheet(sheet)
		}

		//填充内容
		var key string
		for aID, as := range sheetDatas {
			for bID, value := range as {
				//获取坐标
				if writeLine {
					key, err = excelize.CoordinatesToCellName(bID+1, aID+1)
				} else {
					key, err = excelize.CoordinatesToCellName(aID+1, bID+1)
				}

				if err != nil {
					return err
				}

				//单元格插入：表，坐标名如A1，值
				f.SetCellValue(sheet, key, value)
			}
		}
	}

	if err = f.SaveAs(filename+".xlsx", opt...); err != nil {
		return err
	}

	return nil
}

// 将映射追加或覆盖到现有表格文件
//
// 如果某单元格存在则是覆盖内容，否则就是追加
//
// data map[string][][]interface{} 表名称、每行或每列的值，支持多种类型值
//
// filename 路径和名称，不带后缀名，打开和保存的都是同一个名称
//
// writeLine 是否按行写，按行map[string][行][列]interface{}，否则map[string][列][行]interface{}
//
// opt ...excelize.Options 写入选项，打开和写入的都是同一个选项
func DataAppendToExcel(data map[string][][]interface{}, filename string, writeLine bool, opt ...excelize.Options) (err error) {
	//打开源表
	f, err := excelize.OpenFile(filename, opt...)
	if err != nil {
		return err
	}

	defer func() {
		// Close the spreadsheet.
		err = f.Close()
	}()

	//现有表
	sheets := f.GetSheetList()

	for sheet, sheetDatas := range data {
		//检测是否存在此表名
		var haveSheet bool
		for _, tabSheet := range sheets {
			if sheet == tabSheet {
				haveSheet = true
			}
		}
		//不存在创建新表
		if !haveSheet {
			//创建新的表
			f.NewSheet(sheet)
		}

		//填充内容
		var key string
		for aID, as := range sheetDatas {
			for bID, value := range as {
				//获取坐标
				if writeLine {
					key, err = excelize.CoordinatesToCellName(bID+1, aID+1)
				} else {
					key, err = excelize.CoordinatesToCellName(aID+1, bID+1)
				}

				if err != nil {
					return err
				}

				//单元格插入：表，坐标名如A1，值
				f.SetCellValue(sheet, key, value)
			}
		}
	}

	if err = f.SaveAs(filename+".xlsx", opt...); err != nil {
		return err
	}

	return nil
}

// 将表格文件读出数据
//
// filename 表格文件路径及其名称和后缀
//
// 以行读取：map[表名][行][列]原生string
// 以列读取：map[表名][列][行]原生string
//
// opt ...excelize.Options 打开选项：excelize.Options{Password: "password"}使用密码
func ExcelToData(filename string, readLine bool, opt ...excelize.Options) (data *map[string][][]string, err error) {
	f, err := excelize.OpenFile(filename, opt...)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		err = f.Close()
	}()

	data = &map[string][][]string{}
	//所有工作表
	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		var res [][]string
		if readLine {
			//按行获取
			res, err = f.GetRows(sheet)
		} else {
			//按列获取
			res, err = f.GetCols(sheet)
		}
		//赋值
		(*data)[sheet] = res

		if err != nil {
			return nil, err
		}

		/*
			data[sheet] = map[string]string{}
			rows, err := f.GetRows(sheet)
			for r := 0; r < len(rows); r++ {
				//每行，r=行号
				for c := 0; c < len(rows[r]); c++ {
					//行中的每个单元格,c=列号
					//fmt.Print(rows[r][c], "\t")

					//将 [X, Y] 形式的列、行索引转换为由字母和数字组合而成的单元格坐标：A1
					key, err := excelize.CoordinatesToCellName(c+1, r+1)
					fmt.Print("[", c+1, r+1, "]", key, "-", rows[r][c], "\t")
					if err != nil {
						fmt.Println(err)
						return nil, err
					}
					//以单元格坐标作为key如A1,C25
					data[sheet][key] = rows[r][c]
				}
				fmt.Println()
			}
		*/
	}

	return data, nil
}
