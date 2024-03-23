package domain

type PostRepository interface {
	GetPost(id int) (*Post, error)
	GetPosts() ([]*Post, error)
	GetPostBySlug(slug string) (*Post, error)
	CreatePost(post *PostCreate) (*Post, error)
	UpdatePost(post *Post) (*Post, error)
}
