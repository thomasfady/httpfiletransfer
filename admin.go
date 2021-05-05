package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func adminListIps(c *gin.Context) {
	files, err := ioutil.ReadDir(uploadDir)
	var ips []string
	if err == nil && len(files) != 0 {
		for _, f := range files {
			if f.IsDir() {
				ips = append(ips, f.Name())
			}
		}
	}
	renderHtmlTemplate(html, gin.H{
		"ips": ips,
	}, c)

}

func adminListFiles(c *gin.Context) {
	filepath := c.Param("filepath")
	fullpath := uploadDir + "/" + filepath
	fi, err := os.Stat(fullpath)

	if err != nil {
		c.String(http.StatusNotFound, "")
		return
	}

	mode := fi.Mode()

	if mode.IsRegular() {
		c.Writer.Header()["Content-Type"] = []string{"application/octet-stream"}
		c.Writer.Header()["Content-Disposition"] = []string{"attachment; filename=" + fi.Name()}
		c.File(fullpath)
		return
	}

	files, err := ioutil.ReadDir(fullpath)
	var dirs_list []string
	var files_list []string
	if err == nil && len(files) != 0 {
		for _, f := range files {
			if f.IsDir() {
				dirs_list = append(dirs_list, filepath+"/"+f.Name())
			} else {
				files_list = append(files_list, filepath+"/"+f.Name())
			}
		}
	}
	renderHtmlTemplate(list_files_html, gin.H{
		"dirs_list":  dirs_list,
		"files_list": files_list,
		"filepath":   filepath,
	}, c)
}

var html string = `
<html>
<head>
  <title>Https Test</title>
  <script src="/assets/app.js"></script>
</head>
<body>
  <h1>IPs</h1>
  <ul>
  {{ range $element := .ips }}
	<li><a href="/admin/list_files/{{$element}}">{{$element}}</a></li>
	{{end}}
</ul>
</body>
</html>
`

var list_files_html string = `
<html>
<head>
  <title>Https Test</title>
  <script src="/assets/app.js"></script>
</head>
<body>
  <h1>{{ .filepath }}</h1>
  <h2>Directories</h2>
  <ul>
  {{ range $element := .dirs_list }}
	<li><a href="/admin/list_files{{$element}}">{{$element}}</a></li>
  {{end}}
  </ul>
  <h2>Files</h2>
  <ul>
  {{ range $element := .files_list }}
	<li><a href="/admin/list_files{{$element}}">{{$element}}</a></li>
  {{end}}
  </ul>
</body>
</html>
`
