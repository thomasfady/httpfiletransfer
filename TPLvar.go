package main

import "github.com/gin-gonic/gin"

type TPLvar struct {
	Context *gin.Context
	Utils   TPLUtil
}
