package domain

import "math/big"

// Trade struct is the standard representation of any trade in this application
// price and quantity are in big.Float to utilize precision with our trades
type Trade struct {
	Currency string
	Price    *big.Float
	Quantity *big.Float
}

