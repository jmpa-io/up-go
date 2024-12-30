package up

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

var (
	tagsTestdata1 = newTestdata("tags-1")
	tagsTestdata2 = newTestdata("tags-2")
	tagsTestdata3 = newTestdata("tags-3")
)

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
					switch {
					case strings.Contains(req.URL.String(), "---2"):
						b = tagsTestdata2.content
					case strings.Contains(req.URL.String(), "---3"):
						b = tagsTestdata3.content
					default:
						b = tagsTestdata1.content
					}
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
		"add tags": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					var b []byte
					switch {
					case strings.Contains(req.URL.String(), "---2"):
						b = tagsTestdata2.content
					case strings.Contains(req.URL.String(), "---3"):
						b = tagsTestdata3.content
					default:
						b = tagsTestdata1.content
					}
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
