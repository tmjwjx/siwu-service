package requests

import "forum/internal/models"

type UserRequest struct {
	User       models.User       `json:"user"`
	UserDetail models.UserDetail `json:"user_detail"`
	UserTags   []string          `json:"user_tags"`
}

type UserResponse struct {
	User        models.User       `json:"user"`
	UserDetail  models.UserDetail `json:"user_detail"`
	UserTags    []string          `json:"user_tags"`
	AllTagNames []string          `json:"all_tag_names"`
	Path        string            `json:"path"`
}
