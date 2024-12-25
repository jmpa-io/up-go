package up

type RelatedLink struct {
	Related string `json:"related"`
}

type SelfLink struct {
	Self string `json:"self"`
}

type PaginationLink struct {
	Next string `json:"next"`
	Prev string `json:"prev"`
}
