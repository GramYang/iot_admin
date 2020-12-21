package sqlx_client

type Admin struct {
	Id int `json:"id" db:"id"`
	UserName string `json:"userName" db:"user_name"`
	Password string `json:"password" db:"password"`
	Privilege int `json:"privilege" db:"privilege"`
	UpdateTime string `json:"updateTime" db:"update_time"`
}