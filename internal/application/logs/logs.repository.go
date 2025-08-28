package logs

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogModel struct {
    ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
    RequestID string                 `bson:"request_id" json:"request_id"`
    Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
    Level     string                 `bson:"level" json:"level"`           // INFO, ERROR, WARN, DEBUG
    Service   string                 `bson:"service" json:"service"`       // เช่น "user-service"
    Module    string                 `bson:"module" json:"module"`         // เช่น "auth"
    Action    string                 `bson:"action" json:"action"`         // เช่น "login"
    Message   string                 `bson:"message" json:"message"`
    Error     string                 `bson:"error,omitempty" json:"error"` // optional
    Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata"`
}

