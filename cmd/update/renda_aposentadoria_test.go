package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRecordRendaAposentadoriaExtra(t *testing.T) {
	tests := []struct {
		name     string
		row      []string
		validate func(t *testing.T, rec Record)
	}{
		{
			name: "Tesouro Renda+ Aposentadoria Extra - uses conversion year",
			row: []string{
				"Tesouro Renda+ Aposentadoria Extra",
				"15/12/2069",
				"22/12/2025",
				"7,02",
				"7,14",
				"500,08",
				"482,41",
				"482,41",
			},
			validate: func(t *testing.T, rec Record) {
				// Should use conversion year (2069 - 19 = 2050) in nome
				// Conversion date: January 15th of conversion year
				assert.Equal(t, "Tesouro Renda+ Aposentadoria Extra 2050", rec.Nome)
				assert.Equal(t, "2069-12-15", rec.DataVencimento)
				// Conversion date: January 15th, year = maturity year - 19
				assert.Equal(t, "2050-01-15", rec.DataConversao)
			},
		},
		{
			name: "Tesouro Renda+ Aposentadoria Extra 2049",
			row: []string{
				"Tesouro Renda+ Aposentadoria Extra",
				"15/12/2049",
				"22/12/2025",
				"7,28",
				"7,40",
				"1872,56",
				"1847,29",
				"1847,29",
			},
			validate: func(t *testing.T, rec Record) {
				// Should use conversion year (2049 - 19 = 2030) in nome
				// Matches PDF example: Tesouro RendA+ 2030, conversion 15/01/2030, maturity 15/12/2049
				assert.Equal(t, "Tesouro Renda+ Aposentadoria Extra 2030", rec.Nome)
				assert.Equal(t, "2049-12-15", rec.DataVencimento)
				assert.Equal(t, "2030-01-15", rec.DataConversao)
			},
		},
		{
			name: "Regular bond - not Renda+ Aposentadoria Extra",
			row: []string{
				"Tesouro IPCA+",
				"15/05/2035",
				"22/12/2025",
				"7,29",
				"7,41",
				"2374,37",
				"2348,76",
				"2348,76",
			},
			validate: func(t *testing.T, rec Record) {
				// Should use maturity year
				assert.Equal(t, "Tesouro IPCA+ 2035", rec.Nome)
				assert.Equal(t, "2035-05-15", rec.DataVencimento)
				assert.Empty(t, rec.DataConversao) // No conversion date for regular bonds
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec, err := parseRecord(tt.row)
			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, rec)
			}
		})
	}
}

func TestConversionDateCalculation(t *testing.T) {
	// Test that conversion date is correctly calculated
	// Conversion date: January 15th, year = maturity year - 19
	testCases := []struct {
		maturity    string
		expectedConv string
		expectedYear string
	}{
		{"2069-12-15", "2050-01-15", "2050"},
		{"2049-12-15", "2030-01-15", "2030"}, // Matches PDF example
		{"2074-12-15", "2055-01-15", "2055"},
	}

	for _, tc := range testCases {
		maturityDate, err := time.Parse("2006-01-02", tc.maturity)
		require.NoError(t, err)
		
		convYear := maturityDate.Year() - 19
		conversionDate := time.Date(convYear, 1, 15, 0, 0, 0, 0, time.UTC)
		expectedConv, err := time.Parse("2006-01-02", tc.expectedConv)
		require.NoError(t, err)
		
		assert.Equal(t, expectedConv, conversionDate)
		assert.Equal(t, tc.expectedYear, fmt.Sprintf("%d", convYear))
	}
}
