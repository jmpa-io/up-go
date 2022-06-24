package up

type LinksObject struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
}

type RelatedLinksObject struct {
	Related string `json:"related"`
}

type SelfLinkObject struct {
	Self string `json:"self"`
}
