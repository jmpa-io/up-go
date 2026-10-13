package up

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// The constants and validators below describe the limits Up enforces today.
// They are advisory, for callers that want to check tags before a bulk run;
// the SDK itself never blocks a request on them, because the API is the
// authority and its limits may change.
const (
	// MaxTagLength is the maximum length of a tag, in characters (not bytes).
	// The limit applies to the raw string: leading and trailing spaces count.
	MaxTagLength = 30

	// MaxTagsPerTransaction is the maximum number of tags a single transaction
	// can have. A request that would take a transaction over this is rejected
	// by the API as a whole.
	MaxTagsPerTransaction = 6
)

// ValidateTag reports whether tag looks acceptable to the Up API. These rules
// were checked against the live API (see the integration test): empty and whitespace-only tags are rejected,
// as are tags longer than MaxTagLength characters.
func ValidateTag(tag string) error {
	if strings.TrimSpace(tag) == "" {
		return fmt.Errorf("invalid tag %q: must not be empty or only whitespace", tag)
	}
	if n := utf8.RuneCountInString(tag); n > MaxTagLength {
		return fmt.Errorf("invalid tag %q: %d characters, maximum is %d", tag, n, MaxTagLength)
	}
	return nil
}

// ValidateTags validates each tag with ValidateTag and checks that the request
// itself does not name more than MaxTagsPerTransaction tags.
func ValidateTags(tags []string) error {
	if len(tags) > MaxTagsPerTransaction {
		return fmt.Errorf("too many tags: %d, maximum is %d per transaction", len(tags), MaxTagsPerTransaction)
	}
	for _, t := range tags {
		if err := ValidateTag(t); err != nil {
			return err
		}
	}
	return nil
}

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
// Up supports a maximum of 6 tags per transaction (MaxTagsPerTransaction) and
// tags of at most 30 characters (MaxTagLength). The API is the authority and
// rejects anything else, so tags are NOT validated here; use ValidateTags for an
// optional pre-flight check. Duplicate tags are silently ignored by the API.
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
