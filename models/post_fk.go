package models

type Post_Like struct {
	PostId uint
	Post   Post
	UserId uint
	User   User
}

type Post_Comment struct {
	PostId    uint
	Post      Post
	CommentId uint
	Comment   Comment
}

type Post_Material struct {
	PostId     uint
	Post       Post
	MaterialId uint
	Material   Material
}
