package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "valid date",
			input:    "15/05/2035",
			expected: "2035-05-15",
			wantErr:  false,
		},
		{
			name:     "valid date with spaces",
			input:    "  01/01/2024  ",
			expected: "2024-01-01",
			wantErr:  false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "2025-01-01",
			wantErr: true,
		},
		{
			name:    "invalid date",
			input:   "32/01/2025",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDate(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseFloatBR(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{
			name:     "simple comma decimal",
			input:    "7,29",
			expected: 7.29,
			wantErr:  false,
		},
		{
			name:     "with thousand separator",
			input:    "1.234,56",
			expected: 1234.56,
			wantErr:  false,
		},
		{
			name:     "integer",
			input:    "1000",
			expected: 1000.0,
			wantErr:  false,
		},
		{
			name:     "with spaces",
			input:    "  7,29  ",
			expected: 7.29,
			wantErr:  false,
		},
		{
			name:     "negative value",
			input:    "-0,01",
			expected: -0.01,
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
			wantErr:  false,
		},
		{
			name:    "invalid format",
			input:   "abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseFloatBR(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 0.0001)
			}
		})
	}
}

func TestFormatFloatBR(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "simple decimal",
			input:    7.29,
			expected: "7,29",
		},
		{
			name:     "integer",
			input:    1000.0,
			expected: "1000",
		},
		{
			name:     "negative",
			input:    -0.01,
			expected: "-0,01",
		},
		{
			name:     "zero",
			input:    0.0,
			expected: "0",
		},
		{
			name:     "large number",
			input:    1234.56,
			expected: "1234,56",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFloatBR(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseRecord(t *testing.T) {
	tests := []struct {
		name     string
		row      []string
		wantErr  bool
		validate func(t *testing.T, rec Record)
	}{
		{
			name: "valid record",
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
			wantErr: false,
			validate: func(t *testing.T, rec Record) {
				assert.Equal(t, "Tesouro IPCA+ 2035", rec.Nome)
				assert.Equal(t, "2035-05-15", rec.DataVencimento)
				assert.Equal(t, "2025-12-22", rec.DataBase)
				assert.InDelta(t, 7.29, rec.TaxaCompraManha, 0.01)
				assert.InDelta(t, 7.41, rec.TaxaVendaManha, 0.01)
				assert.InDelta(t, 2374.37, rec.PUCompraManha, 0.01)
			},
		},
		{
			name: "invalid date",
			row: []string{
				"Tesouro IPCA+",
				"invalid",
				"22/12/2025",
				"7,29",
				"7,41",
				"2374,37",
				"2348,76",
				"2348,76",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec, err := parseRecord(tt.row)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, rec)
				}
			}
		})
	}
}
