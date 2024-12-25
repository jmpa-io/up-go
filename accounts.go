package up

import "time"

type AccountAttributes struct {
	DisplayName   string    `json:"displayName"`
	AccountType   string    `json:"accountType"`
	OwnershipType string    `json:"ownershipType"`
	Balance       Amount    `json:"balance"`
	CreatedAt     time.Time `json:"createdAt"`
}

type AccountRelationships struct {
	Transactions struct {
		Links RelatedLink `json:"links"`
	} `json:"transactions"`
}

// AccountDataWrapper represents an account in Up.
type AccountDataWrapper Data[AccountAttributes, AccountRelationships]

// AccountPaginationWrapper a pagination wrapper for a slice of AccountDataWrapper.
type AccountPaginationWrapper PaginationWrapper[AccountDataWrapper]
