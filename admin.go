package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func adminListFiles(c *gin.Context) {
	baseDir := c.Query("base")
	if baseDir == "" {
		baseDir = uploadDir
	} else if baseDir != "static" && baseDir != "tpl" && baseDir != uploadDir {
		c.String(http.StatusNotFound, "")
		return
	}
	finalPath := filepath.Join(baseDir, filepath.Clean(c.Param("filepath")))
	if !strings.HasPrefix(filepath.Clean(finalPath), baseDir) {
		c.String(http.StatusForbidden, "")
		return
	}
	fi, err := os.Stat(finalPath)

	if err != nil {
		c.String(http.StatusNotFound, "")
		return
	}

	mode := fi.Mode()

	if mode.IsRegular() {
		c.Writer.Header()["Content-Type"] = []string{"application/octet-stream"}
		c.Writer.Header()["Content-Disposition"] = []string{"attachment; filename=" + fi.Name()}
		c.File(finalPath)
		return
	}

	files, err := ioutil.ReadDir(finalPath)
	var dirs_list []string
	var files_list []string
	if err == nil && len(files) != 0 {
		for _, f := range files {
			if f.IsDir() {
				dirs_list = append(dirs_list, strings.Replace(finalPath+"/"+f.Name(), baseDir, "", 1))
			} else {
				files_list = append(files_list, strings.Replace(finalPath+"/"+f.Name(), baseDir, "", 1))
			}
		}
	}
	renderHtmlTemplate(list_files_html, gin.H{
		"dirs_list":  dirs_list,
		"files_list": files_list,
		"filepath":   strings.Replace(finalPath, baseDir, "", 1),
		"baseDir":    baseDir,
	}, c)
}

var list_files_html string = `
<html>
<head>
  <title>HttpFileTransfer</title>
  <script src="/assets/app.js"></script>
</head>
<body>
  <h1>{{ .filepath }} in {{ .baseDir }} </h1>
  <h2><small>Show files in : <a href="/admin/list_files?base=static"> static<a/></small> - <a href="/admin/list_files?base=tpl"> templates<a/></small> - <a href="/admin/list_files"> uploads<a/></small></h2>
  <h2>Directories</h2>
  <ul>
  {{ range $element := .dirs_list }}
	<li><a href="/admin/list_files{{$element}}?base={{ $.baseDir }}">{{$element}}</a></li>
  {{end}}
  </ul>
  <h2>Files</h2>
  <ul>
  {{ range $element := .files_list }}
	<li><a href="/admin/list_files{{$element}}?base={{ $.baseDir }}">{{$element}}</a></li>
  {{end}}
  </ul>
</body>
</html>
`
