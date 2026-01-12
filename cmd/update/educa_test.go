package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRecordEducaPlus(t *testing.T) {
	tests := []struct {
		name     string
		row      []string
		validate func(t *testing.T, rec Record)
	}{
		{
			name: "Tesouro Educa+ - uses conversion year",
			row: []string{
				"Tesouro Educa+",
				"15/12/2034",
				"22/12/2025",
				"5,36", "5,48", "2587,63", "2556,12", "2556,12",
			},
			validate: func(t *testing.T, rec Record) {
				// Should use conversion year (2034 - 4 = 2030) in nome
				// Matches PDF example: Tesouro Educa+ 2030, conversion 15/01/2030, maturity 15/12/2034
				assert.Equal(t, "Tesouro Educa+ 2030", rec.Nome)
				assert.Equal(t, "2034-12-15", rec.DataVencimento)
				assert.Equal(t, "2030-01-15", rec.DataConversao)
			},
		},
		{
			name: "Tesouro Educa+ 2039",
			row: []string{
				"Tesouro Educa+",
				"15/12/2039",
				"22/12/2025",
				"5,40", "5,52", "2600,00", "2580,00", "2580,00",
			},
			validate: func(t *testing.T, rec Record) {
				// Should use conversion year (2039 - 4 = 2035) in nome
				assert.Equal(t, "Tesouro Educa+ 2035", rec.Nome)
				assert.Equal(t, "2039-12-15", rec.DataVencimento)
				assert.Equal(t, "2035-01-15", rec.DataConversao)
			},
		},
		{
			name: "Regular bond - not Educa+ or Renda+",
			row: []string{
				"Tesouro IPCA+",
				"15/05/2035",
				"22/12/2025",
				"7,29", "7,41", "2374,37", "2348,76", "2348,76",
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
			tt.validate(t, rec)
		})
	}
}
