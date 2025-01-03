//api/types.go
package api

import (
    "context"
    "net/http"
    "time"
)

// Method HTTP方法类型
type Method string

const (
    GET     Method = "GET"
    POST    Method = "POST"
    PUT     Method = "PUT"
    DELETE  Method = "DELETE"
    PATCH   Method = "PATCH"
)

// ContentType 内容类型
const (
    ContentTypeJSON       = "application/json"
    ContentTypeXML        = "application/xml"
    ContentTypeForm      = "application/x-www-form-urlencoded"
    ContentTypeMultipart = "multipart/form-data"
)

// Response API响应结构
type Response struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    RequestID string      `json:"request_id"`
    Timestamp time.Time   `json:"timestamp"`
}

// Error API错误
type Error struct {
    Code       int               `json:"code"`
    Message    string            `json:"message"`
    Details    map[string]string `json:"details,omitempty"`
    Stack      string           `json:"stack,omitempty"`
    InnerError error            `json:"-"`
}

func (e *Error) Error() string {
    return e.Message
}
