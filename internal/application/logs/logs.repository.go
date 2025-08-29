package logs

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type LogDetail string 

const (
  INFO  LogDetail = "INFO" 
  Error LogDetail = "ERROR"
  Warn  LogDetail = "WARN"
  Debug LogDetail = "DEBUG"
)

type LogModel struct {
  ID        bson.ObjectID     `bson:"_id,omitempty" json:"id"`
  RequestID string                 `bson:"request_id" json:"request_id"`
  Method    string                 `bson:"method" json:"method"`
  Service   string                 `bson:"service" json:"service"`       // เช่น "user-service"
  Level     LogDetail              `bson:"level" json:"level"`           // INFO, ERROR, WARN, DEBUG
  Endpoint  string                 `bson:"endpoint" json:"endpoint" `
  ClientIP  string                 `bson:"client_ip" json:"client_ip"`
  Username  *string                `bson:"username,omitempty" json:"username"`
  
  StatusCode int                   `bson:"status_code" json:"status_code"`
  Success bool                     `bson:"success" json:"success"`
  RequestTimeAt time.Time          `bson:"request_time_at" json:"request_time_at"`
  RequestPayload string            `bson:"request_payload" json:"request_payload"`
  ResponseTimeMS int               `bson:"response_time_ms" json:"response_time_ms"`
  Error     string                 `bson:"error_message, omitempty" json:"error_message"`
}


// {
//   "request_id": "01J9Y6M65PQD3ZC7Y5F4MZ8C9E",
//   "service": "shopee-partner-api",
//   "method": "POST",
//   "endpoint": "/api/v1/shopee/partner",
//   "request_time": "2025-08-29T12:34:56Z",
//   "request_payload": { "partner_id": "1234" },
//   "client_ip": "192.168.1.10",
//   "user_id": "u_5678",

//   "status_code": 200,
//   "success": true,
//   "response_time_ms": 152,
//   "error_message": null,

//   "trace_id": "abc123-trace",
//   "host": "pod-shopee-api-01"
// }
