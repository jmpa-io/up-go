package up

// "context"
// "encoding/json"
// "net/http"
// "net/http/httptest"
// "reflect"
// "strings"
// "testing"
// "time"

// // setup expected payloads.
// var (
// 	location                     *time.Location
// 	singleTransactionPayload     WrapperOmittable
// 	paginatedTransactionsPayload TransactionWrapper
// )
//
// func init() {
//
// 	// load location.
// 	var err error
// 	location, err = time.LoadLocation("Australia/Sydney")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// populate single transaction payload.
// 	singleTransactionPayload = WrapperOmittable{
// 		Data: Transaction{
// 			Object: Object{
// 				Type: "transactions",
// 				ID:   "6fff09f5-be7d-4ae1-9f71-4a25440bc405",
// 			},
// 			Attributes: TransactionResource{
// 				Status:          TransactionStatusSettled,
// 				RawText:         "WARUNG BEBEK, UBUD INDONES",
// 				Description:     "Warung Bebek Bengil",
// 				IsCategorizable: true,
// 				HoldInfo: TransactionResourceHoldInfo{
// 					Amount: Money{
// 						CurrencyCode:     "AUD",
// 						Value:            "-107.92",
// 						ValueInBaseUnits: -10792,
// 					},
// 				},
// 				RoundUp: TransactionResourceRoundUp{
// 					Amount: Money{
// 						CurrencyCode:     "AUD",
// 						Value:            "-0.08",
// 						ValueInBaseUnits: -8,
// 					},
// 				},
// 				Amount: Money{
// 					CurrencyCode:     "AUD",
// 					Value:            "-107.92",
// 					ValueInBaseUnits: -10792,
// 				},
// 				ForeignAmount: Money{
// 					CurrencyCode:     "IDR",
// 					Value:            "-1053698.77",
// 					ValueInBaseUnits: -105369877,
// 				},
// 				CardPurchaseMethod: TransactionResourceCardPurchaseMethod{
// 					Method:           CardPurchaseMethodCardOnFile,
// 					CardNumberSuffix: "0001",
// 				},
// 				SettledAt: time.Date(2024, 11, 03, 04, 00, 00, 00, location),
// 				CreatedAt: time.Date(2024, 11, 03, 04, 00, 00, 00, location),
// 				PerformingCustomer: TransactionResourcePerformingCustomer{
// 					DisplayName: "Bobby",
// 				},
// 				DeepLinkURL: "up://transaction/VHJhbnNhY3Rpb24tNDY=",
// 			},
// 			Relationships: TransactionRelationships{
// 				Account: Wrapper[Object]{
// 					Data: Object{
// 						Type: "accounts",
// 						ID:   "fe848390-7d39-41fb-b01d-545de29ab74b",
// 					},
// 					Links: Links{
// 						Related: "https://api.up.com.au/api/v1/accounts/fe848390-7d39-41fb-b01d-545de29ab74b",
// 					},
// 				},
// 				Category: Wrapper[Object]{
// 					Links: Links{
// 						Self: "https://api.up.com.au/api/v1/transactions/6fff09f5-be7d-4ae1-9f71-4a25440bc405/relationships/category",
// 					},
// 				},
// 				Tags: WrapperSlice[Object]{
// 					Links: Links{
// 						Self: "https://api.up.com.au/api/v1/transactions/6fff09f5-be7d-4ae1-9f71-4a25440bc405/relationships/tags",
// 					},
// 				},
// 			},
// 		},
// 		Links: Links{
// 			Self: "https://api.up.com.au/api/v1/transactions/6fff09f5-be7d-4ae1-9f71-4a25440bc405",
// 		},
// 	}
// 	// populate paginated transactions payload.
// 	paginatedTransactionsPayload = TransactionWrapper{
// 		Data: []Transaction{
// 			{
// 				Object: Object{
// 					Type: "transaction",
// 					ID:   "b015b19b-32f0-41e0-9817-d5fb22fd259e",
// 				},
// 				Attributes: TransactionResource{
// 					Status:          "SETTLED",
// 					RawText:         "",
// 					Description:     "David Taylor",
// 					Message:         "Money for the pizzas last night.",
// 					IsCategorizable: true,
// 					Amount: Money{
// 						CurrencyCode:     "AUD",
// 						Value:            "-59.98",
// 						ValueInBaseUnits: -5998,
// 					},
// 					SettledAt: time.Date(2024, 11, 5, 7, 25, 12, 0, location),
// 					CreatedAt: time.Date(2024, 11, 5, 7, 25, 12, 0, location),
// 					PerformingCustomer: TransactionResourcePerformingCustomer{
// 						DisplayName: "Bobby",
// 					},
// 					DeepLinkURL: "up://transaction/VHJhbnNhY3Rpb24tMzg=",
// 				},
// 				Relationships: TransactionRelationships{
// 					Account: Wrapper[Object]{
// 						Data: Object{
// 							Type: "accounts",
// 							ID:   "420e0a2d-1e5e-4278-a194-219d3fb5c575",
// 						},
// 						Links: Links{
// 							Related: "https://api.up.com.au/api/v1/accounts/420e0a2d-1e5e-4278-a194-219d3fb5c575",
// 						},
// 					},
// 					TransferAccount: Wrapper[Object]{},
// 					Category: Wrapper[Object]{
// 						Links: Links{
// 							Related: "https://api.up.com.au/api/v1/transactions/b015b19b-32f0-41e0-9817-d5fb22fd259e/relationships/category",
// 						},
// 					},
// 					ParentCategory: Wrapper[Object]{},
// 					Tags: WrapperSlice[Object]{
// 						Data: []Object{
// 							{
// 								Type: "tags",
// 								ID:   "Pizza Night",
// 							},
// 						},
// 						Links: Links{
// 							Self: "https://api.up.com.au/api/v1/transactions/b015b19b-32f0-41e0-9817-d5fb22fd259e/relationships/tags",
// 						},
// 					},
// 					Attachment: Wrapper[Object]{},
// 				},
// 				Links: Links{
// 					Self: "https://api.up.com.au/api/v1/transactions/b015b19b-32f0-41e0-9817-d5fb22fd259e",
// 				},
// 			},
// 		},
// 	}
// }
//
// func Test_GetTransactions(t *testing.T) {
// 	tests := map[string]struct {
// 		server   *httptest.Server
// 		testdata string
// 		want     interface{}
// 		err      string
// 	}{
// 		"read transaction": {
// 			server: httptest.NewUnstartedServer(
// 				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 					b, _ := json.Marshal(singleTransactionPayload)
// 					w.WriteHeader(http.StatusOK)
// 					w.Write(b)
// 				}),
// 			),
// 			testdata: "testdata/transaction.json",
// 			want:     Transaction{},
// 		},
// 		// "read paginated transactions": {
// 		// 	server: httptest.NewUnstartedServer(
// 		// 		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// 			out := paginatedTransactionsPayload
// 		// 			b, _ := json.Marshal(out)
// 		// 			w.WriteHeader(http.StatusOK)
// 		// 			w.Write(b)
// 		// 		}),
// 		// 	),
// 		// 	testdata: "testdata/transactions.json",
// 		// 	want:     []Transaction{},
// 		// },
// 		// "unauthorized transaction": {
// 		// 	server: httptest.NewUnstartedServer(
// 		// 		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// 			out := apiErrorResponse{
// 		// 				Errors: []apiErrorResponseError{
// 		// 					{
// 		// 						Status: "TODO",
// 		// 						Title:  "TODO",
// 		// 						Detail: "TODO",
// 		// 						Source: apiErrorResponseErrorSource{
// 		// 							Parameter: "TODO",
// 		// 						},
// 		// 					},
// 		// 				},
// 		// 			}
// 		// 			b, _ := json.Marshal(out)
// 		// 			w.WriteHeader(http.StatusUnauthorized)
// 		// 			w.Write(b)
// 		// 		}),
// 		// 	),
// 		// 	testdata: "testdata/unauthorized.json",
// 		// 	err:      "error response returned from API",
// 		// },
// 	}
// 	for name, tt := range tests {
//
// 		// tracing context.
// 		ctx := context.Background()
//
// 		// start test server.
// 		tt.server.Start()
// 		defer tt.server.Close()
//
// 		// read + parse testdata.
// 		var td Transaction
// 		if err := readTestdata(tt.testdata, &td); err != nil {
// 			panic(err)
// 		}
// 		tt.want = td.Attributes
//
// 		// setup client with test server.
// 		c, _ := New(ctx, "xxxx",
// 			WithHttpClient(tt.server.Client()),
// 			WithEndpoint(tt.server.URL),
// 		)
//
// 		// run tests.
// 		t.Run(name, func(t *testing.T) {
// 			got, err := c.ListTransactions(context.Background())
// 			if tt.err != "" && err != nil {
// 				if !strings.Contains(err.Error(), tt.err) {
// 					t.Errorf(
// 						"GetTransactions() returned an unexpected error;\nwant=%v\ngot=%v\n",
// 						tt.err,
// 						err,
// 					)
// 				}
// 				return
// 			}
// 			if err != nil {
// 				t.Errorf("GetTransactions() returned an error;\nerror=%v\n", err)
// 				return
// 			}
// 			switch {
// 			case
// 				!reflect.DeepEqual(got, tt.want):
// 				t.Errorf(
// 					"GetTransactions() returned unexpected configuration;\nwant=%+v\ngot=%+v\n",
// 					tt.want,
// 					got,
// 				)
// 				return
// 			}
// 		})
// 	}
// }
