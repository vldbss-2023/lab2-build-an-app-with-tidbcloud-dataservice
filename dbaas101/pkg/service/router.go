package service

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine, api *API) {
	r := engine.Group("/api/v1/")
	r.GET("/tidbclusters/:namespace/:name", api.GetTidbClusterHandler)
	r.DELETE("/tidbclusters/:namespace/:name", api.DeleteTidbClusterHandler)
	r.POST("/tidbclusters", api.CreateTidbClusterHandler)
	r.GET("/tidbclusters", api.ListTidbClusterHandler)
}
