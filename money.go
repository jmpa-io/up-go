package up

// Money represents the default attributes for any money-related data returned
// from the API. This is largely used to capture the amount of funds returned
// for transactions, or for listing the amount of funds available in an account.
type Money struct {
	CurrencyCode     string `json:"currencyCode"`
	Value            string `json:"value"`
	ValueInBaseUnits int64  `json:"valueInBaseUnits"`
}
