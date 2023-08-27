package structures

type User struct {
	Id int `json:"-" binding:"required"`
}
