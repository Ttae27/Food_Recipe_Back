package routes

import (
	"fmt"
	"log"
	"github.com/Ttae27/Food_Recipe_Back/models"
	"gorm.io/gorm"
)

func createUser(db *gorm.DB,user *models.User) {
	result := db.Create(user)

	if result.Error != nil{
		log.Fatalf("Error creating user: %v",result.Error)
	}
	fmt.Println("Create User successfully")
}

func getUser(db *gorm.DB,id uint) *models.User{
	var user models.User;
	result := db.First(&user,id);

	if result.Error != nil{
		log.Fatalf("Error getting user: %v",result.Error);
	}
	return &user;
}

func updateUser(db *gorm.DB,user *models.User) {
	result := db.Save(user)

	if result.Error != nil{
		log.Fatalf("Error updating user: %v",result.Error)
	}
	fmt.Println("Update User successfully")
}

func deleteUser(db *gorm.DB,id uint) {
	var user models.User;
	result := db.Delete(&user,id)

	if result.Error != nil{
		log.Fatalf("Error updating user: %v",result.Error)
	}
	fmt.Println("Update User successfully")
}

