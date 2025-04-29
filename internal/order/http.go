package main

import (
	"fmt"
	"github.com/Kome1jiSatori/gorder-v2/common"
	"github.com/Kome1jiSatori/gorder-v2/order/app/dto"
	"github.com/Kome1jiSatori/gorder-v2/order/convertor"
	"net/http"

	client "github.com/Kome1jiSatori/gorder-v2/common/client/order"
	"github.com/Kome1jiSatori/gorder-v2/order/app"
	"github.com/Kome1jiSatori/gorder-v2/order/app/command"
	"github.com/Kome1jiSatori/gorder-v2/order/app/query"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	common.BaseResponse
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var (
		req  client.CreateOrderRequest
		err  error
		resp dto.CreateOrderResponse
	)
	defer func() {
		H.Response(c, err, resp)
	}()
	if err := c.ShouldBindJSON(&req); err != nil {
		H.Response(c, err, nil)
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		// c.JSON(http.StatusOK, gin.H{"error": err})
		H.Response(c, err, nil)
		return
	}
	resp.CustomerID = req.CustomerID
	resp.RedirectURL = fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID)
	resp.OrderID = r.OrderID
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	var (
		err  error
		resp interface{}
	)
	defer func() {
		H.Response(c, err, resp)
	}()
	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		CustomerID: customerID,
		OrderID:    orderID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	resp = convertor.NewOrderConvertor().EntityToClient(o)
}
