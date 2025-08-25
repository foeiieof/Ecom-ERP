package auth

import (
  "go.mongodb.org/mongo-driver/bson/primitive"
  "github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
    UserID      primitive.ObjectID `json:"user_id"`
    Roles       []string           `json:"roles"`
    Permissions []string           `json:"permissions"`
    jwt.RegisteredClaims
}

