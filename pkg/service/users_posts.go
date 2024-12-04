package service

import (
	"github.com/SanyaWarvar/temple_api/pkg/models"
	"github.com/SanyaWarvar/temple_api/pkg/repository"
	"github.com/google/uuid"
)

type UsersPostsService struct {
	repo repository.IUsersPostsRepo
}

func NewUsersPostsService(repo repository.IUsersPostsRepo) *UsersPostsService {
	return &UsersPostsService{repo: repo}
}

func (s *UsersPostsService) CreatePost(post models.UserPost) (uuid.UUID, error) {
	post.Id = uuid.New()
	return post.Id, s.repo.CreatePost(post)
}
func (s *UsersPostsService) UpdatePost(newPost models.UserPost) error {
	return s.repo.UpdatePost(newPost)
}
func (s *UsersPostsService) GetPostById(postId, userId uuid.UUID) (repository.UserPostOutput, error) {
	return s.repo.GetPostById(postId, userId)
}
func (s *UsersPostsService) DeletePostById(postId, userId uuid.UUID) error {
	return s.repo.DeletePostById(postId, userId)
}

func (s *UsersPostsService) GetPostsByU(username string, page int, userId uuid.UUID) ([]repository.UserPostOutput, error) {
	return s.repo.GetPostsByU(username, page, userId)
}

func (s *UsersPostsService) LikePostById(postId, userId uuid.UUID) error {
	return s.repo.LikePostById(postId, userId)
}
