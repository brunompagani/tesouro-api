package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"strings"
)

func sortRecords(records []Record) {
	sort.Slice(records, func(i, j int) bool {
		if records[i].Nome != records[j].Nome {
			return records[i].Nome < records[j].Nome
		}
		return records[i].DataVencimento < records[j].DataVencimento
	})
}

func writeJSON(records []Record, path string) error {
	tmpPath := path + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(records); err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err := file.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}

func writeCSV(records []Record, path string) error {
	tmpPath := path + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'

	// Write header
	header := []string{"Nome", "Data Inicio", "Data Conversao", "Data Vencimento", "Data Base", "Taxa Compra Manha", "Taxa Venda Manha", "PU Compra Manha", "PU Venda Manha", "PU Base Manha"}
	if err := writer.Write(header); err != nil {
		os.Remove(tmpPath)
		return err
	}

	// Write records
	for _, rec := range records {
		row := []string{
			rec.Nome,
			rec.DataInicio,
			rec.DataConversao,
			rec.DataVencimento,
			rec.DataBase,
			formatFloatBR(rec.TaxaCompraManha),
			formatFloatBR(rec.TaxaVendaManha),
			formatFloatBR(rec.PUCompraManha),
			formatFloatBR(rec.PUVendaManha),
			formatFloatBR(rec.PUBaseManha),
		}
		if err := writer.Write(row); err != nil {
			os.Remove(tmpPath)
			return err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err := file.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}

func formatFloatBR(f float64) string {
	// Format with comma as decimal separator
	s := strconv.FormatFloat(f, 'f', -1, 64)
	s = strings.ReplaceAll(s, ".", ",")
	return s
}
