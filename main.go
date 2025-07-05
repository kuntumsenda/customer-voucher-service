package main

import (
	"customer-voucher-service/routes"
	"github.com/gin-gonic/gin"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	r := gin.Default()

	routes.ApiRoutes(r)

	r.Run(":8080")
}
