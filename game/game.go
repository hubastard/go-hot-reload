package game

import rl "github.com/gen2brain/raylib-go/raylib"

type State struct {
	PlayerSpeed float32
	PlayerPos rl.Vector2
}

func (s *State) Update() {
	if rl.IsKeyDown(rl.KeyW) { s.PlayerPos.Y -= s.PlayerSpeed }
	if rl.IsKeyDown(rl.KeyS) { s.PlayerPos.Y += s.PlayerSpeed }
	if rl.IsKeyDown(rl.KeyA) { s.PlayerPos.X -= s.PlayerSpeed }
	if rl.IsKeyDown(rl.KeyD) { s.PlayerPos.X += s.PlayerSpeed }
}

func (s *State) Render() {
	rl.DrawCircleV(s.PlayerPos, 30, rl.Blue)
}