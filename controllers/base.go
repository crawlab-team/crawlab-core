package controllers

import "github.com/gin-gonic/gin"

type Controller interface {
	Get(c *gin.Context)
	Post(c *gin.Context)
	Put(c *gin.Context)
	Delete(c *gin.Context)
}

type ListController interface {
	Controller
	GetList(c *gin.Context)
	PutList(c *gin.Context)
	PostList(c *gin.Context)
	DeleteList(c *gin.Context)
}

type PostAction struct {
	Name        string
	HandlerFunc gin.HandlerFunc
}

type PostActionController interface {
	Actions() (actions []PostAction)
}
