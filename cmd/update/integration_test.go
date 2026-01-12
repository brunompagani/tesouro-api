package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	// Sample CSV data that mimics the real format
	sampleCSV := `Tipo Titulo;Data Vencimento;Data Base;Taxa Compra Manha;Taxa Venda Manha;PU Compra Manha;PU Venda Manha;PU Base Manha
Tesouro IPCA+;15/05/2035;22/12/2024;7,29;7,41;2374,37;2348,76;2348,76
Tesouro IPCA+;15/05/2035;23/12/2024;7,30;7,42;2375,00;2349,00;2349,00
Tesouro IPCA+;15/05/2035;22/12/2025;7,31;7,43;2376,00;2350,00;2350,00
Tesouro Selic;17/03/2010;22/12/2025;0,00;0,03;3185,95;3183,49;3182,12
Tesouro Prefixado;01/01/2008;22/12/2025;11,29;11,32;986,50;986,47;986,05`

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "tesouro_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Parse CSV
	reader := strings.NewReader(sampleCSV)
	latest, err := parseCSV(reader)
	require.NoError(t, err)
	require.Len(t, latest, 3) // Three unique bonds (IPCA+ 2035, Selic 2010, Prefixado 2008 Jan and Apr are different)

	// Convert to records and set DataInicio
	records := make([]Record, 0, len(latest))
	for _, asset := range latest {
		asset.record.DataInicio = asset.dataBaseMin.Format("2006-01-02")
		records = append(records, asset.record)
	}
	sortRecords(records)

	// Verify records
	assert.Len(t, records, 3)

	// Find Tesouro IPCA+ 2035
	var ipcaRecord *Record
	for i := range records {
		if records[i].Nome == "Tesouro IPCA+ 2035" {
			ipcaRecord = &records[i]
			break
		}
	}
	require.NotNil(t, ipcaRecord, "Should have Tesouro IPCA+ 2035 record")
	assert.Equal(t, "2035-05-15", ipcaRecord.DataVencimento)
	assert.Equal(t, "2025-12-22", ipcaRecord.DataBase) // Latest
	assert.Equal(t, "2024-12-22", ipcaRecord.DataInicio) // Oldest (start date)
	assert.InDelta(t, 7.31, ipcaRecord.TaxaCompraManha, 0.01)

	// Test JSON output
	jsonPath := filepath.Join(tmpDir, "test.json")
	err = writeJSON(records, jsonPath)
	require.NoError(t, err)

	// Verify JSON file exists and is valid
	jsonData, err := os.ReadFile(jsonPath)
	require.NoError(t, err)

	var decodedRecords []Record
	err = json.Unmarshal(jsonData, &decodedRecords)
	require.NoError(t, err)
	assert.Len(t, decodedRecords, 3)

	// Verify JSON structure
	assert.Equal(t, ipcaRecord.Nome, decodedRecords[0].Nome)
	assert.Equal(t, ipcaRecord.DataInicio, decodedRecords[0].DataInicio)
	assert.Equal(t, ipcaRecord.DataVencimento, decodedRecords[0].DataVencimento)
	assert.Equal(t, ipcaRecord.DataBase, decodedRecords[0].DataBase)

	// Test CSV output
	csvPath := filepath.Join(tmpDir, "test.csv")
	err = writeCSV(records, csvPath)
	require.NoError(t, err)

	// Verify CSV file exists and has correct structure
	csvData, err := os.ReadFile(csvPath)
	require.NoError(t, err)

	csvContent := string(csvData)
	lines := strings.Split(strings.TrimSpace(csvContent), "\n")
	require.GreaterOrEqual(t, len(lines), 2, "CSV should have at least header and one data row")

	// Check header
	assert.Contains(t, lines[0], "Nome")
	assert.Contains(t, lines[0], "Data Inicio")
	assert.Contains(t, lines[0], "Data Vencimento")
	assert.Contains(t, lines[0], "Data Base")

	// Check that header uses semicolons
	assert.Contains(t, lines[0], ";")

	// Verify atomic write (file should exist, .tmp should not)
	tmpPath := csvPath + ".tmp"
	_, err = os.Stat(tmpPath)
	assert.True(t, os.IsNotExist(err), "Temporary file should not exist after write")
}

func TestIntegrationWithRealisticData(t *testing.T) {
	// Test with multiple records for same bond to verify latest + min tracking
	csv := `Tipo Titulo;Data Vencimento;Data Base;Taxa Compra Manha;Taxa Venda Manha;PU Compra Manha;PU Venda Manha;PU Base Manha
Tesouro Prefixado;01/01/2008;01/01/2020;10,00;10,10;900,00;899,00;899,00
Tesouro Prefixado;01/01/2008;02/01/2020;10,05;10,15;901,00;900,00;900,00
Tesouro Prefixado;01/01/2020;15/12/2024;11,00;11,10;950,00;949,00;949,00
Tesouro Prefixado;01/01/2008;31/12/2024;10,20;10,30;905,00;904,00;904,00`

	tmpDir, err := os.MkdirTemp("", "tesouro_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	reader := strings.NewReader(csv)
	latest, err := parseCSV(reader)
	require.NoError(t, err)

	// Should have 2 unique bonds (different maturity dates)
	assert.Len(t, latest, 2)

	// Check Tesouro Prefixado 2008
	key2008 := "Tesouro Prefixado|2008-01-01"
	asset2008, exists := latest[key2008]
	require.True(t, exists)

	// Should have latest record from 2024-12-31
	assert.Equal(t, "2024-12-31", asset2008.record.DataBase)
	assert.InDelta(t, 10.20, asset2008.record.TaxaCompraManha, 0.01)

	// Start date should be the oldest (2020-01-01)
	expectedMin, _ := time.Parse("2006-01-02", "2020-01-01")
	assert.Equal(t, expectedMin, asset2008.dataBaseMin)
}
