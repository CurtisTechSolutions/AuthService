package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Status   string
	Email    string
	Password string
	Role     string
}

func UserCreate(user *User) error {
	// Check if the user already exists
	UserGet(&User{Email: user.Email})

	result := DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UserExists(userQuery *User) (bool, error) {
	var count int64 = 0
	result := DB.Model(&User{}).Where(userQuery).Count(&count)
	if result.Error != nil && count == 0 {
		return false, result.Error
	}
	return count > 0, result.Error
}

func UserGet(userQuery *User) (*User, error) {
	var user User
	result := DB.Where(userQuery).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func UserUpdate(currentUser *User, newUser *User) {
	DB.Model(&currentUser).Updates(newUser)
}
