package mgoImport

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataParser struct {
	buf    *bufio.Reader
	DataCh chan []string
	deli   string
}

func InitParser(fileDir string, limit int, deli string) (*DataParser, error) {
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

	return NewDataParser(file, limit, deli), nil
}

func NewDataParser(reader io.Reader, size int, deli string) *DataParser {

	return &DataParser{
		buf:    bufio.NewReader(reader),
		DataCh: make(chan []string, size),
		deli:   deli,
	}
}

func (d *DataParser) readLine() (err error) {
	var counter int
	defer close(d.DataCh)
	for {
		data, _, err := d.buf.ReadLine()
		if err == io.EOF {
			break
		}
		if data == nil || len(data) == 0 {
			continue
		}
		d.DataCh <- strings.Split(string(data), d.deli)
		if counter == 1000 {
			counter = 0
			fmt.Println("已处理1000条记录！")
		}
	}
	return
}
