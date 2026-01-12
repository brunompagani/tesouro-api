package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	url := flag.String("url", defaultURL, "URL to download CSV from")
	outDir := flag.String("outdir", defaultOutDir, "Output directory for generated files")
	flag.Parse()

	if err := run(*url, *outDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(url, outDir string) error {
	// Create output directory
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Download CSV
	resp, err := downloadCSV(url)
	if err != nil {
		return fmt.Errorf("failed to download CSV: %w", err)
	}
	defer resp.Body.Close()

	// Parse CSV and extract latest records
	latest, err := parseCSV(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	// Convert to sorted slice and set DataInicio (start date)
	records := make([]Record, 0, len(latest))
	for _, asset := range latest {
		asset.record.DataInicio = asset.dataBaseMin.Format("2006-01-02")
		records = append(records, asset.record)
	}
	sortRecords(records)

	// Write JSON output
	jsonPath := filepath.Join(outDir, "latest.json")
	if err := writeJSON(records, jsonPath); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	// Write CSV output
	csvPath := filepath.Join(outDir, "latest.csv")
	if err := writeCSV(records, csvPath); err != nil {
		return fmt.Errorf("failed to write CSV: %w", err)
	}

	fmt.Printf("Successfully processed %d records\n", len(records))
	return nil
}
