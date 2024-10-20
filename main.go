package main

import (
	"math"
	"math/rand"

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
	pos       rl.Vector2
	vel       rl.Vector2
	rot       float32
	deathTime float32
}

func (s Ship) isDead() bool {
	return s.deathTime != 0.0
}

type Asteroid struct {
	pos  rl.Vector2
	vel  rl.Vector2
	size AsteroidSize
	seed int64
}

type State struct {
	now       float32
	delta     float32
	ship      Ship
	asteroids []Asteroid
	particles []Particle
}

type ParticleType int

const (
	LINE ParticleType = iota
	DOT
)

type Particle struct {
	pos rl.Vector2
	vel rl.Vector2
	ttl float32

	pType  ParticleType
	rot    float32
	len    float32
	radius float32
}

var state = State{
	now:   0.0,
	delta: 0.0,
	ship: Ship{
		pos:       rl.Vector2Scale(SIZE, 0.5),
		vel:       rl.NewVector2(0, 0),
		rot:       0.0,
		deathTime: 0.0,
	},
	asteroids: []Asteroid{},
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

type AsteroidSize int

const (
	SMALL AsteroidSize = iota
	MEDIUM
	BIG
)

func (s AsteroidSize) size() float32 {
	switch s {
	case SMALL:
		return SCALE * 0.9
	case MEDIUM:
		return SCALE * 1.5
	case BIG:
		return SCALE * 3.0
	}
	return 0.0
}

func (s AsteroidSize) collisionScale() float32 {
	switch s {
	case SMALL:
		return 1.0
	case MEDIUM:
		return 0.8
	case BIG:
		return 0.5
	}
	return 0.0
}

func (s AsteroidSize) velocity() float32 {
	switch s {
	case SMALL:
		return 1.8
	case MEDIUM:
		return 1.4
	case BIG:
		return 0.8
	}
	return 0.0
}

func drawAsteroid(pos rl.Vector2, s AsteroidSize, seed int64) {
	random := rand.New(rand.NewSource(seed))
	var points []rl.Vector2
	numPoints := 8 + random.Intn(8)

	for i := 0; i < numPoints; i++ {
		radius := 0.3 + (0.2 * random.Float64())
		if random.Float64() < 0.2 {
			radius -= 0.2
		}
		angle := ((math.Pi * 2 / float64(numPoints)) * float64(i)) + (math.Pi * 0.125 * random.Float64())
		points = append(points, rl.Vector2Scale(rl.NewVector2(float32(math.Cos(angle)), float32(math.Sin(angle))), float32(radius)))
	}

	drawLines(pos, s.size(), 0, points)
}

func update() {
	if !state.ship.isDead() {
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

	for i := 0; i < len(state.asteroids); i++ {
		asteroid := &state.asteroids[i]
		asteroid.pos = rl.Vector2Add(asteroid.pos, asteroid.vel)

		if asteroid.pos.X < 0 {
			asteroid.pos.X = SIZE.X
		} else if asteroid.pos.X > SIZE.X {
			asteroid.pos.X = 0
		} else if asteroid.pos.Y < 0 {
			asteroid.pos.Y = SIZE.Y
		} else if asteroid.pos.Y > SIZE.Y {
			asteroid.pos.Y = 0
		}

		if !state.ship.isDead() && rl.Vector2Distance(asteroid.pos, state.ship.pos) < asteroid.size.size()*asteroid.size.collisionScale() {
			state.ship.deathTime = state.now

			for i := 0; i < 5; i++ {
				angle := 2 * math.Pi * rand.Float32()
				state.particles = append(state.particles, Particle{
					pos:    rl.Vector2Add(state.ship.pos, rl.NewVector2(rand.Float32()*3, rand.Float32()*3)),
					vel:    rl.Vector2Scale(rl.NewVector2(float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))), 2*rand.Float32()),
					ttl:    3 + rand.Float32(),
					pType:  LINE,
					rot:    angle,
					len:    SCALE * (0.6 + (0.4 * rand.Float32())),
					radius: 0,
				})
			}
		}
	}

	for i := 0; i < len(state.particles); i++ {
		particle := &state.particles[i]
		particle.pos = rl.Vector2Add(particle.pos, particle.vel)

		if particle.ttl > state.delta {
			particle.ttl -= state.delta
		} else {
			state.particles = append(state.particles[:i], state.particles[i+1:]...)
			i--
		}
	}

	if state.ship.isDead() && state.now-state.ship.deathTime > 3.0 {
		resetGame()
	}
}

func render() {
	if !state.ship.isDead() {
		drawLines(state.ship.pos, SCALE, state.ship.rot, []rl.Vector2{
			rl.NewVector2(-0.4, -0.5),
			rl.NewVector2(0.0, 0.5),
			rl.NewVector2(0.4, -0.5),
			rl.NewVector2(0.3, -0.4),
			rl.NewVector2(-0.3, -0.4),
		})

		if int(state.now*20)%2 == 0 && rl.IsKeyDown(rl.KeyUp) {
			drawLines(state.ship.pos, SCALE, state.ship.rot, []rl.Vector2{
				rl.NewVector2(-0.3, -0.4),
				rl.NewVector2(0.0, -0.73),
				rl.NewVector2(0.3, -0.4),
			})
		}
	}

	for _, asteroid := range state.asteroids {
		drawAsteroid(asteroid.pos, asteroid.size, asteroid.seed)
	}

	for _, particle := range state.particles {
		switch particle.pType {
		case LINE:
			drawLines(particle.pos, particle.len, particle.rot, []rl.Vector2{
				rl.NewVector2(-0.5, 0),
				rl.NewVector2(0.5, 0),
			})
		case DOT:
			rl.DrawCircleV(particle.pos, particle.radius, rl.White)
		}
	}
}

func resetAsteroids() {
	state.asteroids = []Asteroid{}

	for i := 0; i < 20; i++ {
		angle := 2 * math.Pi * rand.Float64()
		size := AsteroidSize(rand.Intn(3))
		state.asteroids = append(state.asteroids, Asteroid{
			pos:  rl.NewVector2(rand.Float32()*SIZE.X, rand.Float32()*SIZE.Y),
			vel:  rl.Vector2Scale(rl.NewVector2(float32(math.Cos(angle)), float32(math.Sin(angle))), size.velocity()*3*rand.Float32()),
			size: size,
			seed: rand.Int63(),
		})
	}
}

func resetGame() {
	state.ship.deathTime = 0.0
	state.ship = Ship{
		pos: rl.Vector2Scale(SIZE, 0.5),
		vel: rl.NewVector2(0, 0),
		rot: 0.0,
	}
}

func main() {
	rl.InitWindow(int32(SIZE.X), int32(SIZE.Y), "Asteroids Game")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	resetGame()
	resetAsteroids()

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
