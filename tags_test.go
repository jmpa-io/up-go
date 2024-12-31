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

		// setup client with mock.
		c, _ := New(ctx, "xxxx",
			WithHttpClient(&http.Client{
				Transport: tt.mock,
			}),
		)

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

		// setup client with mock.
		c, _ := New(ctx, "xxxx",
			WithHttpClient(&http.Client{
				Transport: tt.mock,
			}),
		)

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

		// setup client with mock.
		c, _ := New(ctx, "xxxx",
			WithHttpClient(&http.Client{
				Transport: tt.mock,
			}),
		)

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
