package up

import (
	"time"
)

// TransactionStatus represents the status of a transaction.
type TransactionStatus string

const (
	TransactionStatusHeld    TransactionStatus = "HELD"    // Transaction amount has not yet left your account.
	TransactionStatusSettled TransactionStatus = "SETTLED" // Transaction amount has left your account.
)

// TransactionCardPurchaseMethod defines the method used to complete a purchase.
type TransactionCardPurchaseMethod string

const (
	TransactionCardPurchaseMethodBarCode       TransactionCardPurchaseMethod = "BAR_CODE"       // Purchased via barcode.
	TransactionCardPurchaseMethodOCR           TransactionCardPurchaseMethod = "OCR"            // Purchased via Optical Character Recognition (OCR).
	TransactionCardPurchaseMethodCardPin       TransactionCardPurchaseMethod = "CARD_PIN"       // Purchased via card PIN.
	TransactionCardPurchaseMethodCardDetails   TransactionCardPurchaseMethod = "CARD_DETAILS"   // Purchased via card details.
	TransactionCardPurchaseMethodCardOnFile    TransactionCardPurchaseMethod = "CARD_ON_FILE"   // Purchased via stored card on file.
	TransactionCardPurchaseMethordEcommerce    TransactionCardPurchaseMethod = "ECOMMERCE"      // Purchased via online purchase (e-commerce).
	TransactionCardPurchaseMethodMagneticStrip TransactionCardPurchaseMethod = "MAGNETIC_STRIP" // Purchased via magnetic stripe.
	TransactionCardPurchaseMethodContactless   TransactionCardPurchaseMethod = "CONTACTLESS"    // Purchased via contactless payment.
)

// TransactionResourceHoldInfo defines details about a held transaction.
type TransactionResourceHoldInfo struct {
	Amount        Money `json:"amount"`
	ForeignAmount Money `json:"foreignAmount"`
}

// TransactionResourceRoundUp defines details about the round-up and
// boost-portion amounts associated with a transaction.
type TransactionResourceRoundUp struct {
	Amount       Money `json:"amount"`
	BoostPortion Money `json:"boostPortion"`
}

// TransactionResourceCashback defines details about any cashbacks earned with
// a transaction.
type TransactionResourceCashback struct {
	Description string `json:"description"`
	Amount      Money  `json:"amount"`
}

// TransactionResourceCardPurchaseMethod defines details about the card used,
// and purchase method, for a transaction.
type TransactionResourceCardPurchaseMethod struct {
	CardNumberSuffix string                        `json:"cardNumberSuffix"`
	Method           TransactionCardPurchaseMethod `json:"method"`
}

// TransactionResourceNote represents a note attached to a transaction.
type TransactionResourceNote struct {
	Text string `json:"text"`
}

// TransactionResourcePerformingCustomer defines details about the customer
// who performed the transaction.
type TransactionResourcePerformingCustomer struct {
	DisplayName string `json:"displayName"`
}

// TransactionResource defines the core details of a transaction.
type TransactionResource struct {
	Status             TransactionStatus                     `json:"status"`
	RawText            string                                `json:"rawText"`
	Description        string                                `json:"description"`
	Message            string                                `json:"message"`
	IsCategorizable    bool                                  `json:"isCategorizable"`
	HoldInfo           TransactionResourceHoldInfo           `json:"holdInfo"`
	RoundUp            TransactionResourceRoundUp            `json:"roundUp"`
	Cashback           TransactionResourceCashback           `json:"cashback"`
	Amount             Money                                 `json:"amount"`
	ForeignAmount      Money                                 `json:"foreignAmount"`
	CardPurchaseMethod TransactionResourceCardPurchaseMethod `json:"cardPurchaseMethod"`
	CreatedAt          time.Time                             `json:"createdAt"`
	SettledAt          time.Time                             `json:"settledAt"`
	TransactionType    string                                `json:"transactionType"`
	Note               TransactionResourceNote               `json:"note"`
	PerformingCustomer TransactionResourcePerformingCustomer `json:"performingCustomer"`
	DeepLinkURL        string                                `json:"deepLinkURL"`
}

// TransactionRelationships defines the relationships to other resources for
// a transaction.
type TransactionRelationships struct {
	Account         Wrapper[Object]      `json:"account"`
	TransferAccount Wrapper[Object]      `json:"transferAccount"`
	Category        Wrapper[Object]      `json:"category"`
	ParentCategory  Wrapper[Object]      `json:"parentCategory"`
	Tags            WrapperSlice[Object] `json:"tags"`
	Attachment      Wrapper[Object]      `json:"attachment"`
}

// TransactionDataWrapper wraps the resources and relationships for transaction
// data returned from the API.
type TransactionDataWrapper Data[TransactionResource, TransactionRelationships]
