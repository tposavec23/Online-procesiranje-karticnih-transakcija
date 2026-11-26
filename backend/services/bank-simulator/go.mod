module sevenpay/backend/services/bank-simulator

go 1.27.3

replace sevenpay/backend/libs/h2h => ../../libs/h2h

require (
	github.com/moov-io/iso8583 v0.26.4
	github.com/moov-io/iso8583-connection v0.9.1
	github.com/stretchr/testify v1.12.1
	sevenpay/backend/libs/h2h v0.0.0-00010101000000-000000000000
)

require (
	github.com/yerden/go-util v1.1.4 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/text v0.43.0 // indirect
)
