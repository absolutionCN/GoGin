package models

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
	Role     string `gorm:"column:role" json:"role"`
	Avatar   string `gorm:"column:avatar" json:"avatar"`
}

func CheckAuthToken(username, password string) *Auth {
	var auth Auth
	db.Where(Auth{Username: username, Password: password}).First(&auth)
	if auth.ID > 0 {
		return &auth
	}
	return nil
}

func GetUserInfoByName(name string) *Auth {
	var auth Auth
	db.Where("username = ?", name).First(&auth)
	return &auth
}
