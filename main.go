package main

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"os"
	"strconv"
)

var total = 0
var cache map[string]string
var cache1 map[int]string

type H map[string]interface{}

func mainHandler(r render.Render) {
	r.Text(http.StatusOK, "URL Shortener")
}

func createHandler(r render.Render, params martini.Params) {
	url := params["url"]
	if govalidator.IsURL(url) {
		shortURL := ""
		if _, ok := cache[url]; ok {
			shortURL = cache[url]
		} else {
			shortURL = fmt.Sprintf("/%d", total)
			cache[url] = shortURL
			cache1[total] = url
			total++
		}
		r.JSON(http.StatusOK, H{"original_url": url, "short_url": shortURL})
	} else {
		r.JSON(http.StatusBadRequest, H{"error": "Wrong url format, make sure you have a valid protocol and real site."})
	}
}

func redirectHandler(r render.Render, params martini.Params) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		r.JSON(http.StatusBadRequest, H{"error": "This url is not on the database."})
		return
	}
	if _, ok := cache1[id]; ok {
		r.Redirect(cache1[id], http.StatusPermanentRedirect)
	} else {
		r.JSON(http.StatusBadRequest, H{"error": "This url is not on the database."})
	}
}

func main() {
	cache = make(map[string]string)
	cache1 = make(map[int]string)
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/", mainHandler)
	m.Get("/:id", redirectHandler)
	m.Get("/new/(?P<url>.*)", createHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	m.RunOnAddr(fmt.Sprintf(":%s", port))
}
