package structures

type User struct {
	Id int `json:"user_id" binding:"required"`
}

type UserHistory struct {
	Id        int    `json:"user_id" binding:"required"`
	YearMonth string `json:"year_month" binding:"required" example:"YYYY-MM"`
}
