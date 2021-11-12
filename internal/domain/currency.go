package domain

// Currencies publishes a list of all the supported currencies coming from newCurrencyRegistry
var Currencies = newCurrencyRegistry()

// Currency is struct that contains the currency name we choose to work with
type Currency struct {
	Name string
}

// currencyRegistry consist of all the currencies we would want to support in the application
// in shorted way possible this works like enums
type currencyRegistry struct {
	BtcUsd *Currency
	EthUsd *Currency
	EthBtc *Currency
}


// newCurrencyRegistry create already defined currencies, this prevents unnecessary typo errors
func newCurrencyRegistry() *currencyRegistry {
	btcUsd := &Currency{Name: "BTC-USD"}
	ethUsd := &Currency{Name: "ETH-USD"}
	ethBtc := &Currency{Name: "ETH-BTC"}

	return &currencyRegistry{btcUsd, ethUsd, ethBtc}
}
