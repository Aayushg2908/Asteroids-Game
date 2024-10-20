package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	THICKNESS  = 2.0
	SCALE      = 30.0
	ROT_SPEED  = 1.3
	SHIP_SPEED = 25
	DRAG       = 0.3
)

var SIZE = rl.NewVector2(640*1.2, 480*1.2)

type Ship struct {
	pos rl.Vector2
	vel rl.Vector2
	rot float32
}

type State struct {
	now   float32
	delta float32
	ship  Ship
}

var state = State{
	now:   0.0,
	delta: 0.0,
	ship: Ship{
		pos: rl.Vector2Scale(SIZE, 0.5),
		vel: rl.NewVector2(0, 0),
		rot: 0.0,
	},
}

type Transformer struct {
	org   rl.Vector2
	scale float32
	rot   float32
}

func (t *Transformer) apply(point rl.Vector2) rl.Vector2 {
	return rl.Vector2Add(rl.Vector2Scale(rl.Vector2Rotate(point, t.rot), t.scale), t.org)
}

func drawLines(org rl.Vector2, scale float32, rot float32, points []rl.Vector2) {
	transformer := Transformer{org, scale, rot}

	for i := 0; i < len(points); i++ {
		rl.DrawLineEx(transformer.apply(points[i]), transformer.apply(points[(i+1)%len(points)]), THICKNESS, rl.White)
	}
}

func update() {
	if rl.IsKeyDown(rl.KeyRight) {
		state.ship.rot += state.delta * ROT_SPEED * math.Pi * 2
	} else if rl.IsKeyDown(rl.KeyLeft) {
		state.ship.rot -= state.delta * ROT_SPEED * math.Pi * 2
	}

	dirAngle := state.ship.rot + math.Pi*0.5
	shipDir := rl.NewVector2(float32(math.Cos(float64(dirAngle))), float32(math.Sin(float64(dirAngle))))

	if rl.IsKeyDown(rl.KeyUp) {
		state.ship.vel = rl.Vector2Add(state.ship.vel, rl.Vector2Scale(shipDir, SHIP_SPEED*state.delta))
	}

	state.ship.vel = rl.Vector2Scale(state.ship.vel, 1-DRAG*state.delta)

	if state.ship.pos.X < 0 {
		state.ship.pos.X = SIZE.X
	} else if state.ship.pos.X > SIZE.X {
		state.ship.pos.X = 0
	} else if state.ship.pos.Y < 0 {
		state.ship.pos.Y = SIZE.Y
	} else if state.ship.pos.Y > SIZE.Y {
		state.ship.pos.Y = 0
	}
	state.ship.pos = rl.Vector2Add(state.ship.pos, state.ship.vel)
}

func render() {
	drawLines(state.ship.pos, SCALE, state.ship.rot, []rl.Vector2{
		rl.NewVector2(-0.4, -0.5),
		rl.NewVector2(0.0, 0.5),
		rl.NewVector2(0.4, -0.5),
		rl.NewVector2(0.3, -0.4),
		rl.NewVector2(-0.3, -0.4),
	})

	// make the thruster blink
	if int(state.now*20)%2 == 0 && rl.IsKeyDown(rl.KeyUp) {
		drawLines(state.ship.pos, SCALE, state.ship.rot, []rl.Vector2{
			rl.NewVector2(-0.3, -0.4),
			rl.NewVector2(0.0, -0.73),
			rl.NewVector2(0.3, -0.4),
		})
	}
}

func main() {
	rl.InitWindow(int32(SIZE.X), int32(SIZE.Y), "Asteroids Game")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		state.delta = rl.GetFrameTime()
		state.now += state.delta

		update()

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		render()

		rl.EndDrawing()
	}
}
