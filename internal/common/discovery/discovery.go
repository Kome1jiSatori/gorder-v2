package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	// 服务唯一键、服务名、ip地址
	Register(ctx context.Context, instanceID, serviceName, hostPort string) error
	Deregister(ctx context.Context, instanceID, serviceName string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	// 探测是否存活
	HealthCheck(instanceID, serviceName string) error
}

func GenerateInstanceID(serviceName string) string {
	x := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s-%d", serviceName, x)
}
