package mgoImport

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

type DataParser struct {
	//buf    *bufio.Reader
	buf    *csv.Reader
	DataCh chan []string
	deli   string
}

func InitParser(fileDir string, limit int, deli string, readerSize int) (*DataParser, error) {
	file, err := os.OpenFile(fileDir, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}

	if deli == "" {
		return nil, errors.New("delimiter is empty string")
	}

	if len([]byte(deli)) > 1 {
		return nil, errors.New("delimiter size greater than one byte")
	}

	return NewDataParser(file, limit, deli, readerSize), nil
}

func NewDataParser(reader io.Reader, size int, deli string, readerSize int) *DataParser {

	r := csv.NewReader(bufio.NewReaderSize(reader, readerSize))
	r.Comma = '\t'
	r.LazyQuotes = true
	r.FieldsPerRecord = 0

	return &DataParser{
		buf:    r,
		DataCh: make(chan []string, size),
		deli:   deli,
	}
}

func (d *DataParser) readLine() (err error) {
	var counter int
	defer close(d.DataCh)
	for {
		data, err := d.buf.Read()
		if err == io.EOF {
			break
		}
		if len(data) == 0 {
			continue
		}

		if err != nil {
			fmt.Println("[ERROR] csv parse failed! err:   ", err, " data:", data)
			continue
		}

		d.DataCh <- data

		if counter == 1000 {
			counter = 0
			fmt.Println("[INFO] 已处理1000条记录！")
		}
		counter++
	}
	return
}
