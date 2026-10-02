package up

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func toJSON(v any) string {
	// A nil slice marshals as null; callers (and models) expect an empty list.
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Slice && rv.IsNil() {
		return "[]"
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		safe, _ := json.Marshal(err.Error())
		return `{"error":` + string(safe) + `}`
	}
	return string(b)
}

func (c *Client) PingString(ctx context.Context) (string, error) {
	p, err := c.Ping(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Up API is reachable. Status: %s (id: %s)", p.Meta.StatusEmoji, p.Meta.ID), nil
}

func (c *Client) ListAccountsString(ctx context.Context, opts ...ListAccountsOption) (string, error) {
	accounts, err := c.ListAccounts(ctx, opts...)
	if err != nil {
		return "", err
	}
	return toJSON(accounts), nil
}

func (c *Client) ListTransactionsString(ctx context.Context, opts ...ListTransactionsOption) (string, error) {
	txns, err := c.ListTransactions(ctx, opts...)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\n\n(total: %d)", toJSON(txns), len(txns)), nil
}

func (c *Client) ListTransactionsByAccountString(ctx context.Context, accountID string, opts ...ListTransactionsOption) (string, error) {
	txns, err := c.ListTransactionsByAccount(ctx, accountID, opts...)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\n\n(total: %d)", toJSON(txns), len(txns)), nil
}

func (c *Client) GetTransactionString(ctx context.Context, id string) (string, error) {
	txn, err := c.GetTransaction(ctx, id)
	if err != nil {
		return "", err
	}
	return toJSON(txn), nil
}

func (c *Client) ListCategoriesString(ctx context.Context) (string, error) {
	cats, err := c.ListCategories(ctx)
	if err != nil {
		return "", err
	}
	return toJSON(cats), nil
}

func (c *Client) SetTransactionCategoryString(ctx context.Context, transactionID, categoryID string) (string, error) {
	if err := c.SetTransactionCategory(ctx, transactionID, categoryID); err != nil {
		return "", err
	}
	if categoryID == "" {
		return fmt.Sprintf("Transaction %s de-categorised.", transactionID), nil
	}
	return fmt.Sprintf("Transaction %s category set to %q.", transactionID, categoryID), nil
}

func (c *Client) ListTagsString(ctx context.Context, opts ...ListTagsOption) (string, error) {
	tags, err := c.ListTags(ctx, opts...)
	if err != nil {
		return "", err
	}
	if len(tags) == 0 {
		return "No tags found.", nil
	}
	var sb strings.Builder
	for _, t := range tags {
		fmt.Fprintf(&sb, "%s\n", t.ID)
	}
	return fmt.Sprintf("Tags (%d):\n%s", len(tags), sb.String()), nil
}

func (c *Client) AddTagsToTransactionString(ctx context.Context, transactionID string, tags []string) (string, error) {
	if err := c.AddTagsToTransaction(ctx, transactionID, tags); err != nil {
		return "", err
	}
	return fmt.Sprintf("Added tags %v to transaction %s.", tags, transactionID), nil
}

func (c *Client) RemoveTagsFromTransactionString(ctx context.Context, transactionID string, tags []string) (string, error) {
	if err := c.RemoveTagsFromTransaction(ctx, transactionID, tags); err != nil {
		return "", err
	}
	return fmt.Sprintf("Removed tags %v from transaction %s.", tags, transactionID), nil
}
