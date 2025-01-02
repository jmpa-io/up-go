package up

import (
	"time"
)

// AccountType represents the type of account.
type AccountType string

const (
	AccountTypeSaver         AccountType = "SAVER"         // A savings account.
	AccountTypeTransactional AccountType = "TRANSACTIONAL" // A transactional account.
	AccountTypeHomeLoan      AccountType = "HOME_LOAN"     // A home_loan account.
)

// AccountOwnershipType represents the type of ownership for an account.
type AccountOwnershipType string

const (
	AccountOwnershipTypeIndividual AccountOwnershipType = "INDIVIDUAL" // An account owned by a single person.
	AccountOwnershipTypeJoint      AccountOwnershipType = "JOINT"      // An account owned by multiple people.
)

// AccountResource defines the core details of an account.
type AccountResource struct {
	DisplayName   string               `json:"displayName"`
	AccountType   AccountType          `json:"accountType"`
	OwnershipType AccountOwnershipType `json:"ownershipType"`
	Balance       Money                `json:"balance"`
	CreatedAt     time.Time            `json:"createdAt"`
}

// AccountRelationships defines the relationships to other resources for
// an account.
type AccountRelationships struct {
	Transactions WrapperOmittable `json:"transactions"`
}

// AccountDataWrapper wraps the resources and relationships for account data
// returned from the API.
type AccountDataWrapper Data[AccountResource, AccountRelationships]
