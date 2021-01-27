package main

import (
	"encoding/csv"
	"log"
	"os"
)

func readCsv(csvName string) ([][]string, error) {
	csvFile, err := os.Open(csvName)
	if err != nil {
		return nil, err
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	// Read the first line to make sure we skip the headers
	reader.Read()

	// Read rest of file
	csvLines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return csvLines, nil
}

func handleIncomingWebhooks(dbHandler *DBHandler) {
	log.Print("Got webhook, handling body now")
	voucherRequest, err := VoucherRequestFromEmailBody([]byte{})
	if err != nil {
		panic(err)
	}

	log.Print("Getting vouchers for ", voucherRequest)
	vouchers, err := dbHandler.GetAvailableVouchers(voucherRequest.QuantityVouchers)
	if err != nil {
		panic(err)
	}

	log.Printf("Found vouchers: %v. Marking them as sent to %s.", vouchers, voucherRequest.Email)
	err = dbHandler.MarkVouchersSent(vouchers, voucherRequest.Email)
	if err != nil {
		panic(err)
	}

	vouchers2, err := dbHandler.GetAvailableVouchers(1)
	if err != nil {
		panic(err)
	}
	log.Print("Got last voucher", vouchers2)
}

func main() {
	log.Print("Creating database file")
	err := CreateDatabaseFile()
	if err != nil {
		panic(err)
	}

	log.Print("Creating db handler")
	dbHandler, err := NewDBHandler()
	if err != nil {
		panic(err)
	}

	log.Print("Initializing db handler")
	err = dbHandler.Init()
	if err != nil {
		panic(err)
	}

	csvName := "codes.csv"
	log.Print("Reading csv with name: ", csvName)
	csvLines, err := readCsv(csvName)
	if err != nil {
		panic(err)
	}

	for _, line := range csvLines {
		voucher := VoucherFromCsvLine(line)
		log.Print("Creating voucher from csv: ", voucher)
		dbHandler.CreateVoucher(voucher)
	}

	handleIncomingWebhooks(&dbHandler)

	// hvem email skal sendes til, antall promos, og hvilken promotion kode/leverandør

	// 1: hent antall koder de skal ha.
}
