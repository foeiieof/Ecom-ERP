package users

import "time"

type IReqCreateUserDTO struct {
    Username string `json:"username" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
    AvatarURL string `json:"avatar_url,omitempty"`
    FullName string `json:"full_name,omitempty"`
    Status   string    `json:"status"`
    
}

type IReqUpdateUserDTO struct {
    Username  string `json:"username" validate:"required"`
    FullName  string `json:"full_name,omitempty"`
    AvatarURL string `json:"avatar_url,omitempty"`
    Status    string `json:"status,omitempty"`
}

// Response DTOs
type UserDTO struct {
    ID        string    `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    FullName  string    `json:"full_name,omitempty"`
    AvatarURL string    `json:"avatar_url,omitempty"`
    Roles     []string  `json:"roles,omitempty"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    LastLogin *time.Time `json:"last_login_at,omitempty"`
}
