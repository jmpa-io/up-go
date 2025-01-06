package up

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
)

// TagsPaginationWrapper is a pagination wrapper for a slice of TagResource. It
// is used to organize paginated tag data received from the API.
type TagsPaginationWrapper WrapperSlice[TagResource]

// wrapTags wraps the given slice of tags into a TagsPaginationWrapper. Each
// tag string is transformed into a TagResource before being sent to the API.
func wrapTags(tags []string) (wrappedTags TagsPaginationWrapper) {
	for _, t := range tags {
		wrappedTags.Data = append(
			wrappedTags.Data,
			TagResource{Object: Object{Type: "tags", ID: t}},
		)
	}
	return wrappedTags
}

// ListTagsOption defines the available options used to configure the ListTags
// function when listing tags from the API.
type ListTagsOption struct {
	listOption
}

// ListTagsOptionPageSize sets the page size used when when listing tags from
// the API. This option affects how many tags are returned at once - increasing
// this can improve performance as the number of API calls is reduced.
func ListTagsOptionPageSize(size int) ListTagsOption {
	return ListTagsOption{newListOption("page[size]", strconv.Itoa(size))}
}

// ListTags returns a list of ALL tags associated with the user from the API.
// This function supports pagination, and is configurable by the given
// ListTagsOption options.
// https://developer.up.com.au/#tags.
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
		var resp TagsPaginationWrapper
		if _, err := c.sender(newCtx, sr, &resp); err != nil {
			return nil, err
		}

		// extract response data.
		tags = append(tags, resp.Data...)

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

// AddTagsToTransaction adds the given tags to a transaction, via the given
// transaction id. This function only returns an error, and won't return
// anything if successful.
// https://developer.up.com.au/#post_transactions_transactionId_relationships_tags.
func (c *Client) AddTagsToTransaction(ctx context.Context, id string, tags []string) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "AddTagsToTransaction")
	defer span.End()

	// add tags.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transactions/%s/relationships/tags", id),
		body:   wrapTags(tags),
	}, nil)
	return err
}

// RemoveTagsFromTransaction removes the given tags from a transaction,
// identified by the given transaction id. This function only returns an error,
// and won't return anything if successful.
// https://developer.up.com.au/#delhttps://developer.up.com.au/#tagsete_transactions_transactionId_relationships_tags.
func (c *Client) RemoveTagsFromTransaction(ctx context.Context, id string, tags []string) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "RemoveTagsToTransaction")
	defer span.End()

	// remove tags.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/transactions/%s/relationships/tags", id),
		body:   wrapTags(tags),
	}, nil)
	return err
}
