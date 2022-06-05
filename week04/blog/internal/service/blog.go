package service

import (
	"context"
	"fmt"

	v1 "github.com/blog/api/blog/v1"
)

type BlogService struct {
	v1.UnimplementedBlogServer
}

func (s *BlogService) GetArticle(ctx context.Context, in *v1.GetArticleRequest) (*v1.GetArticleReply, error) {
	return &v1.GetArticleReply{Message: fmt.Sprintf("GetArticleRequest Id:%s", in.Id)}, nil
}

func NewBlogService() *BlogService {
	return &BlogService{}
}
