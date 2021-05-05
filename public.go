package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
)

func templateHandler(c *gin.Context) {
	filePath := c.Param("filepath")
	dat, err := ioutil.ReadFile("./tpl/" + filePath)
	u := TPLUtil{c}
	tpl_var := TPLvar{c, u}
	t, _ := template.New("tpl").Parse(string(dat))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, tpl_var); err != nil {
		c.String(http.StatusOK, "OK")
	}
	mime := mimetype.Detect(tpl.Bytes())
	if err != nil {
		c.String(http.StatusOK, "OK")
	}
	c.Header("Content-Type", mime.String())
	c.String(http.StatusOK, tpl.String())
}

func getUploadHandler(c *gin.Context) {
	filePath := c.Param("filepath")
	if c.Query("c") != "" {
		content, _ := base64.StdEncoding.DecodeString(c.Query("c"))
		saveFile(c, filePath, content, false)
	} else if c.Query("ca") != "" {
		content, _ := base64.StdEncoding.DecodeString(c.Query("ca"))
		saveFile(c, filePath, content, true)
	} else {
		c.String(http.StatusOK, "OK")
	}
}

func putFileHandler(c *gin.Context) {
	filePath := c.Param("filepath")
	if strings.HasPrefix(filePath, "/admin/") {
		_, auth := checkBasicAuth(c)
		if !auth {
			c.String(http.StatusUnauthorized, "")
		} else {
			var dir = strings.Split(filePath, "/")[2]
			if dir != "static" && dir != "tpl" {
				c.String(http.StatusUnauthorized, "")
				return
			}
			data, _ := c.GetRawData()
			saveFileTo(c, strings.Replace(c.Param("filepath"), "/admin/"+dir, "", 1), data, false, dir)
		}
		return
	}
	data, _ := c.GetRawData()
	saveFile(c, filePath, data, false)
}

func noRouteHandler(c *gin.Context) {
	c.String(http.StatusOK, "")
}

func corsHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "PUT, GET, POST, DELETE")
	c.String(http.StatusOK, "")
}
