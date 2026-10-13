package up

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
)

var tagsTestdata []*testdata

func init() {

	// populate testdata.
	for i := 1; i <= 3; i++ {
		tagsTestdata = append(tagsTestdata, newTestdata(fmt.Sprintf("tags-%v", i)))
	}
}

func Test_ListTags(t *testing.T) {
	tests := map[string]struct {
		mock *mockRoundTripper
		want []TagResource
		err  string
	}{
		"list tags": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					b := tagsTestdata[0].content
					for i := 0; i < len(tagsTestdata); i++ {
						if !strings.Contains(req.URL.String(), fmt.Sprintf("---%v", i+1)) {
							continue
						}
						b = tagsTestdata[i].content
						break
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(b)),
						Header:     make(http.Header),
					}
				},
			},
			want: []TagResource{
				{Object: Object{Type: "tags", ID: "Holiday"}},
				{Object: Object{Type: "tags", ID: "Pizza Night"}},
				{Object: Object{Type: "tags", ID: "Dining Out"}},
				{Object: Object{Type: "tags", ID: "Shopping"}},
				{Object: Object{Type: "tags", ID: "Fitness"}},
			},
		},
	}
	for name, tt := range tests {

		// tracing context.
		ctx := context.Background()

		// setup client.
		c := newTestClient(t, tt.mock)

		// run tests.
		t.Run(name, func(t *testing.T) {
			got, err := c.ListTags(ctx)

			// any errors?
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf(
						"ListTags() returned an unexpected error;\nwant=%v\ngot=%v\n",
						tt.err,
						err,
					)
				}
				return
			}
			if err != nil {
				t.Errorf("ListTags() returned an error;\nerror=%v\n", err)
				return
			}

			// do the lengths match?
			if len(got) != len(tt.want) {
				t.Errorf(
					"ListTags() returned unexpected number of results;\nwant=%d\ngot=%d\n",
					len(tt.want),
					len(got),
				)
				return
			}

			// is there a mismatch from what we're expecting vs what we've got?
			if !slices.Equal(got, tt.want) {
				t.Errorf(
					"ListTags() returned unexpected configuration;\nwant=%+v\ngot=%+v\n",
					tt.want,
					got,
				)
			}
		})
	}
}

func Test_AddTagsToTransaction(t *testing.T) {
	tests := map[string]struct {
		mock *mockRoundTripper
		id   string
		tags []string
		err  string
	}{
		"add tags": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					var b []byte
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(b)),
						Header:     make(http.Header),
					}
				},
			},
			id:   "1",
			tags: []string{"hello", "world"},
		},
	}
	for name, tt := range tests {

		// tracing context.
		ctx := context.Background()

		// setup client.
		c := newTestClient(t, tt.mock)

		// run tests.
		t.Run(name, func(t *testing.T) {
			err := c.AddTagsToTransaction(ctx, tt.id, tt.tags)
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf(
						"AddTagsToTransaction() returned an unexpected error;\nwant=%v\ngot=%v\n",
						tt.err,
						err,
					)
				}
				return
			}
			if err != nil {
				t.Errorf("AddTagsToTransaction() returned an error;\nerror=%v\n", err)
				return
			}
		})
	}
}

func Test_RemoveTagsFromTransaction(t *testing.T) {
	tests := map[string]struct {
		mock *mockRoundTripper
		id   string
		tags []string
		err  string
	}{
		"remove tags": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					var b []byte
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(b)),
						Header:     make(http.Header),
					}
				},
			},
			tags: []string{"hello", "world"},
		},
	}
	for name, tt := range tests {

		// tracing context.
		ctx := context.Background()

		// setup client.
		c := newTestClient(t, tt.mock)

		// run tests.
		t.Run(name, func(t *testing.T) {
			err := c.RemoveTagsFromTransaction(ctx, tt.id, tt.tags)
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf(
						"RemoveTagsToTransaction() returned an unexpected error;\nwant=%v\ngot=%v\n",
						tt.err,
						err,
					)
				}
				return
			}
			if err != nil {
				t.Errorf("RemoveTagsToTransaction() returned an error;\nerror=%v\n", err)
				return
			}
		})
	}
}

