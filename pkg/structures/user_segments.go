package structures

type UserSegments struct {
	UserId                  int      `json:"user_id" binding:"required"`
	SegmentsToAdd           []string `json:"segments_to_add" binding:"required"`
	SegmentsToAddExpiration *string  `json:"segments_to_add_expiration" example:"2023-08-30 12:00:00"`
	SegmentsToDelete        []string `json:"segments_to_delete" binding:"required"`
}
