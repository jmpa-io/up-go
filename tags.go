package up

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
)

// Tags represents tags in Up.
type TagsWrapper WrapperSlice[TagResource]

// wrapTags wraps the given tags in the data wrapper, ready to be sent to the API.
func wrapTags(tags []string) (wrappedTags TagsWrapper) {
	for _, t := range tags {
		wrappedTags.Data = append(
			wrappedTags.Data,
			TagResource{Object: Object{Type: "tags", ID: t}},
		)
	}
	return wrappedTags
}

type ListTagsOption struct {
	ListOption
}

func ListTagsOptionPageSize(size int) ListTagsOption {
	return ListTagsOption{NewListOption("page[size]", strconv.Itoa(size))}
}

// https://developer.up.com.au/#tags
func (c *Client) ListTags(
	ctx context.Context,
	opts ...ListTagsOption,
) (tags []TagResource, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListTags")
	defer span.End()

	// setup request.
	sr := senderRequest{
		method:  http.MethodGet,
		path:    "/tags",
		queries: setupQueries(opts),
	}

	// retrieve tags.
	for {

		// get response.
		var resp TagsWrapper
		if _, err := c.sender(newCtx, sr, &resp); err != nil {
			return nil, err
		}

		// extract response data.
		for _, t := range resp.Data {
			tags = append(tags, t)
		}

		// paginate?
		if resp.Links.Next == "" {
			break
		}
		c.logger.Debug(resp.Links.Next)
		sr.path = strings.Replace(resp.Links.Next, c.endpoint, "", 1)
		sr.queries = nil
	}
	return tags, nil

}

// AddTagsToTransaction, using the given transaction id, adds the given tags to
// the transaction.
// https://developer.up.com.au/#post_transactions_transactionId_relationships_tags
func (c *Client) AddTagsToTransaction(ctx context.Context, id string, tags []string) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "AddTagsToTransaction")
	defer span.End()

	// add tags.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transactions/%s/relationships/tags", id),
		data:   wrapTags(tags),
	}, nil)
	return err
}

// RemoveTagsFromTransaction, using the given transaction id, removes the given
// tags from a transaction.
// https://developer.up.com.au/#delhttps://developer.up.com.au/#tagsete_transactions_transactionId_relationships_tags
func (c *Client) RemoveTagsFromTransaction(ctx context.Context, id string, tags []string) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "RemoveTagsToTransaction")
	defer span.End()

	// remove tags.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/transactions/%s/relationships/tags", id),
		data:   wrapTags(tags),
	}, nil)
	return err
}