// Test_ValidateTag pins the advisory helper to the rules observed on the live Up
// API (checked 2026-10-13; see the integration test): a tag is at most 30 characters (not bytes), counted on the raw
// string, and must not be empty or whitespace-only.
func Test_ValidateTag(t *testing.T) {
	tests := map[string]struct {
		tag string
		err string
	}{
		"simple":                  {tag: "groceries"},
		"namespaced":              {tag: "ps-cat:Food"},
		"spaces and symbols":      {tag: "Transport & Auto / a,b \"q\""},
		"exactly 30 characters":   {tag: strings.Repeat("x", 30)},
		"31 characters":           {tag: strings.Repeat("x", 31), err: "31 characters, maximum is 30"},
		"30 multibyte characters": {tag: strings.Repeat("é", 30)}, // 60 bytes: length is in characters.
		"31 multibyte characters": {tag: strings.Repeat("é", 31), err: "maximum is 30"},
		"emoji":                   {tag: "🙂ok"},
		"leading space counts":    {tag: " " + strings.Repeat("y", 29)},
		"30 chars plus trailing":  {tag: strings.Repeat("y", 30) + " ", err: "maximum is 30"},
		"empty":                   {tag: "", err: "must not be empty"},
		"whitespace only":         {tag: "  \t", err: "must not be empty"},
		"long namespaced":         {tag: "pocketsmith-category:Transport & Auto", err: "maximum is 30"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateTag(tt.tag)
			switch {
			case tt.err == "" && err != nil:
				t.Errorf("ValidateTag(%q) returned an unexpected error: %v", tt.tag, err)
			case tt.err != "" && err == nil:
				t.Errorf("ValidateTag(%q) returned no error; want one containing %q", tt.tag, tt.err)
			case tt.err != "" && !strings.Contains(err.Error(), tt.err):
				t.Errorf("ValidateTag(%q) error = %v; want it to contain %q", tt.tag, err, tt.err)
			}
		})
	}
}

func Test_ValidateTags(t *testing.T) {
	six := []string{"a", "b", "c", "d", "e", "f"}
	tests := map[string]struct {
		tags []string
		err  string
	}{
		"none":               {tags: nil},
		"one":                {tags: []string{"a"}},
		"six is the maximum": {tags: six},
		"seven is too many":  {tags: append(append([]string{}, six...), "g"), err: "too many tags: 7, maximum is 6"},
		"one bad tag":        {tags: []string{"ok", strings.Repeat("x", 31)}, err: "maximum is 30"},
		"empty tag":          {tags: []string{"ok", ""}, err: "must not be empty"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateTags(tt.tags)
			switch {
			case tt.err == "" && err != nil:
				t.Errorf("ValidateTags() returned an unexpected error: %v", err)
			case tt.err != "" && err == nil:
				t.Errorf("ValidateTags() returned no error; want one containing %q", tt.err)
			case tt.err != "" && !strings.Contains(err.Error(), tt.err):
				t.Errorf("ValidateTags() error = %v; want it to contain %q", err, tt.err)
			}
		})
	}
}

// Test_AddTagsToTransaction_apiIsTheAuthority checks that the SDK does not
// second-guess the API: even tags ValidateTag would reject are sent, and the
// API's own error is returned to the caller.
func Test_AddTagsToTransaction_apiIsTheAuthority(t *testing.T) {
	tooLong := strings.Repeat("x", 31)
	calls := 0
	mock := &mockRoundTripper{
		MockFunc: func(req *http.Request) *http.Response {
			calls++
			return &http.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Body:       io.NopCloser(strings.NewReader(`{"errors":[{"status":"422","title":"Invalid Request Payload","detail":"tag too long"}]}`)),
				Header:     make(http.Header),
			}
		},
	}
	c := newTestClient(t, mock)
	err := c.AddTagsToTransaction(context.Background(), "1", []string{tooLong})
	if err == nil {
		t.Fatal("AddTagsToTransaction() returned no error; want the API's rejection")
	}
	if calls != 1 {
		t.Errorf("AddTagsToTransaction() made %d request(s); want 1 (the SDK must not pre-block)", calls)
	}
	if !strings.Contains(err.Error(), "Invalid Request Payload") {
		t.Errorf("AddTagsToTransaction() error = %v; want the API's message to be surfaced", err)
	}
}

// RemoveTagsFromTransaction is deliberately not validated: removing a tag that
// cannot exist is a harmless no-op, and it lets callers clean up odd tags.
func Test_RemoveTagsFromTransaction_notValidated(t *testing.T) {
	calls := 0
	mock := &mockRoundTripper{
		MockFunc: func(req *http.Request) *http.Response {
			calls++
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(nil)),
				Header:     make(http.Header),
			}
		},
	}
	c := newTestClient(t, mock)
	if err := c.RemoveTagsFromTransaction(context.Background(), "1", []string{strings.Repeat("x", 40)}); err != nil {
		t.Fatalf("RemoveTagsFromTransaction() returned an error: %v", err)
	}
	if calls != 1 {
		t.Errorf("RemoveTagsFromTransaction() made %d request(s); want 1", calls)
	}
}
