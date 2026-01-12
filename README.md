# Tesouro Direto Latest Snapshot API

A small service that publishes a daily snapshot of the latest Tesouro Direto bond prices from the official Brazilian Treasury CSV dataset. The original CSV contains historical data for all bonds, which can be large (100k+ rows). This project generates a compact snapshot containing only the most recent data point per asset, making it easy to consume in Google Sheets or other tools.

The service runs automatically via GitHub Actions once per day and publishes the results as static JSON and CSV files via GitHub Pages.

## Output Files

After running, the service generates two files in the `public/` directory:

- **latest.json** - Latest snapshot in JSON format
- **latest.csv** - Latest snapshot in CSV format (semicolon-delimited, PT-BR number format)

These files are published to GitHub Pages and accessible at:

- `https://<user>.github.io/<repo>/latest.json`
- `https://<user>.github.io/<repo>/latest.csv`

## Local Development

### Prerequisites

- Go 1.21 or later

### Running Locally

To run the updater locally:

```bash
go run ./cmd/update
```

This will:
1. Download the CSV from Tesouro Direto
2. Parse it and extract the latest record per asset (by `Data Base`)
3. Track the oldest `Data Base` date per asset (start date)
4. Generate `public/latest.json` and `public/latest.csv`

### Command-line Options

```bash
go run ./cmd/update --url <custom-url> --outdir <output-directory>
```

- `--url`: Override the CSV URL (default: official Tesouro Direto URL)
- `--outdir`: Output directory (default: `public/`)

## GitHub Pages Setup

1. Go to your repository settings on GitHub
2. Navigate to **Pages** in the left sidebar
3. Under **Source**, select:
   - Branch: `gh-pages` (This branch is created after the first workflow run - trigger the workflow first if it doesn't exist yet)
   - Folder: `/ (root)`
4. Click **Save**

After the first workflow run, your files will be available at the URLs shown above.

## Using in Google Sheets

### Method 1: Custom Function (Recommended)

Install the `TESOURODIRETO()` custom function for easy access to bond data:

1. Open your Google Sheet
2. Go to **Extensions** → **Apps Script**
3. Create a new file called `TesouroDireto.gs`
4. Copy and paste the code from `GOOGLE_SHEETS_FUNCTION.gs` in this repository
5. **Important**: Replace `YOUR_USERNAME` and `YOUR_REPO_NAME` in the `apiUrl` variable with your actual GitHub Pages URL
   - Example: `https://brunompagani.github.io/tesouro-api/latest.json`
6. Save the script (Ctrl+S or Cmd+S)
7. Close the Apps Script editor

#### Usage

```excel
=TESOURODIRETO("Tesouro IPCA+ 2035", "taxa_compra_manha")
=TESOURODIRETO("Tesouro IPCA+ 2035", "pu_venda_manha")
=TESOURODIRETO("Tesouro Prefixado 2008", "taxa_compra_manha", "2008-01-01")
```

**Parameters:**
- `nome` (required): Bond name, e.g., "Tesouro IPCA+ 2035"
- `campo` (required): Field to return - `data_vencimento`, `data_base`, `data_inicio`, `taxa_compra_manha`, `taxa_venda_manha`, `pu_compra_manha`, `pu_venda_manha`, `pu_base_manha`
- `data_vencimento` (optional): Maturity date in ISO format (yyyy-mm-dd) to differentiate bonds with the same name

### Method 2: Import CSV

In a Google Sheets cell, use:

```
=IMPORTDATA("https://<user>.github.io/<repo>/latest.csv")
```

This will import the CSV data into your spreadsheet. You can then use QUERY or FILTER functions to look up specific bonds.

## Data Format

### CSV Schema

- **Nome**: Combined bond name (e.g., "Tesouro IPCA+ 2035")
- **Data Inicio**: Start date - oldest Data Base date for this bond (ISO format: yyyy-mm-dd)
- **Data Vencimento**: Maturity date (ISO format: yyyy-mm-dd)
- **Data Base**: Latest reference date for the prices (ISO format: yyyy-mm-dd)
- **Taxa Compra Manha**: Morning buy rate (PT-BR format: comma as decimal separator)
- **Taxa Venda Manha**: Morning sell rate
- **PU Compra Manha**: Morning buy price (PU = Preço Unitário)
- **PU Venda Manha**: Morning sell price
- **PU Base Manha**: Morning base price

### JSON Schema

The JSON output uses snake_case field names:
- `nome`: Combined bond name (tipo_titulo + year)
- `data_inicio`: Start date - oldest Data Base date for this bond (ISO format: yyyy-mm-dd)
- `data_vencimento`: Maturity date (ISO format: yyyy-mm-dd)
- `data_base`: Latest reference date for the prices (ISO format: yyyy-mm-dd)
- `taxa_compra_manha`: Morning buy rate (float)
- `taxa_venda_manha`: Morning sell rate (float)
- `pu_compra_manha`: Morning buy price (float)
- `pu_venda_manha`: Morning sell price (float)
- `pu_base_manha`: Morning base price (float)

All numeric values are floats, and dates are ISO strings (yyyy-mm-dd).

## How It Works

1. **Daily Schedule**: The GitHub Action runs automatically at 07:00 UTC (04:00 BRT) every day
2. **Manual Trigger**: You can also trigger it manually via the GitHub Actions UI (workflow_dispatch)
3. **Processing**:
   - Downloads the full CSV from Tesouro Direto
   - Parses it using streaming (memory-efficient)
   - Groups records by asset key (Tipo Titulo + Data Vencimento)
   - Keeps only the latest record per asset (by Data Base date)
   - Tracks the oldest Data Base date per asset (start date)
   - Sorts the output for stable diffs
4. **Output**: Generates JSON and CSV files
5. **Deployment**: Publishes to the `gh-pages` branch via GitHub Actions

## License

MIT
