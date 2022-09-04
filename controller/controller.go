package controller

import (
	_ "draughts/docs"
	"draughts/service"
	"net/http"
	"os"
	// "strconv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const address = "localhost"
const port = "8080"
const boardCookie = "DRAUGHTSTATUS"

func setBoardStringCookie(c *gin.Context, cookieValue string) {
	c.SetCookie(boardCookie, cookieValue, 3600*24*365, "/", address, false, true)
}

func getStatusCookie(c *gin.Context) string {
	statusCookie, _ := c.Cookie(boardCookie)
	return statusCookie
}

// @Router       /api/v1/board [post]
func getNewBoard(c *gin.Context) {
	result, statusString := service.GetNewBoard()
	setBoardStringCookie(c, statusString)
	c.IndentedJSON(http.StatusOK, result)
}

// @Router       /api/v1/board [get]
func getBoard(c *gin.Context) {
	cookie := getStatusCookie(c)
	result, statusString := service.GetBoard(cookie)
	setBoardStringCookie(c, statusString)
	c.IndentedJSON(http.StatusOK, result)
}

// // @Router       /api/v1/move/{column}/{row} [post]
// func doMove(c *gin.Context) {
// 	col, _ := strconv.Atoi(c.Param("column"))
// 	row, _ := strconv.Atoi(c.Param("row"))
// 	cookie := getStatusCookie(c)
// 	result, statusString := service.DoMove(cookie, col, row)
// 	setBoardStringCookie(c, statusString)
// 	c.IndentedJSON(http.StatusOK, result)
// }

// // @Router       /api/v1/move/passmove [post]
// func doPassMove(c *gin.Context) {
// 	cookie := getStatusCookie(c)
// 	result, statusString := service.DoPassMove(cookie)
// 	setBoardStringCookie(c, statusString)
// 	c.IndentedJSON(http.StatusOK, result)
// }

// // @Router       /api/v1/move/takeback/ [post]
// func takeBackLastMove(c *gin.Context) {
// 	cookie := getStatusCookie(c)
// 	result, statusString := service.TakeBackLastMove(cookie)
// 	setBoardStringCookie(c, statusString)
// 	c.IndentedJSON(http.StatusOK, result)
// }

// @Router       / [get]
func getHtml(c *gin.Context) {
	html, err := os.ReadFile("./view/draughts10x10.html")
	if err != nil {
		panic(err)
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", html)
}

//----------------------------------------------------------------------------------------------------------------------

func getRouter() *gin.Engine {
	var router = gin.New()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}
	return router
}

func setHandlers(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", getHtml)

	router.GET("/api/v1/board", getBoard)
	router.POST("/api/v1/board", getNewBoard)
	// router.POST("/api/v1/move/:column/:row/", doMove)
	// router.POST("/api/v1/move/passmove/", doPassMove)
	// router.POST("/api/v1/move/takeback/", takeBackLastMove)
}

func startRouter(router *gin.Engine) {

	err := router.Run(address + ":" + port)
	if err != nil {
		panic(err)
	}
}

func RunController() {
	var router = getRouter()
	setHandlers(router)
	startRouter(router)
}
