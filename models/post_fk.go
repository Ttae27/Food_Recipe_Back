package models

type Post_Like struct {
	PostID uint
	Post   Post
	UserID uint
	User   User
}

type Post_Comment struct {
	PostID    uint
	CommentID uint
	Comment   Comment
}

type Post_Ingredient struct {
	PostID       uint
	Post         Post
	IngredientID uint
	Ingredient   Ingredient
	Quantity     uint
}
