package handlers

import (
	"encoding/csv"
	"instrumentFolderSized/utils"
	"net/http"
	"strconv"
)

func CreateCsvHandler(w http.ResponseWriter, r *http.Request) {
	year, month, err := utils.ParseForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	csv := csv.NewWriter(w)
	defer csv.Flush()

	// Headline
	period := year + month
	csv.Write([]string{"Instrument Output for ", period})
	csv.Write([]string{}) //Empty row for spacing

	//Column Headers
	columnHeaders := []string{"Instrument_ID", "Data (in GB)"}
	csv.Write(columnHeaders)

	// Fill rows
	for _, instrument_id := range InstrumentIds {
		size_bytes, err := utils.API_Call(instrument_id, year, month)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		size_int := utils.FromBytesToGib(size_bytes)
		size_str := strconv.FormatInt(size_int, 10)

		csv.Write([]string{instrument_id, size_str})
	}
	if err := csv.Error(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+year+month+".csv\"")
	w.WriteHeader(http.StatusOK)
}
