//discovery/types.go
package discovery

import (
    "context"
    "time"
)

// ServiceInstance 服务实例
type ServiceInstance struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Version     string                 `json:"version"`
    Endpoint    string                 `json:"endpoint"`
    Status      ServiceStatus          `json:"status"`
    Metadata    map[string]interface{} `json:"metadata"`
    Tags        []string               `json:"tags"`
    Weight      int                    `json:"weight"`
    RegisterTime time.Time             `json:"register_time"`
    LastHeartbeat time.Time           `json:"last_heartbeat"`
}

// ServiceStatus 服务状态
type ServiceStatus string

const (
    StatusUp      ServiceStatus = "up"
    StatusDown    ServiceStatus = "down"
    StatusStarting ServiceStatus = "starting"
    StatusStopping ServiceStatus = "stopping"
    StatusUnknown  ServiceStatus = "unknown"
)

// ServiceEvent 服务事件
type ServiceEvent struct {
    Type      EventType
    Instance  *ServiceInstance
    Timestamp time.Time
}

type EventType string

const (
    EventRegister   EventType = "register"
    EventDeregister EventType = "deregister"
    EventModify     EventType = "modify"
    EventFailure    EventType = "failure"
)
