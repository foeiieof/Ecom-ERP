package response

type APIResponse[T any] struct {
    Success         bool        `json:"success"`
    Message         string      `json:"message"`
    RequestID       string      `json:"request_id"`
    TimestampUnix   int64       `json:"timestamp_unix"`
    TimestampUTC    string      `json:"timestamp_utc"`
    TimestampLocal  string      `json:"timestamp_local"`
    Data            T           `json:"data,omitempty"`
    Error           any         `json:"error,omitempty"`
}


