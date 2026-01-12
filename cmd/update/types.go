package main

import "time"

const (
	defaultURL   = "https://www.tesourotransparente.gov.br/ckan/dataset/df56aa42-484a-4a59-8184-7676580c81e3/resource/796d2059-14e9-44e3-80c9-2d9e30b405c1/download/precotaxatesourodireto.csv"
	defaultOutDir = "public"
)

type Record struct {
	Nome             string  `json:"nome"`              // Combined: tipo_titulo + year (conversion year for Renda+ Aposentadoria Extra, maturity year otherwise)
	DataInicio       string  `json:"data_inicio"`       // ISO format: yyyy-mm-dd (oldest Data Base for this bond)
	DataConversao    string  `json:"data_conversao"`    // ISO format: yyyy-mm-dd (conversion date for Renda+ Aposentadoria Extra, empty otherwise)
	DataVencimento   string  `json:"data_vencimento"`   // ISO format: yyyy-mm-dd
	DataBase         string  `json:"data_base"`         // ISO format: yyyy-mm-dd (latest Data Base)
	TaxaCompraManha  float64 `json:"taxa_compra_manha"`
	TaxaVendaManha   float64 `json:"taxa_venda_manha"`
	PUCompraManha    float64 `json:"pu_compra_manha"`
	PUVendaManha     float64 `json:"pu_venda_manha"`
	PUBaseManha      float64 `json:"pu_base_manha"`
	tipoTitulo       string  // Internal: used for grouping only
}

type assetRecord struct {
	dataBaseMax time.Time // Latest Data Base (for keeping the most recent record)
	dataBaseMin time.Time // Oldest Data Base (start date)
	record      Record
}
