package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
)

var uploadDir = "uploads"

func saveFile(c *gin.Context, filePath string, content []byte, append bool) {
	baseDir := filepath.Join(uploadDir, c.ClientIP())
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

type TPLvar struct {
	Context *gin.Context
	Utils   TPLUtil
}

type TPLUtil struct {
	Context *gin.Context
}

func (u TPLUtil) Sleep() string {
	return "<img src=\"/f/sleep\"/>"
}

func (u TPLUtil) SleepTime(sleepTime string) string {
	return "<img src=\"/f/sleep?time=" + sleepTime + "\"/>"
}

func setupRouter(static bool, getUpload bool) *gin.Engine {
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})
	if getUpload {
		fmt.Println("Enable static to /s/")
		r.Static("s/", "./")
		fmt.Println("Enable sleeping to /sleep")
		r.GET("/sleep", func(c *gin.Context) {
			i, _ := strconv.Atoi(c.DefaultQuery("time", "1"))
			time.Sleep(time.Duration(i) * time.Second)
		})
		fmt.Println("Enable GET file upload to /f/*filepath?c=BASE64_CONTENT")
		r.GET("/f/*filepath", func(c *gin.Context) {
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
		})
	}
	r.GET("/tpl/*filepath", func(c *gin.Context) {
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
	})
	r.PUT("/*filepath", func(c *gin.Context) {
		filePath := c.Param("filepath")
		data, _ := c.GetRawData()
		saveFile(c, filePath, data, false)
	})
	r.OPTIONS("/*filepath", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "PUT, GET, POST, DELETE")
		c.String(http.StatusOK, "")
	})

	if static {
		fmt.Println("Enable static to /")
		r.Static("", "./")
	}

	return r
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

func renderHtmlTemplate(html string, data interface{}, c *gin.Context) {
	t, _ := template.New("tpl").Parse(html)
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		c.String(http.StatusOK, "OK")
	}
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, tpl.String())
}

func setupAdminRouter(authorized *gin.RouterGroup) {

	authorized.GET("list_ips", func(c *gin.Context) {
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

	})
	authorized.GET("/list_files/*filepath", func(c *gin.Context) {
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
	})
}

func main() {
	var ip string
	var static bool
	var getUpload bool
	var basic string
	flag.StringVar(&ip, "listen", "0.0.0.0:8080", "Ip and Port to listening on.")
	flag.StringVar(&basic, "basic", "", "Enable authorized part with basic auth")
	flag.BoolVar(&static, "static", true, "Enable static file serving.")
	flag.BoolVar(&getUpload, "getupload", false, "Enable GET upload (/*filename?c=BASE64). Move static serving to /s/.")
	flag.Parse()
	if getUpload {
		static = false
	}
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter(static, getUpload)
	if basic != "" {
		authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
			strings.Split(basic, ":")[0]: strings.Split(basic, ":")[1],
		}))
		setupAdminRouter(authorized)
	}
	log.Printf("Listening on %s", ip)
	err := r.Run(ip)
	if err != nil {
		log.Fatal(err)
	}
}
