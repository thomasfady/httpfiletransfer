package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var uploadDir = "uploads"

func saveFile(c *gin.Context, filePath string, content []byte) {
	baseDir := filepath.Join(uploadDir, c.ClientIP())
	finalPath := filepath.Join(baseDir, filepath.Clean(filePath))
	if strings.HasPrefix(filepath.Clean(finalPath), baseDir) {
		if _, err := os.Stat(filepath.Dir(finalPath)); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(finalPath), 0700)
		}
		file, err := os.Create(finalPath)
		if err != nil {
			c.String(http.StatusForbidden, "")
			return
		}
		defer file.Close()
		file.Write(content)
		c.String(http.StatusOK, "OK")
	} else {
		c.String(http.StatusForbidden, "")
	}
}

func setupRouter(static bool, getUpload bool) *gin.Engine {
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})
	if getUpload {
		r.GET("/*filepath", func(c *gin.Context) {
			filePath := c.Param("filepath")
			if c.Query("c") != "" {
				content, _ := base64.StdEncoding.DecodeString(c.Query("c"))
				saveFile(c, filePath, content)

			} else {
				c.String(http.StatusOK, "OK")
			}
		})
	}
	r.PUT("/*filepath", func(c *gin.Context) {
		filePath := c.Param("filepath")
		data, _ := c.GetRawData()
		saveFile(c, filePath, data)
	})

	if static {
		fmt.Println("Enable static")
		r.Static("", "./")
	}

	return r
}

func main() {
	var ip string
	var static bool
	var getUpload bool
	flag.StringVar(&ip, "listen", "0.0.0.0:8080", "Ip and Port to listening on.")
	flag.BoolVar(&static, "static", true, "Enable static file serving.")
	flag.BoolVar(&getUpload, "getupload", false, "Enable GET upload (/*filename?c=BASE64). Disable static serving.")
	flag.Parse()
	if getUpload {
		static = false
	}
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter(static, getUpload)
	log.Printf("Listening on %s", ip)
	err := r.Run(ip)
	if err != nil {
		log.Fatal(err)
	}

}
