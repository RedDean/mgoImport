package mgoImport

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type DataParser struct {
	buf *bufio.Reader
}

func NewDataParser(reader io.Reader) *DataParser {
	return &DataParser{
		buf: bufio.NewReader(reader),
	}
}

func (d DataParser) readLine() (ret []string, err error) {
	var data []byte

	for {
		data, _, err = d.buf.ReadLine()
		if err == io.EOF {
			break
		}
		ret = append(ret, string(data))
	}

	return ret, err
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
