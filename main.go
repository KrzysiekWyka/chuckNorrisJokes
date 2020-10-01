package main

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

type Joke struct {
	IconUrl string `json:"icon_url"`
	Value   string `json:"value"`
}

func main() {
	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: cfg,
	}

	router := gin.Default()

	router.LoadHTMLGlob("views/*")

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/", func(context *gin.Context) {
		resp, _ := http.Get("https://api.chucknorris.io/jokes/random")

		body, _ := ioutil.ReadAll(resp.Body)

		var joke Joke

		if err := json.Unmarshal(body, &joke); err != nil {
			context.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": err.Error()})

			return
		}

		var hostname string

		hostnameQueryParam := context.Query("displayHostname")

		if hostnameQueryParam == "true" {
			hostnameEnv, _ := os.Hostname()

			hostname = hostnameEnv
		}

		context.HTML(http.StatusOK, "index.tmpl", gin.H{"joke": joke, "hostname": hostname})
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	router.Run(":" + port)
}
