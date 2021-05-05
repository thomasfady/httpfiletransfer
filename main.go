package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

var uploadDir = "uploads"

var adminAccounts gin.Accounts = gin.Accounts{}

func setupRouter(nostatic bool) *gin.Engine {
	r := gin.Default()

	r.NoRoute(noRouteHandler)

	fmt.Println("Enable GET file upload on /f/*filepath?c=BASE64_CONTENT")
	r.GET("/f/*filepath", getUploadHandler)
	fmt.Println("Enable templates on /tpl/*filepath")
	r.GET("/tpl/*filepath", templateHandler)
	fmt.Println("Enable sleeping on /sleep")
	r.GET("/sleep", sleep)

	r.PUT("/*filepath", putFileHandler)

	r.OPTIONS("/*filepath", corsHandler)

	if !nostatic {
		fmt.Println("Enable static on /s/")
		r.Static("s/", "./static")
	}

	return r
}

func setupAdminRouter(authorized *gin.RouterGroup) {

	authorized.GET("list_ips", adminListIps)
	authorized.GET("/list_files/*filepath", adminListFiles)
}

func main() {

	var ip string
	var nostatic bool
	var basic string

	flag.StringVar(&ip, "listen", "0.0.0.0:8080", "Ip and Port to listening on.")
	flag.StringVar(&basic, "basic", "", "Enable authorized part with basic auth")
	flag.BoolVar(&nostatic, "nostatic", false, "Disable static file serving.")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	r := setupRouter(nostatic)

	if basic != "" {
		adminAccounts = gin.Accounts{
			strings.Split(basic, ":")[0]: strings.Split(basic, ":")[1],
		}
		authorized := r.Group("/admin", gin.BasicAuth(adminAccounts))
		setupAdminRouter(authorized)
	}

	log.Printf("Listening on %s", ip)
	err := r.Run(ip)
	if err != nil {
		log.Fatal(err)
	}
}
