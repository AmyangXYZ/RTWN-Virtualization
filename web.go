package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/AmyangXYZ/sgo"
	"github.com/AmyangXYZ/sgo/middlewares"
	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	WSOnline = false
)

func RunWeb() {
	web := sgo.New()
	web.USE(middlewares.CORS(middlewares.CORSOpt{}))
	web.GET("/ws", ws)

	web.GET("/api/networks", getAllNetworks)
	web.GET("/api/network/:id", getNetwork)

	web.PUT("/api/new/topo/:seed", putNewTopo)
	web.OPTIONS("/api/new/topo/:seed", sgo.PreflightHandler)
	web.PUT("/api/new/app/:seed", putNewApps)
	web.OPTIONS("/api/new/app/:seed", sgo.PreflightHandler)
	if err := web.Run(":8000"); err != nil {
		panic(err)
	}
}

func getAllNetworks(ctx *sgo.Context) error {
	ids := []int{}
	for _, t := range Networks {
		ids = append(ids, t.ID)
	}
	sort.Ints(ids)
	return ctx.JSON(200, 1, "success", ids)
}

func getNetwork(ctx *sgo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return err
	}
	return ctx.JSON(200, 1, "success", Networks[id])
}

func ws(ctx *sgo.Context) error {
	// ws, err := upgrader.Upgrade(ctx.Resp, ctx.Req, nil)
	// breakSig := make(chan bool)
	// if err != nil {
	// 	return err
	// }
	// WSOnline = true
	// // fmt.Println("ws/comm connected")
	// defer func() {
	// 	ws.Close()
	// 	WSOnline = false
	// 	// fmt.Println("ws/comm client closed")
	// }()
	// go func() {
	// 	for {
	// 		_, _, err := ws.ReadMessage()
	// 		if err != nil {
	// 			breakSig <- true
	// 		}
	// 	}
	// }()
	// for {
	// 	select {
	// 	case l := <-NetworkRes:
	// 		ws.WriteJSON(l)
	// 	case <-breakSig:
	// 		WSOnline = false
	// 		return nil
	// 	}
	// }
	return nil
}

func putNewTopo(ctx *sgo.Context) error {
	seed, err := strconv.Atoi(ctx.Param("seed"))
	if err != nil {
		return err
	}
	Networks[0].Settings.RAND_SEED_TOPO = int64(seed)
	Networks[0].Run()
	return ctx.JSON(200, 1, "success", 1)
}

func putNewApps(ctx *sgo.Context) error {
	seed, err := strconv.Atoi(ctx.Param("seed"))
	if err != nil {
		return err
	}
	Networks[0].Settings.RAND_SEED_APP = int64(seed)
	Networks[0].Run()
	return ctx.JSON(200, 1, "success", 1)
}
