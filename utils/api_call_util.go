package utils

import (
	"io"
	"net/http"
	"strconv"
)

const (
	base_address = "http://192.168.1.113:8888/summary/"
	size_call    = "/size"
)

func API_Call(instrumentSymbol, year, month string) (int64, error) {
	url := base_address + instrumentSymbol + year + month + size_call
	req, err := http.Get(url)
	if err != nil {
		return -1, err
	}

	defer req.Body.Close()

	size_bytes, err := io.ReadAll(req.Body)
	if err != nil {
		return -1, err
	}
	size_string := string(size_bytes)

	size, err := strconv.ParseInt(size_string, 10, 64)
	if err != nil {
		return -1, err
	}

	return size, nil
}
