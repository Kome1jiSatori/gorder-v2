package main

import (
	"github.com/Kome1jiSatori/gorder-v2/order/app"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	app app.Application
}

func (s HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {

}

func (s HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {

}
