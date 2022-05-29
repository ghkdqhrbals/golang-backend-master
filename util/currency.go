package util

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	WON = "WON"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, WON:
		return true
	}
	return false
}
