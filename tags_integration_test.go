//go:build integration

package up

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
)

// These tests talk to the REAL Up API and temporarily add tags to a real
// transaction, removing them again afterwards. They check that the advisory
// ValidateTag/ValidateTags helpers still agree with what the API enforces, so
// they will tell you if Up ever changes its tag limits.
//
// Run with:
//
//	UP_INTEGRATION_TOKEN=... UP_INTEGRATION_TRANSACTION_ID=... \
//	  go test -tags integration -run Test_Integration -v .
//
// Use a transaction that accepts tags (an ordinary card purchase) and has few
// tags already; the cap test needs room for at least one throwaway tag.

func integrationClient(t *testing.T) (*Client, string) {
	t.Helper()
	token, id := os.Getenv("UP_INTEGRATION_TOKEN"), os.Getenv("UP_INTEGRATION_TRANSACTION_ID")
	if token == "" || id == "" {
		t.Skip("set UP_INTEGRATION_TOKEN and UP_INTEGRATION_TRANSACTION_ID to run")
	}
	c, err := New(context.Background(), token)
	if err != nil {
		t.Fatalf("New() returned an error: %v", err)
	}
	return c, id
}

// cleanupTags removes tags when the test ends and fails the test loudly if a
// throwaway tag could not be removed, so nothing is silently left behind.
func cleanupTags(t *testing.T, c *Client, id string, tags ...string) {
	t.Helper()
	t.Cleanup(func() {
		for _, tag := range tags {
			// the API trims and collapses whitespace, so remove that form too.
			for _, form := range []string{tag, strings.Join(strings.Fields(tag), " ")} {
				if strings.TrimSpace(form) == "" {
					continue
				}
				_ = c.RemoveTagsFromTransaction(context.Background(), id, []string{form})
			}
		}
		got, err := c.GetTransaction(context.Background(), id)
		if err != nil {
			t.Errorf("could not verify cleanup: %v", err)
			return
		}
		for _, o := range got.Relationships.Tags.Data {
			if strings.HasPrefix(o.ID, "up-go-it") {
				t.Errorf("throwaway tag %q was left on transaction %s; remove it by hand", o.ID, id)
			}
		}
	})
}

func Test_Integration_TagRules(t *testing.T) {
	c, id := integrationClient(t)
	ctx := context.Background()

	tests := map[string]string{
		"30 characters":          strings.Repeat("a", 30),
		"31 characters":          strings.Repeat("a", 31),
		"30 multibyte":           strings.Repeat("é", 30),
		"31 multibyte":           strings.Repeat("é", 31),
		"empty":                  "",
		"whitespace only":        "   ",
		"symbols and spaces":     "up-go-it a & b / c: d,e",
		"emoji":                  "up-go-it 🙂",
		"30 plus trailing space": strings.Repeat("b", 30) + " ",
	}
	for name, tag := range tests {
		t.Run(name, func(t *testing.T) {
			cleanupTags(t, c, id, tag)
			apiErr := c.AddTagsToTransaction(ctx, id, []string{tag})
			helperErr := ValidateTag(tag)
			if (apiErr == nil) != (helperErr == nil) {
				t.Errorf("tag %q: the API accepted it = %v, but ValidateTag accepted it = %v; "+
					"the advisory limits in tags.go no longer match Up (API error: %v, helper error: %v)",
					tag, apiErr == nil, helperErr == nil, apiErr, helperErr)
			}
		})
	}
}

func Test_Integration_TagCap(t *testing.T) {
	c, id := integrationClient(t)
	ctx := context.Background()

	got, err := c.GetTransaction(ctx, id)
	if err != nil {
		t.Fatalf("GetTransaction() returned an error: %v", err)
	}
	room := MaxTagsPerTransaction - len(got.Relationships.Tags.Data)
	if room < 1 {
		t.Skipf("transaction %s already has %d tags; pick one with room", id, len(got.Relationships.Tags.Data))
	}

	var added []string
	for i := 0; i < room+1; i++ {
		added = append(added, fmt.Sprintf("up-go-it-cap-%d", i))
	}
	cleanupTags(t, c, id, added...)

	for i := 0; i < room; i++ {
		if err := c.AddTagsToTransaction(ctx, id, []string{added[i]}); err != nil {
			t.Fatalf("adding tag %d of %d (still within the cap of %d) was rejected: %v", i+1, room, MaxTagsPerTransaction, err)
		}
	}
	if err := c.AddTagsToTransaction(ctx, id, []string{added[room]}); err == nil {
		t.Errorf("adding a tag beyond the cap of %d was accepted; MaxTagsPerTransaction is out of date", MaxTagsPerTransaction)
	}
}
