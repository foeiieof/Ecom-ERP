package users

import "time"

//-//go:generate stringer -type=StatusUser with bash go generate ./...
type StatusUser string  

const (
  StatusActive StatusUser   = "active"
  StatusInactive StatusUser = "inactive"
  StatusLocked StatusUser   = "locked"
) 

// var stateuser = map[StatusUser]string{
//   StatusActive : "active",
//   StatusInactive : "inactive",
//   StatusLocked : "locked",
// }


type UserEntity struct {
    ID            string
    Username      string
    Email         string
    PasswordHash  string
    FullName      string
    AvatarURL     string
    Roles         []RoleEntity
    Permissions   []PermissionEntity
    Status        StatusUser
    IsDeleted     bool
    TenantID      *string
    CreatedAt     time.Time
    UpdatedAt     time.Time
    LastLoginAt   *time.Time
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
