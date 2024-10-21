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

var SIZE = rl.NewVector2(640*1.5, 480*1.5)

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
	pos    rl.Vector2
	vel    rl.Vector2
	size   AsteroidSize
	seed   int64
	remove bool
}

type AlienSize int

const (
	TINY AlienSize = iota
	HUGE
)

type Alien struct {
	pos      rl.Vector2
	dir      rl.Vector2
	size     AlienSize
	remove   bool
	lastShot float32
	lastDir  float32
}

type State struct {
	now       float32
	delta     float32
	ship      Ship
	asteroids []Asteroid
	particles []Particle
	bullets   []Bullet
	aliens    []Alien
	lives     int
	score     int
	reset     bool
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

type Bullet struct {
	pos    rl.Vector2
	vel    rl.Vector2
	ttl    float32
	remove bool
	spawn  float32
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
	particles: []Particle{},
	bullets:   []Bullet{},
	lives:     3,
	score:     0,
	reset:     false,
}

type Transformer struct {
	org   rl.Vector2
	scale float32
	rot   float32
}

func (t *Transformer) apply(point rl.Vector2) rl.Vector2 {
	return rl.Vector2Add(rl.Vector2Scale(rl.Vector2Rotate(point, t.rot), t.scale), t.org)
}

var NUMBERS = [][][]rl.Vector2{
	{
		{rl.NewVector2(-0.5, -0.5), rl.NewVector2(0.5, -0.5)},
		{rl.NewVector2(0.5, -0.5), rl.NewVector2(0.5, 0.5)},
		{rl.NewVector2(0.5, 0.5), rl.NewVector2(-0.5, 0.5)},
		{rl.NewVector2(-0.5, 0.5), rl.NewVector2(-0.5, -0.5)},
	},
	{
		{rl.NewVector2(0.5, -0.5), rl.NewVector2(0.5, 0.5)},
	},
	{
		{rl.NewVector2(-0.5, -0.5), rl.NewVector2(0.5, -0.5)},
		{rl.NewVector2(0.5, -0.5), rl.NewVector2(0.5, 0)},
		{rl.NewVector2(0.5, 0), rl.NewVector2(-0.5, 0)},
		{rl.NewVector2(-0.5, 0), rl.NewVector2(-0.5, 0.5)},
		{rl.NewVector2(-0.5, 0.5), rl.NewVector2(0.5, 0.5)},
	},
	{
		{rl.NewVector2(-0.5, -0.5), rl.NewVector2(0.5, -0.5)},
		{rl.NewVector2(0.5, -0.5), rl.NewVector2(0.5, 0.5)},
		{rl.NewVector2(0.5, 0.5), rl.NewVector2(-0.5, 0.5)},
		{rl.NewVector2(-0.5, 0), rl.NewVector2(0.5, 0)},
	},
	{
		{rl.NewVector2(-0.5, -0.5), rl.NewVector2(-0.5, 0)},
		{rl.NewVector2(-0.5, 0), rl.NewVector2(0.5, 0)},
		{rl.NewVector2(0.5, -0.5), rl.NewVector2(0.5, 0.5)},
	},
	{
		{rl.NewVector2(0.5, -0.5), rl.NewVector2(-0.5, -0.5)},
		{rl.NewVector2(-0.5, -0.5), rl.NewVector2(-0.5, 0)},
		{rl.NewVector2(-0.5, 0), rl.NewVector2(0.5, 0)},
		{rl.NewVector2(0.5, 0), rl.NewVector2(0.5, 0.5)},
		{rl.NewVector2(0.5, 0.5), rl.NewVector2(-0.5, 0.5)},
	},
	{
		{rl.NewVector2(0.5, -0.5), rl.NewVector2(-0.5, -0.5)},
		{rl.NewVector2(-0.5, -0.5), rl.NewVector2(-0.5, 0.5)},
		{rl.NewVector2(-0.5, 0.5), rl.NewVector2(0.5, 0.5)},
		{rl.NewVector2(0.5, 0.5), rl.NewVector2(0.5, 0)},
		{rl.NewVector2(0.5, 0), rl.NewVector2(-0.5, 0)},
	},
	{
		{rl.NewVector2(-0.5, -0.5), rl.NewVector2(0.5, -0.5)},
		{rl.NewVector2(0.5, -0.5), rl.NewVector2(0.5, 0.5)},
	},
	{
		{rl.NewVector2(-0.5, -0.5), rl.NewVector2(0.5, -0.5)},
		{rl.NewVector2(0.5, -0.5), rl.NewVector2(0.5, 0.5)},
		{rl.NewVector2(0.5, 0.5), rl.NewVector2(-0.5, 0.5)},
		{rl.NewVector2(-0.5, 0.5), rl.NewVector2(-0.5, -0.5)},
		{rl.NewVector2(-0.5, 0), rl.NewVector2(0.5, 0)},
	},
	{
		{rl.NewVector2(-0.5, -0.5), rl.NewVector2(0.5, -0.5)},
		{rl.NewVector2(0.5, -0.5), rl.NewVector2(0.5, 0.5)},
		{rl.NewVector2(0.5, 0.5), rl.NewVector2(-0.5, 0.5)},
		{rl.NewVector2(-0.5, 0), rl.NewVector2(0.5, 0)},
		{rl.NewVector2(0.5, 0), rl.NewVector2(0.5, -0.5)},
		{rl.NewVector2(-0.5, 0), rl.NewVector2(-0.5, -0.5)},
	},
}

func drawLines(org rl.Vector2, scale float32, rot float32, points []rl.Vector2) {
	transformer := Transformer{org, scale, rot}

	for i := 0; i < len(points); i++ {
		rl.DrawLineEx(transformer.apply(points[i]), transformer.apply(points[(i+1)%len(points)]), THICKNESS, rl.White)
	}
}

func drawNumber(score int, pos rl.Vector2) {
	val := score
	digits := []int{}

	for val > 0 {
		digits = append(digits, val%10)
		val /= 10
	}

	pos.X -= SCALE * 1.1 * float32(len(digits))
	for i := len(digits) - 1; i >= 0; i-- {
		for _, line := range NUMBERS[digits[i]] {
			drawLines(pos, 25, 0, line)
		}
		pos.X += SCALE * 1.1
	}

	if len(digits) == 0 {
		for _, line := range NUMBERS[0] {
			drawLines(pos, SCALE, 0, line)
		}
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

func (s AsteroidSize) score() int {
	switch s {
	case SMALL:
		return 100
	case MEDIUM:
		return 50
	case BIG:
		return 20
	}
	return 0
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
		return 1.5
	case MEDIUM:
		return 1.1
	case BIG:
		return 0.7
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
	if state.reset {
		state.reset = false
		resetGame()
	}

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

		if rl.IsKeyPressed(rl.KeySpace) {
			state.bullets = append(state.bullets, Bullet{
				pos:   rl.Vector2Add(state.ship.pos, rl.Vector2Scale(shipDir, SCALE*0.5)),
				vel:   rl.Vector2Scale(shipDir, 8.0),
				ttl:   2.0,
				spawn: state.now,
			})

			state.ship.vel = rl.Vector2Add(state.ship.vel, rl.Vector2Scale(shipDir, -0.7))
		}

		for i := 0; i < len(state.bullets); i++ {
			bullet := &state.bullets[i]

			if !bullet.remove && state.now-bullet.spawn > 0.05 && rl.Vector2Distance(bullet.pos, state.ship.pos) < SCALE*0.7 {
				bullet.remove = true
				state.ship.deathTime = state.now

				for i := 0; i < 5; i++ {
					angle := 2 * math.Pi * rand.Float32()
					state.particles = append(state.particles, Particle{
						pos:    rl.Vector2Add(state.ship.pos, rl.NewVector2(rand.Float32()*3, rand.Float32()*3)),
						vel:    rl.Vector2Scale(rl.NewVector2(float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))), 2*rand.Float32()),
						ttl:    2.0,
						pType:  LINE,
						rot:    angle,
						len:    SCALE * (0.6 + (0.4 * rand.Float32())),
						radius: 0,
					})
				}
			}
		}
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

		if !state.ship.isDead() && state.now-state.ship.deathTime > 0.10 && rl.Vector2Distance(asteroid.pos, state.ship.pos) < asteroid.size.size()*asteroid.size.collisionScale() {
			state.ship.deathTime = state.now

			for i := 0; i < 5; i++ {
				angle := 2 * math.Pi * rand.Float32()
				state.particles = append(state.particles, Particle{
					pos:    rl.Vector2Add(state.ship.pos, rl.NewVector2(rand.Float32()*3, rand.Float32()*3)),
					vel:    rl.Vector2Scale(rl.NewVector2(float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))), 2*rand.Float32()),
					ttl:    2.0,
					pType:  LINE,
					rot:    angle,
					len:    SCALE * (0.6 + (0.4 * rand.Float32())),
					radius: 0,
				})
			}

			hitAsteroid(asteroid, rl.Vector2Normalize(state.ship.vel))
		}

		for j := 0; j < len(state.aliens); j++ {
			alien := &state.aliens[j]

			if !alien.remove && rl.Vector2Distance(asteroid.pos, alien.pos) < asteroid.size.size()*asteroid.size.collisionScale() {
				alien.remove = true
				hitAsteroid(asteroid, rl.Vector2Normalize(rl.Vector2Subtract(asteroid.pos, alien.pos)))
			}
		}

		for j := 0; j < len(state.bullets); j++ {
			bullet := &state.bullets[j]
			if !bullet.remove && rl.Vector2Distance(asteroid.pos, bullet.pos) < asteroid.size.size()*asteroid.size.collisionScale() {
				bullet.remove = true
				hitAsteroid(asteroid, rl.Vector2Normalize(bullet.vel))
			}
		}

		if asteroid.remove {
			state.asteroids = append(state.asteroids[:i], state.asteroids[i+1:]...)
			i--
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

	for i := 0; i < len(state.bullets); i++ {
		bullet := &state.bullets[i]
		bullet.pos = rl.Vector2Add(bullet.pos, bullet.vel)

		if bullet.pos.X < 0 {
			bullet.pos.X = SIZE.X
		} else if bullet.pos.X > SIZE.X {
			bullet.pos.X = 0
		} else if bullet.pos.Y < 0 {
			bullet.pos.Y = SIZE.Y
		} else if bullet.pos.Y > SIZE.Y {
			bullet.pos.Y = 0
		}

		if !bullet.remove && bullet.ttl > state.delta {
			bullet.ttl -= state.delta
		} else {
			state.bullets = append(state.bullets[:i], state.bullets[i+1:]...)
			i--
		}
	}

	for i := 0; i < len(state.aliens); i++ {
		alien := &state.aliens[i]
		size := func() float32 {
			switch alien.size {
			case TINY:
				return SCALE * 0.5
			case HUGE:
				return SCALE * 0.8
			}
			return 0.0
		}()

		for j := 0; j < len(state.bullets); j++ {
			bullet := &state.bullets[j]
			if !bullet.remove && state.now-bullet.spawn > 0.05 && rl.Vector2Distance(alien.pos, bullet.pos) < size {
				bullet.remove = true
				alien.remove = true
			}
		}

		if !alien.remove && rl.Vector2Distance(alien.pos, state.ship.pos) < size {
			alien.remove = true
			state.ship.deathTime = state.now
		}

		if !alien.remove {
			dirChangeTime := func() float32 {
				switch alien.size {
				case TINY:
					return 0.35
				case HUGE:
					return 0.85
				}
				return 0.0
			}()

			shotTime := func() float32 {
				switch alien.size {
				case TINY:
					return 0.75
				case HUGE:
					return 1.25
				}
				return 0.0
			}()

			speed := func() float32 {
				switch alien.size {
				case TINY:
					return 6
				case HUGE:
					return 3
				}
				return 0.0
			}()

			if state.now-alien.lastDir > dirChangeTime {
				alien.lastDir = state.now
				angle := 2 * math.Pi * rand.Float32()
				alien.dir = rl.NewVector2(float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle))))
			}

			alien.pos = rl.Vector2Add(alien.pos, rl.Vector2Scale(alien.dir, speed))
			if alien.pos.X < 0 {
				alien.pos.X = SIZE.X
			} else if alien.pos.X > SIZE.X {
				alien.pos.X = 0
			} else if alien.pos.Y < 0 {
				alien.pos.Y = SIZE.Y
			} else if alien.pos.Y > SIZE.Y {
				alien.pos.Y = 0
			}

			if state.now-alien.lastShot > shotTime {
				alien.lastShot = state.now
				dir := rl.Vector2Normalize(rl.Vector2Subtract(state.ship.pos, alien.pos))

				state.bullets = append(state.bullets, Bullet{
					pos:   rl.Vector2Add(alien.pos, rl.Vector2Scale(dir, SCALE*0.5)),
					vel:   rl.Vector2Scale(dir, 6.0),
					ttl:   2.0,
					spawn: state.now,
				})
			}
		}

		if alien.remove {
			state.aliens = append(state.aliens[:i], state.aliens[i+1:]...)
			i--
		}
	}

	if state.ship.isDead() && state.now-state.ship.deathTime > 2.0 {
		resetStage()
	}

	if len(state.asteroids) == 0 && len(state.aliens) == 0 {
		resetAsteroids()
	}

	if state.score%5000 == 0 && state.score != 0 {
		state.aliens = append(state.aliens, Alien{
			pos:  rl.NewVector2(SIZE.X/2, SIZE.Y/2),
			dir:  rl.NewVector2(0, 0),
			size: HUGE,
		})
	}

	if state.score%8000 == 0 && state.score != 0 {
		state.aliens = append(state.aliens, Alien{
			pos:  rl.NewVector2(SIZE.X/2, SIZE.Y/2),
			dir:  rl.NewVector2(0, 0),
			size: TINY,
		})
	}
}

