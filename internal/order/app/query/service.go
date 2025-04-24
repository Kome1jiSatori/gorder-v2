package query

import (
	"context"
	"github.com/Kome1jiSatori/gorder-v2/common/genproto/orderpb"
)

// 定义到stock的通信协议
type StockService interface {
	// 检查itemid在不在
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}
