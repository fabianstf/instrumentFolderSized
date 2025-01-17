package handlers

import (
	"bytes"
	"github.com/jung-kurt/gofpdf"
	"io"
	"net/http"
	"strconv"
)

const (
	base_address = "http://192.168.1.113:8888/summary/"
	size_call    = "/size"
)

var (
	instrument_ids = []string{"F", "G", "I", "M", "O", "V"}
)

func CreatePdfHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	year := r.Form.Get("year")
	month := r.Form.Get("month")

	pdf := CreatePDFTempl()

	for _, instrument_id := range instrument_ids {
		size_bytes, err := API_Call(instrument_id, year, month)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		size_int := fromBytesToGib(size_bytes)
		size_str := strconv.FormatInt(size_int, 10)
		if len(size_str) > 6 {
			size_str = size_str[:len(size_str)-9]
		}

		pdf.Cell(40, 10, instrument_id+": ")
		pdf.Cell(40, 10, size_str)
		pdf.Ln(10)
	}

	// Output PDF to byte slice
	var buf bytes.Buffer
	err = pdf.Output(&buf)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+year+month+".pdf\"")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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

func fromBytesToGib(num int64) int64 {
	return (num + 500_000_000) / 1_000_000_000 * 1_000_000_000
}

func CreatePDFTempl() gofpdf.Pdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Instr_ID")
	pdf.Cell(40, 10, "Data (in GiB)")
	pdf.Ln(15)
	return pdf
}
