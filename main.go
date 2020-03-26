package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var uploadDir = "uploads"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})
	r.PUT("/*filepath", func(c *gin.Context) {
		filePath := c.Param("filepath")
		baseDir := filepath.Join(uploadDir, c.ClientIP())
		finalPath := filepath.Join(baseDir, filepath.Clean(filePath))
		if strings.HasPrefix(filepath.Clean(finalPath), baseDir) {
			if _, err := os.Stat(filepath.Dir(finalPath)); os.IsNotExist(err) {
				os.MkdirAll(filepath.Dir(finalPath), 0700)
			}
			file, err := os.Create(finalPath)
			if err != nil {
				return
			}
			defer file.Close()
			data, _ := c.GetRawData()
			file.Write(data)
			c.String(http.StatusOK, "OK")
		} else {
			c.String(http.StatusForbidden, "")
		}

	})
	r.Static("", "./")
	return r
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	log.Print("Listening on :8080")
	r.Run(":8080")
}
