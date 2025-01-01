package up

// Object represents the default attributes for data types returned from the
// API. This struct is largely used as composition for other data within this
// package (eg. Accounts, Tags, Transactions, etc).
type Object struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}
