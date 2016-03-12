package main

import (
	"flag"
	"log"
	// "net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

func main() {
	flag.Parse()

	router := fasthttprouter.New()
	router.GET("/css/*path", cssRequestHandler)

	if err := fasthttp.ListenAndServe(":8080", router.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func cssRequestHandler(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	// FIXME
	minifier := minify.New()
	minifier.AddFunc("style/css", css.Minify)

	path := ps.ByName("path")
	path = path[1:]
	// recovering path.
	if path[0:5] == "https" {
		path = "https://" + path[7:]
	}
	if path[0:4] == "http" {
		path = "http://" + path[6:]
	}
	// resp, err := http.Get(path)
	// if err != nil {
	// 	ctx.Error("Invalid request", 400)
	// 	ctx.Close()
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	_, body, err := fasthttp.Get(nil, path)
	if err != nil {
		log.Printf("%s\n", err)
		ctx.Error("Invalid request", 400)
	}
	res, err := minifier.Bytes("style/css", body)
	if err != nil {
		log.Printf("%s\n", err)
	}
	ctx.Write(res)
}
