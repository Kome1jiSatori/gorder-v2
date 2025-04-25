package domain

import (
	"context"
	"github.com/Kome1jiSatori/gorder-v2/common/genproto/orderpb"
)

// 创建支付连接
type Processor interface {
	CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
}
