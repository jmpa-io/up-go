package up

import (
	"time"
)

type TransactionStatus string

const (
	TransactionStatusHeld    TransactionStatus = "HELD"
	TransactionStatusSettled                   = "SETTLED"
)

type TransactionCardPurchaseMethod string

const (
	TransactionCardPurchaseMethodBarCode       TransactionCardPurchaseMethod = "BAR_CODE"
	TransactionCardPurchaseMethodOCR                                         = "OCR"
	TransactionCardPurchaseMethodCardPin                                     = "CARD_PIN"
	TransactionCardPurchaseMethodCardDetails                                 = "CARD_DETAILS"
	TransactionCardPurchaseMethodCardOnFile                                  = "CARD_ON_FILE"
	TransactionCardPurchaseMethordEcommerce                                  = "ECOMMERCE"
	TransactionCardPurchaseMethodMagneticStrip                               = "MAGNETIC_STRIP"
	TransactionCardPurchaseMethodContactless                                 = "CONTACTLESS"
)

type TransactionResourceHoldInfo struct {
	Amount        Money `json:"amount"`
	ForeignAmount Money `json:"foreignAmount"`
}

type TransactionResourceRoundUp struct {
	Amount       Money `json:"amount"`
	BoostPortion Money `json:"boostPortion"`
}

type TransactionResourceCashback struct {
	Description string `json:"description"`
	Amount      Money  `json:"amount"`
}

type TransactionResourceCardPurchaseMethod struct {
	CardNumberSuffix string                        `json:"cardNumberSuffix"`
	Method           TransactionCardPurchaseMethod `json:"method"`
}

type TransactionResourceNote struct {
	Text string `json:"text"`
}

type TransactionResourcePerformingCustomer struct {
	DisplayName string `json:"displayName"`
}

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

type TransactionRelationships struct {
	Account         Wrapper[Object]      `json:"account"`
	TransferAccount Wrapper[Object]      `json:"transferAccount"`
	Category        Wrapper[Object]      `json:"category"`
	ParentCategory  Wrapper[Object]      `json:"parentCategory"`
	Tags            WrapperSlice[Object] `json:"tags"`
	Attachment      Wrapper[Object]      `json:"attachment"`
}

// Transaction represents a transaction in Up.
type Transaction Data[TransactionResource, TransactionRelationships]
