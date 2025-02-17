package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Aliucord/Aliucord-backend/bot"
	"github.com/Aliucord/Aliucord-backend/common"
	"github.com/Aliucord/Aliucord-backend/database"
	"github.com/Aliucord/Aliucord-backend/updateTracker"
	"github.com/valyala/fasthttp"
)

func main() {
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var config *common.Config
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	database.InitDB(config.Database)
	if config.Bot.Enabled {
		bot.StartBot(config)
		defer bot.StopBot()
	}
	updateTracker.StartUpdateTracker(config)

	log.Println("Starting http server at port", config.Port)
	err = fasthttp.ListenAndServe(":"+config.Port, func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(
				ctx,
				"<html><head><title>Aliucord</title></head><body>High quality™ temp page.<br><a href=\"/links/github\">[GitHub Org]</a> <a href=\"/links/discord\">[Discord Server]</a></body></html>",
			)
		case "/links/github":
			ctx.Redirect("https://github.com/Aliucord", fasthttp.StatusMovedPermanently)
		case "/links/discord":
			ctx.Redirect("https://discord.gg/EsNDvBaHVU", fasthttp.StatusMovedPermanently)
		case "/download/discord":
			v := string(ctx.QueryArgs().Peek("v"))
			if v == "" {
				missingParams(ctx, []string{"v - discord version code"})
				return
			}
			version, err := strconv.Atoi(v)
			if err != nil {
				missingParams(ctx, []string{"v - discord version code"})
				return
			}
			url, err := updateTracker.GetDownloadURL(version, false)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				ctx.WriteString("apk not found: " + err.Error())
				return
			}
			ctx.Redirect(url, fasthttp.StatusFound)
		default:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.WriteString(fasthttp.StatusMessage(fasthttp.StatusNotFound))
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}

func missingParams(ctx *fasthttp.RequestCtx, params []string) {
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	ctx.WriteString("missing required query params:\n   " + strings.Join(params, "\n   "))
}
