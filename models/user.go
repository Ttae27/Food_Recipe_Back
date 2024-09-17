package models

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string
	Password       string
	FirstName      string
	LastName       string
	PictureProfile string
	Bookmarks      []Post `gorm:"many2many:Bookmark;"`
}

func createUser(db *gorm.DB,user *User) {
	result := db.Create(user)

	if result.Error != nil{
		log.Fatalf("Error creating user: %v",result.Error)
	}
	fmt.Println("Create User successfully")
}

func getUser(db *gorm.DB,id uint) *User{
	var user User;
	result := db.First(&user,id);

	if result.Error != nil{
		log.Fatalf("Error getting user: %v",result.Error);
	}
	return &user;
}

func updateUser(db *gorm.DB,user *User) {
	result := db.Save(user)

	if result.Error != nil{
		log.Fatalf("Error updating user: %v",result.Error)
	}
	fmt.Println("Update User successfully")
}

func deleteUser(db *gorm.DB,id uint) {
	var user User;
	result := db.Delete(&user,id)

	if result.Error != nil{
		log.Fatalf("Error updating user: %v",result.Error)
	}
	fmt.Println("Update User successfully")
}