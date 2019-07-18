package mgoImport

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

type DataParser struct {
	buf *bufio.Reader
	DataCh chan string
}

func InitParser(filedir string, limit int) *DataParser {
	file, err := os.OpenFile(dir, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	return NewDataParser(file, limit)
}

func NewDataParser(reader io.Reader,size int) *DataParser {
	return &DataParser{
		buf: bufio.NewReader(reader),
		DataCh: make(chan string, size),
	}
}

func (d *DataParser) readLine()(err error) {
	defer close(d.DataCh)
	for {
		data, _, err := d.buf.ReadLine()
		if err == io.EOF {
			break
		}
		d.DataCh <- string(data)
	}
	return
}

func splitByDelimiter(str string, deli string) ([]string, error) {

	if str == "" {
		return nil, nil
	}

	if deli == "" {
		return nil, errors.New("delimiter is empty string")
	}

	if len([]byte(deli)) > 1 {
		return nil, errors.New("delimiter size greater than one byte")
	}

	return strings.Split(str, deli), nil
}
