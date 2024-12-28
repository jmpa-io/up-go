package up

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
)

type TransactionStatus string

const (
	TransactionStatusHeld    TransactionStatus = "HELD"
	TransactionStatusSettled                   = "SETTLED"
)

type CardPurchaseMethod string

const (
	CardPurchaseMethodBarCode       CardPurchaseMethod = "BAR_CODE"
	CardPurchaseMethodOCR                              = "OCR"
	CardPurchaseMethodCardPin                          = "CARD_PIN"
	CardPurchaseMethodCardDetails                      = "CARD_DETAILS"
	CardPurchaseMethodCardOnFile                       = "CARD_ON_FILE"
	CardPurchaseMethordEcommerce                       = "ECOMMERCE"
	CardPurchaseMethodMagneticStrip                    = "MAGNETIC_STRIP"
	CardPurchaseMethodContactless                      = "CONTACTLESS"
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
	CardNumberSuffix string             `json:"cardNumberSuffix"`
	Method           CardPurchaseMethod `json:"method"`
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
	SettledAt          time.Time                             `json:"settledAt"`
	CreatedAt          time.Time                             `json:"createdAt"`
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

// TransactionsWrapper is a pagination wrapper for a slice of TransactionData.
type TransactionWrapper WrapperSlice[Transaction]

// ListTransactionsOption defines the options for listing transactions in the
// authed account.
type ListTransactionsOption struct {
	name  string
	value string
}

func ListTransactionsOptionPageSize(size int) ListTransactionsOption {
	return ListTransactionsOption{
		name:  "page[size]",
		value: strconv.Itoa(size),
	}
}

type FilterStatus string

const (
	FilterStatusHeld    FilterStatus = "HELD"
	FilterStatusSettled FilterStatus = "SETTLED"
)

func ListTransactionsOptionStatus(status FilterStatus) ListTransactionsOption {
	return ListTransactionsOption{
		name:  "filter[status]",
		value: string(status),
	}
}

func ListTransactionsOptionSince(since time.Time) ListTransactionsOption {
	return ListTransactionsOption{
		name:  "filter[since]",
		value: since.Format(time.RFC3339),
	}
}

func ListTransactionsOptionUntil(until time.Time) ListTransactionsOption {
	return ListTransactionsOption{
		name:  "filter[until]",
		value: until.Format(time.RFC3339),
	}
}

func ListTransactionsOptionCategory(category string) ListTransactionsOption {
	return ListTransactionsOption{
		name:  "filter[category]",
		value: category,
	}
}

func ListTransactionsOptionTag(tag string) ListTransactionsOption {
	return ListTransactionsOption{
		name:  "filter[tag]",
		value: tag,
	}
}

// ListTransactions list all transactions for the authed account.
// https://developer.up.com.au/#get_transactions
func (c *Client) ListTransactions(ctx context.Context,
	options ...ListTransactionsOption,
) (transactions []TransactionResource, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListTransactions")
	defer span.End()

	// setup queries.
	queries := make(url.Values)
	for _, o := range options {
		queries[o.name] = []string{o.value}
	}

	// default queries.
	if _, ok := queries["page[size]"]; !ok {
		queries["page[size]"] = []string{"100"}
	}

	// setup request.
	sr := senderRequest{
		method:  http.MethodGet,
		path:    "/transactions",
		queries: queries,
	}

	// retrieve transactions.
	for {

		// get response.
		var resp TransactionWrapper
		if _, err := c.sender(newCtx, sr, &resp); err != nil {
			return nil, err
		}

		// extract response data.
		for _, t := range resp.Data {
			transactions = append(transactions, t.Attributes)
		}

		// paginate?
		if resp.Links.Next == "" {
			break
		}
		c.logger.Debug(resp.Links.Next)
		sr.path = strings.Replace(resp.Links.Next, c.endpoint, "", 1)
		sr.queries = nil
	}
	return transactions, nil
}
