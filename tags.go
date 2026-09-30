package up

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// TagsPaginationWrapper is a pagination wrapper for a slice of TagResource.
type TagsPaginationWrapper WrapperSlice[TagResource]

// wrapTags converts a slice of tag label strings into the API's required
// TagsPaginationWrapper format.
func wrapTags(tags []string) TagsPaginationWrapper {
	var w TagsPaginationWrapper
	for _, t := range tags {
		w.Data = append(w.Data, TagResource{Object: Object{Type: "tags", ID: t}})
	}
	return w
}

// ListTagsOption configures a ListTags call.
type ListTagsOption struct {
	listOption
}

// ListTagsOptionPageSize sets the number of tags returned per page.
func ListTagsOptionPageSize(size int) ListTagsOption {
	return ListTagsOption{newListOption("page[size]", strconv.Itoa(size))}
}

// ListTags returns all tags currently in use for the authenticated user.
// https://developer.up.com.au/#get_tags.
func (c *Client) ListTags(
	ctx context.Context,
	opts ...ListTagsOption,
) (tags []TagResource, err error) {

	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListTags")
	defer span.End()

	sr := senderRequest{
		method:  http.MethodGet,
		path:    "/tags",
		queries: setupQueries(opts),
	}

	for {
		var resp TagsPaginationWrapper
		if _, err := c.sender(newCtx, sr, &resp); err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("failed to list tags: %v", err))
			span.RecordError(err)
			return nil, fmt.Errorf("listing tags: %w", err)
		}
		tags = append(tags, resp.Data...)
		if resp.Links.Next == "" {
			break
		}
		sr.path = strings.Replace(resp.Links.Next, c.endpoint, "", 1)
		sr.queries = nil
	}
	return tags, nil
}

// AddTagsToTransaction adds the given tags to a transaction.
// Up supports a maximum of 6 tags per transaction. Duplicate tags are silently
// ignored by the API.
// https://developer.up.com.au/#post_transactions_transactionId_relationships_tags.
func (c *Client) AddTagsToTransaction(ctx context.Context, id string, tags []string) error {

	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "AddTagsToTransaction")
	defer span.End()

	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transactions/%s/relationships/tags", id),
		body:   wrapTags(tags),
	}, nil)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to add tags to transaction %s: %v", id, err))
		span.RecordError(err)
		return fmt.Errorf("adding tags to transaction %s: %w", id, err)
	}
	return nil
}

// RemoveTagsFromTransaction removes the given tags from a transaction.
// Tags not present on the transaction are silently ignored by the API.
// https://developer.up.com.au/#delete_transactions_transactionId_relationships_tags.
func (c *Client) RemoveTagsFromTransaction(ctx context.Context, id string, tags []string) error {

	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "RemoveTagsFromTransaction")
	defer span.End()

	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/transactions/%s/relationships/tags", id),
		body:   wrapTags(tags),
	}, nil)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to remove tags from transaction %s: %v", id, err))
		span.RecordError(err)
		return fmt.Errorf("removing tags from transaction %s: %w", id, err)
	}
	return nil
}
