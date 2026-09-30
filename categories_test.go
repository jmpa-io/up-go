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
	categoriesTestdata = newTestdata("categories")
	categoryTestdata   = newTestdata("category")
)

func Test_ListCategories(t *testing.T) {
	tests := map[string]struct {
		mock      *mockRoundTripper
		wantCount int
		err       string
	}{
		"list categories": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(categoriesTestdata.content)),
						Header:     make(http.Header),
					}
				},
			},
			wantCount: 2,
		},
	}
	for name, tt := range tests {
		ctx := context.Background()
		c := newTestClient(t, tt.mock)
		t.Run(name, func(t *testing.T) {
			got, err := c.ListCategories(ctx)
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("ListCategories() returned an unexpected error;\nwant=%v\ngot=%v\n", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("ListCategories() returned an error;\nerror=%v\n", err)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("ListCategories() returned unexpected count;\nwant=%d\ngot=%d\n", tt.wantCount, len(got))
			}
		})
	}
}

func Test_GetCategory(t *testing.T) {
	tests := map[string]struct {
		mock     *mockRoundTripper
		id       string
		wantID   string
		wantName string
		err      string
	}{
		"get category": {
			id:       "home",
			wantID:   "home",
			wantName: "Home",
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(categoryTestdata.content)),
						Header:     make(http.Header),
					}
				},
			},
		},
	}
	for name, tt := range tests {
		ctx := context.Background()
		c := newTestClient(t, tt.mock)
		t.Run(name, func(t *testing.T) {
			got, err := c.GetCategory(ctx, tt.id)
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("GetCategory() returned an unexpected error;\nwant=%v\ngot=%v\n", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("GetCategory() returned an error;\nerror=%v\n", err)
				return
			}
			if got.ID != tt.wantID || got.Attributes.Name != tt.wantName {
				t.Errorf("GetCategory() returned unexpected result;\nwantID=%s wantName=%s\ngotID=%s gotName=%s\n",
					tt.wantID, tt.wantName, got.ID, got.Attributes.Name)
			}
		})
	}
}

func Test_SetTransactionCategory(t *testing.T) {
	tests := map[string]struct {
		mock          *mockRoundTripper
		transactionID string
		categoryID    string
		err           string
	}{
		"set category": {
			transactionID: "6fff09f5-be7d-4ae1-9f71-4a25440bc405",
			categoryID:    "home",
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusNoContent,
						Body:       io.NopCloser(bytes.NewBuffer(nil)),
						Header:     make(http.Header),
					}
				},
			},
		},
		"de-categorise": {
			transactionID: "6fff09f5-be7d-4ae1-9f71-4a25440bc405",
			categoryID:    "",
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusNoContent,
						Body:       io.NopCloser(bytes.NewBuffer(nil)),
						Header:     make(http.Header),
					}
				},
			},
		},
	}
	for name, tt := range tests {
		ctx := context.Background()
		c := newTestClient(t, tt.mock)
		t.Run(name, func(t *testing.T) {
			err := c.SetTransactionCategory(ctx, tt.transactionID, tt.categoryID)
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("SetTransactionCategory() returned an unexpected error;\nwant=%v\ngot=%v\n", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("SetTransactionCategory() returned an error;\nerror=%v\n", err)
			}
		})
	}
}
