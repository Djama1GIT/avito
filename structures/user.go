package structures

type User struct {
	Id int `json:"user_id" binding:"required"`
}
