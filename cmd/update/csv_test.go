package main

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCSV(t *testing.T) {
	tests := []struct {
		name     string
		csv      string
		validate func(t *testing.T, result map[string]*assetRecord)
		wantErr  bool
	}{
		{
			name: "single bond, single record",
			csv: `Tipo Titulo;Data Vencimento;Data Base;Taxa Compra Manha;Taxa Venda Manha;PU Compra Manha;PU Venda Manha;PU Base Manha
Tesouro IPCA+;15/05/2035;22/12/2025;7,29;7,41;2374,37;2348,76;2348,76`,
			validate: func(t *testing.T, result map[string]*assetRecord) {
				require.Len(t, result, 1)
				key := "Tesouro IPCA+|2035-05-15"
				asset, exists := result[key]
				require.True(t, exists)
				assert.Equal(t, "Tesouro IPCA+ 2035", asset.record.Nome)
				assert.Equal(t, "2035-05-15", asset.record.DataVencimento)
				assert.Equal(t, "2025-12-22", asset.record.DataBase)
				assert.Empty(t, asset.record.DataConversao) // Not a Renda+ Aposentadoria Extra bond
				// Start date should be same as latest for single record
				expectedDate, _ := time.Parse("2006-01-02", "2025-12-22")
				assert.Equal(t, expectedDate, asset.dataBaseMin)
				assert.Equal(t, expectedDate, asset.dataBaseMax)
			},
			wantErr: false,
		},
		{
			name: "single bond, multiple records - keeps latest",
			csv: `Tipo Titulo;Data Vencimento;Data Base;Taxa Compra Manha;Taxa Venda Manha;PU Compra Manha;PU Venda Manha;PU Base Manha
Tesouro IPCA+;15/05/2035;22/12/2024;7,29;7,41;2374,37;2348,76;2348,76
Tesouro IPCA+;15/05/2035;23/12/2024;7,30;7,42;2375,00;2349,00;2349,00
Tesouro IPCA+;15/05/2035;22/12/2025;7,31;7,43;2376,00;2350,00;2350,00`,
			validate: func(t *testing.T, result map[string]*assetRecord) {
				require.Len(t, result, 1)
				key := "Tesouro IPCA+|2035-05-15"
				asset, exists := result[key]
				require.True(t, exists)
				// Should keep the latest record (2025-12-22)
				assert.Equal(t, "2025-12-22", asset.record.DataBase)
				assert.InDelta(t, 7.31, asset.record.TaxaCompraManha, 0.01)
				// Start date should be the oldest (2024-12-22)
				expectedMin, _ := time.Parse("2006-01-02", "2024-12-22")
				expectedMax, _ := time.Parse("2006-01-02", "2025-12-22")
				assert.Equal(t, expectedMin, asset.dataBaseMin)
				assert.Equal(t, expectedMax, asset.dataBaseMax)
			},
			wantErr: false,
		},
		{
			name: "multiple bonds",
			csv: `Tipo Titulo;Data Vencimento;Data Base;Taxa Compra Manha;Taxa Venda Manha;PU Compra Manha;PU Venda Manha;PU Base Manha
Tesouro IPCA+;15/05/2035;22/12/2025;7,29;7,41;2374,37;2348,76;2348,76
Tesouro Selic;17/03/2010;22/12/2025;0,00;0,03;3185,95;3183,49;3182,12`,
			validate: func(t *testing.T, result map[string]*assetRecord) {
				require.Len(t, result, 2)
				// Check both bonds exist
				assert.Contains(t, result, "Tesouro IPCA+|2035-05-15")
				assert.Contains(t, result, "Tesouro Selic|2010-03-17")
			},
			wantErr: false,
		},
		{
			name: "handles BOM in header",
			csv:  "\ufeffTipo Titulo;Data Vencimento;Data Base;Taxa Compra Manha;Taxa Venda Manha;PU Compra Manha;PU Venda Manha;PU Base Manha\nTesouro IPCA+;15/05/2035;22/12/2025;7,29;7,41;2374,37;2348,76;2348,76",
			validate: func(t *testing.T, result map[string]*assetRecord) {
				require.Len(t, result, 1)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.csv)
			result, err := parseCSV(reader)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestSortRecords(t *testing.T) {
	records := []Record{
		{Nome: "Tesouro IPCA+ 2035", DataVencimento: "2035-05-15"},
		{Nome: "Tesouro Selic 2010", DataVencimento: "2010-03-17"},
		{Nome: "Tesouro IPCA+ 2035", DataVencimento: "2035-05-20"},
		{Nome: "Tesouro Prefixado 2008", DataVencimento: "2008-01-01"},
	}

	sortRecords(records)

	// Should be sorted by Nome, then DataVencimento
	assert.Equal(t, "Tesouro IPCA+ 2035", records[0].Nome)
	assert.Equal(t, "2035-05-15", records[0].DataVencimento)
	assert.Equal(t, "Tesouro IPCA+ 2035", records[1].Nome)
	assert.Equal(t, "2035-05-20", records[1].DataVencimento)
	assert.Equal(t, "Tesouro Prefixado 2008", records[2].Nome)
	assert.Equal(t, "Tesouro Selic 2010", records[3].Nome)
}
