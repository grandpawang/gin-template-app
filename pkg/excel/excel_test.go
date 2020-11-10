package excel

import (
	"fmt"
	"log"
	"os"
	"testing"
)

type excelTestTable struct {
	Name  string
	Email string
	Phone string
}

// TestRead
func TestRead(t *testing.T) {
	fileName := "readTest.xlsx"
	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("can't not open file , %v", err)
	}
	defer file.Close()
	fileStat, _ := file.Stat()
	fileSize := fileStat.Size()
	excelData := make([]byte, fileSize)
	_, err = file.Read(excelData)
	if err != nil {
		log.Printf("file open err, %v", err)
	}

	table, err := ReadFromByte(excelData, excelTestTable{})
	if err != nil {
		log.Println("fail")
	}

	fmt.Println(table)
	fmt.Printf("%T \r\n", table)

	PrintFromByte(excelData, excelTestTable{})
}

// TestReadFile
func TestReadFile(t *testing.T) {
	fileName := "readTest.xlsx"
	table, err := ReadFromFile(fileName, excelTestTable{})
	if err != nil {
		log.Println("fail")
	}

	fmt.Println(table)
}

// TestWrite
func TestWrite(t *testing.T) {
	da := []excelTestTable{}
	da = append(da, excelTestTable{
		Name:  "123",
		Email: "456",
		Phone: "789",
	})
	da = append(da, excelTestTable{
		Name:  "858",
		Email: "485",
		Phone: "666",
	})

	WriteToFile("writeTest.xlsx", da)
}

// TestWriteByte
func TestWriteByte(t *testing.T) {
	da := []excelTestTable{}
	da = append(da, excelTestTable{
		Name:  "123",
		Email: "456",
		Phone: "789",
	})
	da = append(da, excelTestTable{
		Name:  "858",
		Email: "485",
		Phone: "666",
	})

	out, _ := WriteToByte(da)
	file, err := os.OpenFile("writeTestBytes.xlsx", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		t.Error("open fail")
	}
	defer file.Close()
	file.Write(out)
}
