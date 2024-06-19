package main

import (
	"context"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sort"
	_ "star_dust/config"
	"star_dust/consts"
	"star_dust/inter"
	_ "star_dust/things"
)

type Game struct {
}

func (g *Game) Update() error {
	select {
	case <-consts.ClosingSignal.Done():
		return context.Canceled
	default:
	}
	if ebiten.IsWindowBeingClosed() {
		consts.ClosingCommand()
		return nil
	}
	inter.GameThings.Range(func(key int64, value inter.Thing) bool {
		td := value
		td.Update()
		return true
	})
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	orderedThingsKeys := make([]inter.Thing, 0)
	inter.GameThings.Range(func(key int64, value inter.Thing) bool {
		orderedThingsKeys = append(orderedThingsKeys, value)
		return true
	})
	sort.Slice(orderedThingsKeys, func(i, j int) bool {
		return orderedThingsKeys[i].GetLayers() < orderedThingsKeys[j].GetLayers()
	})
	for _, thing := range orderedThingsKeys {
		thing.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return inter.ScreenWidth, inter.ScreenHeight
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	g := new(Game)

	ebiten.SetWindowSize(inter.ScreenWidth, inter.ScreenHeight)
	ebiten.SetWindowTitle("Star Dust")
	ebiten.SetWindowClosingHandled(true)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
