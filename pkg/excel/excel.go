package excel

import (
	"bytes"
	"gin-template-app/pkg/utils"
	"reflect"

	"github.com/modood/table"
	"github.com/tealeg/xlsx"
	"github.com/unknwon/com"
)

// File xlsx file struct
type File = xlsx.File

// Sheet xlsx sheet struct
type Sheet = xlsx.Sheet

// Get sheet's index (First Line)
func Index(sheet *Sheet) map[int]string {
	index := map[int]string{}
	for i, cell := range sheet.Row(0).Cells {
		index[i] = cell.Value
	}
	return index
}

//// Get sheet's index (First Line)
//func IndexStr(sheet *Sheet) map[string]int {
//
//}

// Excel2Slice Get sheet's content
// index -> map[  column  ] name
func Excel2Slice(sheet *Sheet, index map[int]string, obj interface{}, SkipFirstLine bool) (interface{}, error) {
	table := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(obj)), 0, 0)
	for i, row := range sheet.Rows {
		if i == 0 && SkipFirstLine {
			continue
		}
		line := reflect.New(reflect.TypeOf(obj))
		for j, cell := range row.Cells {
			line.Elem().FieldByName(index[j]).SetString(cell.Value)
		}
		table = reflect.Append(table, reflect.ValueOf(line.Elem().Interface()))
	}
	return table.Interface(), nil
}

// Slice2Excel []interface{} -> Sheet
func Slice2Excel(sheet *Sheet, obj interface{}) (*Sheet, error) {
	row := sheet.AddRow()

	nilStruct := utils.Struct1DMapHadNil(reflect.New(reflect.TypeOf(obj).Elem()).Elem().Interface(), nil)
	index := utils.GetKeys(nilStruct)

	for _, v := range index {
		cells := row.AddCell()
		cells.Value = v
	}
	srcValue := reflect.ValueOf(obj)

	for i := 0; i < srcValue.Len(); i++ {
		line := utils.Struct1DMapHadNil(srcValue.Index(i).Interface(), nil)
		row = sheet.AddRow()
		for _, v := range index {
			cell := row.AddCell()
			cell.Value = com.ToStr(line[v])
		}
	}

	return sheet, nil
}

// Read
func Read(file *File, obj interface{}) (interface{}, error) {
	sheet := file.Sheets[0]
	index := Index(sheet)
	return Excel2Slice(sheet, index, obj, true)
}

// ReadFromByte Read from []byte
func ReadFromByte(bs []byte, obj interface{}) (interface{}, error) {
	file, err := xlsx.OpenBinary(bs)
	if err != nil {
		return nil, err
	}

	return Read(file, obj)
}

// ReadFromFile Read from file
func ReadFromFile(fileName string, obj interface{}) (interface{}, error) {
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	return Read(file, obj)
}

// PrintFromByte Print table to output
func PrintFromByte(bs []byte, obj interface{}) {
	list, _ := ReadFromByte(bs, obj)
	table.Output(list)
}

// PrintFromFile Print table to output
func PrintFromFile(fileName string, obj interface{}) {
	list, _ := ReadFromFile(fileName, obj)
	table.Output(list)
}

// WriteToFile Write excel to file
func WriteToFile(fileName string, obj interface{}) error {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("output")

	sheet, _ = Slice2Excel(sheet, obj)
	file.Save(fileName)
	//file.Write()
	return nil
}

// WriteToByte write excel to []byte
func WriteToByte(obj interface{}) ([]byte, error) {
	bs := new(bytes.Buffer)
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("output")

	sheet, _ = Slice2Excel(sheet, obj)
	file.Write(bs)
	//fmt.Println(bs)

	return bs.Bytes(), nil
}
