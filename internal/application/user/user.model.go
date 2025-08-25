package user

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
    ID            primitive.ObjectID   `bson:"_id,omitempty"`
    Username      string               `bson:"username"`             // unique
    Email         string               `bson:"email"`                // unique, used for login/reset
    PasswordHash  string               `bson:"password_hash"`        // bcrypt hash
    FullName      string               `bson:"full_name,omitempty"`  // optional
    AvatarURL     string               `bson:"avatar_url,omitempty"` // optional
    RoleIDs       []primitive.ObjectID `bson:"role_ids,omitempty"`      
    PermissionIDs []primitive.ObjectID `bson:"permission_ids,omitempty"` 
    Status        string               `bson:"status"`               // active, inactive, suspended
    IsDeleted     bool                 `bson:"is_deleted"`           // soft delete
    CreatedAt     time.Time            `bson:"created_at"`
    UpdatedAt     time.Time            `bson:"updated_at"`
    LastLoginAt   *time.Time           `bson:"last_login_at,omitempty"`
    // TenantID      *primitive.ObjectID  `bson:"tenant_id,omitempty"`  // multi-tenant
}

type RoleModel struct {
    ID            primitive.ObjectID   `bson:"_id,omitempty"`
    Name          string               `bson:"name"`                 // unique within tenant
    Description   string               `bson:"description,omitempty"`
    PermissionIDs []primitive.ObjectID `bson:"permission_ids,omitempty"`
    IsDefault     bool                 `bson:"is_default"`           // default role for new users
    TenantID      *primitive.ObjectID  `bson:"tenant_id,omitempty"`  // multi-tenant
    CreatedAt     time.Time            `bson:"created_at"`
    UpdatedAt     time.Time            `bson:"updated_at"`
}

type PermissionModel struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Name        string             `bson:"name"`                  // unique identifier
    Description string             `bson:"description,omitempty"`
    Resource    string             `bson:"resource,omitempty"`    // e.g., "user", "project"
    Action      string             `bson:"action,omitempty"`      // e.g., "read", "write", "delete"
    // TenantID    *primitive.ObjectID `bson:"tenant_id,omitempty"`  // optional multi-tenant
    CreatedAt   time.Time           `bson:"created_at"`
    UpdatedAt   time.Time           `bson:"updated_at"`
}
