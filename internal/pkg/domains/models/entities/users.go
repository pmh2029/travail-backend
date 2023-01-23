package entities

// UsersTableName TableName
var UsersTableName = "users"

type User struct {
	BaseEntity
	Username string `gorm:"column:username;not null;unique"`
	Email    string `gorm:"column:email;not null;unique"`
	Password string `gorm:"column:password;not null"`
}

// TableName func
func (i *User) TableName() string {
	return UsersTableName
}
