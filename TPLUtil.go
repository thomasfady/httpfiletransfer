package main

import "github.com/gin-gonic/gin"

type TPLUtil struct {
	Context *gin.Context
}

func (u TPLUtil) Sleep() string {
	return "<img src=\"/f/sleep\"/>"
}

func (u TPLUtil) SleepTime(sleepTime string) string {
	return "<img src=\"/f/sleep?time=" + sleepTime + "\"/>"
}
