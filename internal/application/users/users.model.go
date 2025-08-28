package users

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserModel struct {
    ID            bson.ObjectID   `bson:"_id,omitempty"`
    Username      string               `bson:"username"`             // unique
    Email         string               `bson:"email"`                // unique, used for login/reset
    PasswordHash  string               `bson:"password_hash"`        // bcrypt hash
    FullName      string               `bson:"full_name,omitempty"`  // optional
    AvatarURL     string               `bson:"avatar_url,omitempty"` // optional
    RoleIDs       []bson.ObjectID     `bson:"role_ids,omitempty"`      
    PermissionIDs []bson.ObjectID     `bson:"permission_ids,omitempty"` 
    Status        string               `bson:"status"`               // active, inactive, suspended
    IsDeleted     bool                 `bson:"is_deleted"`           // soft delete
    CreatedAt     time.Time            `bson:"created_at"`
    UpdatedAt     time.Time            `bson:"updated_at"`
    LastLoginAt   *time.Time           `bson:"last_login_at,omitempty"`
    TenantID      *bson.ObjectID  `bson:"tenant_id,omitempty"`  // multi-tenant
}

type RoleModel struct {
    ID            bson.ObjectID   `bson:"_id,omitempty"`
    Name          string               `bson:"name"`                 // unique within tenant
    Description   string               `bson:"description,omitempty"`
    PermissionIDs []bson.ObjectID `bson:"permission_ids,omitempty"`
    IsDefault     bool                 `bson:"is_default"`           // default role for new users
    TenantID      *bson.ObjectID  `bson:"tenant_id,omitempty"`  // multi-tenant
    CreatedAt     time.Time            `bson:"created_at"`
    UpdatedAt     time.Time            `bson:"updated_at"`
}

type PermissionModel struct {
    ID          bson.ObjectID `bson:"_id,omitempty"`
    Name        string             `bson:"name"`                  // unique identifier
    Description string             `bson:"description,omitempty"`
    Resource    string             `bson:"resource,omitempty"`    // e.g., "user", "project"
    Action      string             `bson:"action,omitempty"`      // e.g., "read", "write", "delete"
    TenantID    *bson.ObjectID `bson:"tenant_id,omitempty"`  // optional multi-tenant
    CreatedAt   time.Time           `bson:"created_at"`
    UpdatedAt   time.Time           `bson:"updated_at"`
}

func UserModelToEntity(e UserModel) UserEntity {

  var tenant_id string
  if e.TenantID != nil {
    tenant_id = e.TenantID.Hex()
  }
  // var RolesParse []RoleEntity  
  // var PermissionsParse []PermissionEntity

  return UserEntity{
    ID:e.ID.Hex(),
    Username: e.Username,
    Email: e.Email,
    PasswordHash: e.PasswordHash,
    FullName: e.FullName,
    AvatarURL: e.AvatarURL,
    // Roles: RolesParse,
    // Permissions: PermissionsParse,
    Status: StatusUser(e.Status),
    IsDeleted: e.IsDeleted,
    CreatedAt: e.CreatedAt,
    UpdatedAt: e.UpdatedAt,
    LastLoginAt: e.LastLoginAt,
    TenantID: &tenant_id,
  }
}


