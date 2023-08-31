package structures

type Segment struct {
	Slug       string `json:"slug" binding:"required"`
	Percentage *int   `json:"percentage"`
}
