package user

import "time"

type UserEntity struct {
    ID            string
    Username      string
    Email         string
    PasswordHash  string
    FullName      string
    AvatarURL     string
    Roles         []RoleEntity
    Permissions   []PermissionEntity
    Status        string
    IsDeleted     bool
    CreatedAt     time.Time
    UpdatedAt     time.Time
    LastLoginAt   *time.Time
    TenantID      *string
}

type RoleEntity struct {
    ID          string
    Name        string
    Description string
    Permissions []PermissionEntity
    IsDefault   bool
    TenantID    *string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type PermissionEntity struct {
    ID          string
    Name        string
    Description string
    Resource    string
    Action      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
