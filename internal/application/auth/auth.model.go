package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JWTClaims struct {
    UserID      bson.ObjectID `json:"user_id"`
    Roles       []string           `json:"roles"`
    Permissions []string           `json:"permissions"`
    jwt.RegisteredClaims
}

