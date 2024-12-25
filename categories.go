package up

type CategoryAttributes struct {
	Name string `json:"name"`
}

type CategoryRelationships struct {
	Parent struct {
	} `json:"parent"`
	Child struct {
	} `json:"child"`
	Link RelatedLink `json:"links"`
}

// CategoryDataWrapper represents a category in Up.
type CategoryDataWrapper Data[CategoryAttributes, CategoryRelationships]

// CategoryPaginationWrapper is a pagination wrapper for a slice of CategoryData
type CategoryPagenationWrapper PaginationWrapper[CategoryDataWrapper]
