package up

import (
	"time"
)

type AccountType string

const (
	AccountTypeSaver         AccountType = "SAVER"
	AccountTypeTransactional             = "TRANSACTIONAL"
	AccountTypeHomeLoan                  = "HOME_LOAN"
)

type AccountOwnershipType string

const (
	AccountOwnershipTypeIndividual AccountOwnershipType = "INDIVIDUAL"
	AccountOwnershipTypeJoint                           = "JOINT"
)

type AccountResource struct {
	DisplayName   string               `json:"displayName"`
	AccountType   AccountType          `json:"accountType"`
	OwnershipType AccountOwnershipType `json:"ownershipType"`
	Balance       Money                `json:"balance"`
	CreatedAt     time.Time            `json:"createdAt"`
}

type AccountRelationships struct {
	Transactions WrapperOmittable `json:"transactions"`
}

// Account represents an account in Up.
type Account Data[AccountResource, AccountRelationships]
