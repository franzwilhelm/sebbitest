package main

type Voucher struct {
	ID    int
	Code  string
	Sent  int
	Email string
}

func VoucherFromCsvLine(line []string) Voucher {
	var code = line[0]
	var sent int

	if line[1] == "Send" {
		sent = 1
	}

	return Voucher{
		Code: code,
		Sent: sent,
	}
}

type VoucherRequest struct {
	Email            string
	QuantityVouchers int
}

// VoucherRequestFromEmailBody ...
func VoucherRequestFromEmailBody(body []byte) (VoucherRequest, error) {
	// TODO: parse body and generate request from it
	return VoucherRequest{
		Email:            "digg@sleep.com",
		QuantityVouchers: 2,
	}, nil
}
