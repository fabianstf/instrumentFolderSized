package handlers

import (
	"bytes"
	"github.com/jung-kurt/gofpdf"
	"instrumentFolderSized/utils"
	"net/http"
	"strconv"
)

func CreatePdfHandler(w http.ResponseWriter, r *http.Request) {
	year, month, err := utils.ParseForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	pdf := CreatePDFTempl()

	for _, instrument_id := range InstrumentIds {
		size_bytes, err := utils.API_Call(instrument_id, year, month)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		size_int := utils.FromBytesToGib(size_bytes)
		size_str := strconv.FormatInt(size_int, 10)

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

func CreatePDFTempl() gofpdf.Pdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Instr_ID")
	pdf.Cell(40, 10, "Data (in GB)")
	pdf.Ln(15)
	return pdf
}
