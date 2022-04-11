package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//todo 需要完善任务状态等， 把运行的任务管理起来
type Server struct {
	*http.Server
	Running     bool // 标记有任务在运行
	*gin.Engine      //router
}

//todo
func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() error {
	return s.ListenAndServe()
}
