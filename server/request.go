package server

import (
	"bufio"
	"fmt"
	"strings"
)

type Request struct {
	Method string
	Path   string
	Header map[string]string
	Body   *bufio.Reader
}

func ReadRequest(requestReader *bufio.Reader) (*Request, error) {
	request := new(Request)

	// read first line
	reqLine, err := readLine(requestReader)
	if err != nil {
		return nil, err
	}
	if len(reqLine) < 2 {
		return nil, fmt.Errorf("invalid header format")
	}
	reqLines := strings.Split(reqLine, " ")
	request.Method = reqLines[0]
	request.Path = reqLines[1]

	// read header
	header, err := readHeader(requestReader)
	if err != nil {
		return nil, err
	}
	request.Header = header

	fmt.Println((*request))
	return request, nil
}

func readHeader(requestReader *bufio.Reader) (map[string]string, error) {
	header := make(map[string]string)
	for {
		line, err := readLine(requestReader)
		if err != nil {
			return nil, err
		}
		if len(line) == 0 {
			break
		}
		lines := strings.Split(line, ":")
		header[lines[0]] = strings.TrimLeft(lines[1], " \t")
	}
	return header, nil
}

func readLine(requestReader *bufio.Reader) (string, error) {
	out := make([]byte, 0, 1024)
	for {
		line, isMore, err := requestReader.ReadLine()
		if err != nil {
			return string(out), err
		}
		// break if it an empty line
		if len(line) == 0 {
			break
		}
		// append to the currentt line buffer
		out = append(out, line...)
		// break when the line ends
		if !isMore {
			break
		}
	}
	return string(out), nil
}