func drawAlien(pos rl.Vector2, size AlienSize) {
	scale := func() float32 {
		switch size {
		case TINY:
			return SCALE * 0.8
		case HUGE:
			return SCALE * 1.3
		}
		return 0.0
	}()

	drawLines(pos, scale, 0, []rl.Vector2{
		rl.NewVector2(-0.5, 0.0),
		rl.NewVector2(-0.3, -0.3),
		rl.NewVector2(0.3, -0.3),
		rl.NewVector2(0.5, 0.0),
		rl.NewVector2(0.3, 0.3),
		rl.NewVector2(-0.3, 0.3),
		rl.NewVector2(-0.5, 0.0),
		rl.NewVector2(0.5, 0.0),
	})

	drawLines(pos, scale, 0, []rl.Vector2{
		rl.NewVector2(-0.2, -0.3),
		rl.NewVector2(-0.1, -0.5),
		rl.NewVector2(0.1, -0.5),
		rl.NewVector2(0.2, -0.3),
	})
}

func render() {
	for i := 0; i < state.lives; i++ {
		drawLines(rl.NewVector2(20+float32(i)*SCALE, 20), SCALE, -math.Pi, []rl.Vector2{
			rl.NewVector2(-0.4, -0.5),
			rl.NewVector2(0.0, 0.5),
			rl.NewVector2(0.4, -0.5),
			rl.NewVector2(0.3, -0.4),
			rl.NewVector2(-0.3, -0.4),
		})
	}

	drawNumber(state.score, rl.NewVector2(SIZE.X-SCALE, SCALE))

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

	for _, alien := range state.aliens {
		drawAlien(alien.pos, alien.size)
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

	for _, bullet := range state.bullets {
		rl.DrawCircleV(bullet.pos, max(SCALE*0.06, 1), rl.White)
	}
}

func hitAsteroid(a *Asteroid, impact rl.Vector2) {
	state.score += a.size.score()
	a.remove = true

	for i := 0; i < 10; i++ {
		angle := 2 * math.Pi * rand.Float32()
		state.particles = append(state.particles, Particle{
			pos:    rl.Vector2Add(a.pos, rl.NewVector2(rand.Float32()*3, rand.Float32()*3)),
			vel:    rl.Vector2Scale(rl.NewVector2(float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))), 2.0+(4.0*rand.Float32())),
			ttl:    0.6 + (0.4 * rand.Float32()),
			pType:  DOT,
			rot:    0,
			len:    0,
			radius: SCALE * 0.025,
		})
	}

	if a.size == SMALL {
		return
	}

	for i := 0; i < 2; i++ {
		dir := rl.Vector2Normalize(a.vel)
		size := func() AsteroidSize {
			switch a.size {
			case BIG:
				return MEDIUM
			case MEDIUM:
				return SMALL
			default:
				return SMALL
			}
		}()
		state.asteroids = append(state.asteroids, Asteroid{
			pos:  a.pos,
			vel:  rl.Vector2Add(rl.Vector2Scale(dir, a.size.velocity()*3*rand.Float32()), rl.Vector2Scale(impact, 1.5)),
			size: size,
			seed: rand.Int63(),
		})
	}
}

func resetAsteroids() {
	state.asteroids = []Asteroid{}

	for i := 0; i < 20+state.score/2000; i++ {
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
	state.lives = 3
	state.score = 0
	resetStage()
	resetAsteroids()
}

func resetStage() {
	if state.ship.isDead() {
		if state.lives == 0 {
			state.reset = true
		} else {
			state.lives--
		}
	}

	state.ship.deathTime = 0.0
	state.ship = Ship{
		pos: rl.Vector2Scale(SIZE, 0.5),
		vel: rl.NewVector2(0, 0),
		rot: 0.0,
	}

	random := rand.Intn(2)
	randomEntry := func() float32 {
		if random == 0 {
			return SCALE
		} else {
			return SIZE.X - SCALE
		}
	}()
	state.aliens = append(state.aliens, Alien{
		pos:    rl.NewVector2(randomEntry, rand.Float32()*SIZE.Y),
		dir:    rl.NewVector2(0, 0),
		size:   HUGE,
		remove: false,
	})
}

func main() {
	rl.InitWindow(int32(SIZE.X), int32(SIZE.Y), "Asteroids Game")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	resetGame()

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
