package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/hubastard/hot-reload/game"
)

const screenWidth int32 = 800
const screenHeight int32 = 450
const targetFPS int32 = 60

var gameState game.State

func main() {
	hotreload := flag.Bool("hot", false, "reload previous game state if true")
	hotProcessID := flag.Int("hotpid", 0, "process id of the hot reloader")
	flag.Parse()

	title := "Game"
	if *hotProcessID != 0 {
		title = "Game [HOT]"
	}

	rl.InitWindow(screenWidth, screenHeight, title)

	rl.SetTargetFPS(targetFPS)

	if *hotreload {
		bytes, err := os.ReadFile("hotstate.json")
		if err != nil {
			panic(err)
		}

		json.Unmarshal(bytes, &gameState)
	} else {
		gameState = game.State{
			PlayerSpeed: 3.0,
			PlayerPos: rl.Vector2{X: float32(screenWidth / 2), Y: float32(screenHeight / 2)},
		}
	}

	for !rl.WindowShouldClose() {
		gameState.Update()

		if *hotProcessID != 0 && rl.IsKeyPressed(rl.KeyF11) {
			hotReload()
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		gameState.Render()
		rl.EndDrawing()
	}

	rl.CloseWindow()

	if *hotProcessID != 0 {
		process, _ := os.FindProcess(*hotProcessID)
		if process != nil {
			process.Kill()
		}
	}
}

func hotReload() {
	file, err := json.MarshalIndent(gameState, "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("hotstate.json", file, 0644)
	if err != nil {
		panic("Error writing state file to disk")
	}
}