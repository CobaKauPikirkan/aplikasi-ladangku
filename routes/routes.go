package routes

import (
	"github.com/CobaKauPikirkan/aplikasi-ladangku/controller"
	"github.com/gin-gonic/gin"
)

func Routes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/user/signup", controller.SignUp())
	incomingRoutes.POST("/user/login", controller.Login())

	incomingRoutes.POST("/user/commodity/add", controller.AddCommodity())
	incomingRoutes.GET("/user/commodity/all", controller.GetCommodityAll())
	incomingRoutes.GET("/user/commodity/search", controller.SearchCommodityByQuery())
	incomingRoutes.GET("/user/commodity/byid", controller.CommodityById())

	incomingRoutes.POST("/user/ladang/add", controller.AddLadang())
	incomingRoutes.POST("/user/ladang/commodity", controller.AddCommodityToLadang())
	incomingRoutes.POST("/user/ladang/addtodolist", controller.AddTodoList())
	incomingRoutes.PUT("/user/ladang/deletecommodity", controller.DeleteCommodity())

	incomingRoutes.GET("/user/ladang/all", controller.GetLadangAll())
	incomingRoutes.GET("/user/ladang/byindex", controller.GetLadangByarray())
	incomingRoutes.GET("/user/ladang/gettodo", controller.GetAllTodo())
	incomingRoutes.GET("/user/ladang/getonetodo", controller.GetOneTodo())
	incomingRoutes.PUT("/user/ladang/updatetodo", controller.EditTodo())
	incomingRoutes.PUT("/user/ladang/deletetodo", controller.DeleteTodo())
}