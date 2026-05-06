package up

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// CategoryAttributes defines the attributes of an Up category.
type CategoryAttributes struct {
	Name string `json:"name"`
}

// CategoryRelationships defines the parent/child relationships for a category.
type CategoryRelationships struct {
	Parent struct {
		Data *Object `json:"data"` // nil for top-level categories
	} `json:"parent"`
	Children struct {
		Data []Object `json:"data"`
	} `json:"children"`
}

// CategoryData represents a single category returned from the Up API.
type CategoryData Data[CategoryAttributes, CategoryRelationships]

// CategoryPaginationWrapper is a pagination wrapper for a slice of CategoryData.
type CategoryPaginationWrapper WrapperSlice[CategoryData]

// setCategoryBody is the request body for PATCH /transactions/{id}/relationships/category.
type setCategoryBody struct {
	Data *categoryRef `json:"data"` // nil to de-categorise
}

type categoryRef struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// ListCategories returns all categories from the Up API.
// https://developer.up.com.au/#get_categories.
func (c *Client) ListCategories(ctx context.Context) ([]CategoryData, error) {

	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListCategories")
	defer span.End()

	var resp CategoryPaginationWrapper
	if _, err := c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   "/categories",
	}, &resp); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to list categories: %v", err))
		span.RecordError(err)
		return nil, fmt.Errorf("listing categories: %w", err)
	}
	return resp.Data, nil
}

// SetTransactionCategory assigns an Up category to a transaction. Pass an
// empty categoryID to de-categorise the transaction.
// https://developer.up.com.au/#patch_transactions_transactionId_relationships_category.
func (c *Client) SetTransactionCategory(
	ctx context.Context,
	transactionID string,
	categoryID string,
) error {

	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "SetTransactionCategory")
	defer span.End()

	var body setCategoryBody
	if categoryID != "" {
		body.Data = &categoryRef{Type: "categories", ID: categoryID}
	}

	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodPatch,
		path:   fmt.Sprintf("/transactions/%s/relationships/category", transactionID),
		body:   body,
	}, nil)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to set category on transaction %s: %v", transactionID, err))
		span.RecordError(err)
		return fmt.Errorf("setting category on transaction %s: %w", transactionID, err)
	}
	return nil
}

// GetCategory retrieves a single category by its Up category ID.
// https://developer.up.com.au/#get_categories_id.
func (c *Client) GetCategory(ctx context.Context, id string) (*CategoryData, error) {

	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "GetCategory")
	defer span.End()

	var resp struct {
		Data CategoryData `json:"data"`
	}
	if _, err := c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/categories/%s", id),
	}, &resp); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get category %s: %v", id, err))
		span.RecordError(err)
		return nil, fmt.Errorf("getting category %s: %w", id, err)
	}
	return &resp.Data, nil
}
