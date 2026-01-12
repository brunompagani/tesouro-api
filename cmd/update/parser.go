package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseCSV(r io.Reader) (map[string]*assetRecord, error) {
	br := bufio.NewReader(r)
	csvReader := csv.NewReader(br)
	csvReader.Comma = ';'
	csvReader.FieldsPerRecord = -1 // Tolerant

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Handle UTF-8 BOM in first header field
	if len(header) > 0 {
		header[0] = strings.TrimPrefix(header[0], "\ufeff")
	}

	// Validate header
	expectedHeaders := []string{"Tipo Titulo", "Data Vencimento", "Data Base", "Taxa Compra Manha", "Taxa Venda Manha", "PU Compra Manha", "PU Venda Manha", "PU Base Manha"}
	if len(header) < len(expectedHeaders) {
		return nil, fmt.Errorf("insufficient columns in header: got %d, expected %d", len(header), len(expectedHeaders))
	}

	latest := make(map[string]*assetRecord)
	lineNum := 2 // Start at 2 (header is line 1)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading line %d: %w", lineNum, err)
		}

		if len(row) < 8 {
			continue // Skip incomplete rows
		}

		record, err := parseRecord(row)
		if err != nil {
			// Log but continue processing
			fmt.Fprintf(os.Stderr, "Warning: failed to parse line %d: %v\n", lineNum, err)
			lineNum++
			continue
		}

		// Create asset key: Tipo Titulo + Data Vencimento (using internal tipoTitulo)
		assetKey := record.tipoTitulo + "|" + record.DataVencimento

		// Parse data base for comparison
		dataBase, err := time.Parse("2006-01-02", record.DataBase)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse data base on line %d: %v\n", lineNum, err)
			lineNum++
			continue
		}

		// Track latest record per asset key and minimum Data Base (start date)
		if existing, exists := latest[assetKey]; !exists {
			latest[assetKey] = &assetRecord{
				dataBaseMax: dataBase,
				dataBaseMin: dataBase,
				record:      record,
			}
		} else {
			// Update minimum if this record has an older Data Base
			if dataBase.Before(existing.dataBaseMin) {
				existing.dataBaseMin = dataBase
			}
			// Update record if this is a newer Data Base
			if !dataBase.Before(existing.dataBaseMax) {
				existing.dataBaseMax = dataBase
				existing.record = record
			}
		}

		lineNum++
	}

	return latest, nil
}

func parseRecord(row []string) (Record, error) {
	var rec Record
	var err error

	rec.tipoTitulo = strings.TrimSpace(row[0])

	// Parse dates (dd/mm/yyyy -> ISO yyyy-mm-dd)
	rec.DataVencimento, err = parseDate(row[1])
	if err != nil {
		return rec, fmt.Errorf("failed to parse Data Vencimento: %w", err)
	}

	rec.DataBase, err = parseDate(row[2])
	if err != nil {
		return rec, fmt.Errorf("failed to parse Data Base: %w", err)
	}

	// Parse float values (PT-BR format: comma as decimal separator)
	rec.TaxaCompraManha, err = parseFloatBR(row[3])
	if err != nil {
		return rec, fmt.Errorf("failed to parse Taxa Compra Manha: %w", err)
	}

	rec.TaxaVendaManha, err = parseFloatBR(row[4])
	if err != nil {
		return rec, fmt.Errorf("failed to parse Taxa Venda Manha: %w", err)
	}

	rec.PUCompraManha, err = parseFloatBR(row[5])
	if err != nil {
		return rec, fmt.Errorf("failed to parse PU Compra Manha: %w", err)
	}

	rec.PUVendaManha, err = parseFloatBR(row[6])
	if err != nil {
		return rec, fmt.Errorf("failed to parse PU Venda Manha: %w", err)
	}

	rec.PUBaseManha, err = parseFloatBR(row[7])
	if err != nil {
		return rec, fmt.Errorf("failed to parse PU Base Manha: %w", err)
	}

	// Compute combined name: tipo_titulo + year from data_vencimento
	// Extract year from ISO date (yyyy-mm-dd format)
	if len(rec.DataVencimento) >= 4 {
		year := rec.DataVencimento[:4]
		rec.Nome = rec.tipoTitulo + " " + year
	} else {
		rec.Nome = rec.tipoTitulo // Fallback if date parsing failed
	}

	return rec, nil
}

func parseDate(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("empty date string")
	}

	t, err := time.Parse("02/01/2006", s)
	if err != nil {
		return "", err
	}

	return t.Format("2006-01-02"), nil
}

func parseFloatBR(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	// Replace comma with dot, remove thousand separators
	s = strings.ReplaceAll(s, ".", "") // Remove thousand separators
	s = strings.ReplaceAll(s, ",", ".") // Replace comma with dot

	return strconv.ParseFloat(s, 64)
}
