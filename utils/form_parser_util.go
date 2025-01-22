package utils

import (
	"errors"
	"net/http"
)

func ParseForm(r *http.Request) (string, string, error) {
	if r.Method != http.MethodPost {
		return "", "", errors.New("only POST method is allowed")
	}

	err := r.ParseForm()
	if err != nil {
		return "", "", err
	}

	year := r.Form.Get("year")
	month := r.Form.Get("month")

	return year, month, nil
}
