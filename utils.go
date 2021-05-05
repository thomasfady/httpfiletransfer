package main

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

func saveFile(c *gin.Context, filePath string, content []byte, append bool) {
	baseDir := filepath.Join(uploadDir, c.ClientIP())
	saveFileTo(c, filePath, content, append, baseDir)
}

func saveFileTo(c *gin.Context, filePath string, content []byte, append bool, baseDir string) {

	finalPath := filepath.Join(baseDir, filepath.Clean(filePath))
	if strings.HasPrefix(filepath.Clean(finalPath), baseDir) {
		if _, err := os.Stat(filepath.Dir(finalPath)); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(finalPath), 0700)
		}
		var file *os.File
		var err error

		if append {
			file, err = os.OpenFile(finalPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			file, err = os.Create(finalPath)
		}

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

func renderHtmlTemplate(html string, data interface{}, c *gin.Context) {
	t, _ := template.New("tpl").Parse(html)
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		c.String(http.StatusOK, "OK")
	}
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, tpl.String())
}

func sleep(c *gin.Context) {
	i, _ := strconv.Atoi(c.DefaultQuery("time", "1"))
	time.Sleep(time.Duration(i) * time.Second)
}

func checkBasicAuth(c *gin.Context) (string, bool) {
	for user, password := range adminAccounts {
		if "Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+password)) == c.GetHeader("Authorization") {
			return user, true
		}
	}
	return "", false
}
