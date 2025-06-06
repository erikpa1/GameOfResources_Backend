package api

import (
	"github.com/gin-gonic/gin"
	"turtle/auth"
	"turtle/ctrl"
	"turtle/models"
	"turtle/tools"
)

func _ListProjects(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, ctrl.ListProjects(user.Org))
}

func _COUProject(c *gin.Context) {
	user := auth.GetUserFromContext(c)

	project := tools.ObjFromJsonPtr[models.TurtleProject](c.PostForm("data"))

	ctrl.COUProject(user, project)
	tools.AutoReturn(c, project.Uid)

}

func initApiProjects(r *gin.Engine) {
	r.GET("/api/projects", auth.LoginOrApp, _ListProjects)
	r.POST("/api/project", auth.LoginRequired, _COUProject)
}
