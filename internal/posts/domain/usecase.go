package domain

type PostUsecase interface {
	Get(id int) (*Post, error)
	GetPosts() ([]*Post, error)
	GetBySlug(slug string) (*Post, error)
	Create(post *PostCreate)
	Update(id int, post *PostUpdate)
}
