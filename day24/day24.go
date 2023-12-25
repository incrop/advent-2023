package day24

import (
	"advent/utils"
	"math/big"
	"strconv"
)

type Coords struct {
	x, y, z int64
}

type Hail struct {
	pos, vel Coords
}

type Storm []Hail

func (storm Storm) countHailsPathsIntersectingXY(minXY, maxXY int) (count int) {
	if len(storm) < 2 {
		return
	}
	minXYr := big.NewRat(int64(minXY), 1)
	maxXYr := big.NewRat(int64(maxXY), 1)
	for i, h1 := range storm[:len(storm)-1] {
		for _, h2 := range storm[i+1:] {
			if h1.isPathIntersectsXY(h2, minXYr, maxXYr) {
				count++
			}
		}
	}
	return
}

// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line
func (h1 Hail) isPathIntersectsXY(h2 Hail, minXY, maxXY *big.Rat) bool {
	a1, b1 := big.NewInt(h1.vel.y), big.NewInt(-h1.vel.x)
	c1 := new(big.Int).Sub(
		new(big.Int).Mul(big.NewInt(h1.vel.x), big.NewInt(h1.pos.y)),
		new(big.Int).Mul(big.NewInt(h1.vel.y), big.NewInt(h1.pos.x)),
	)
	a2, b2 := big.NewInt(h2.vel.y), big.NewInt(-h2.vel.x)
	c2 := new(big.Int).Sub(
		new(big.Int).Mul(big.NewInt(h2.vel.x), big.NewInt(h2.pos.y)),
		new(big.Int).Mul(big.NewInt(h2.vel.y), big.NewInt(h2.pos.x)),
	)
	numX := new(big.Int).Sub(
		new(big.Int).Mul(b1, c2),
		new(big.Int).Mul(b2, c1),
	)
	numY := new(big.Int).Sub(
		new(big.Int).Mul(c1, a2),
		new(big.Int).Mul(c2, a1),
	)
	den := new(big.Int).Sub(
		new(big.Int).Mul(a1, b2),
		new(big.Int).Mul(a2, b1),
	)
	if den.Sign() == 0 {
		// paths are parallel
		return false
	}
	x := new(big.Rat).SetFrac(numX, den)
	y := new(big.Rat).SetFrac(numY, den)

	res := h1.pointsToXY(x, y) && h2.pointsToXY(x, y)
	if !res {
		return false
	}
	if x.Cmp(minXY) < 0 || x.Cmp(maxXY) > 0 || y.Cmp(minXY) < 0 || y.Cmp(maxXY) > 0 {
		return false
	}
	return res
}

func (h Hail) pointsToXY(x, y *big.Rat) bool {
	// fmt.Println()
	// fmt.Println(h, x.FloatString(3), y.FloatString(3))
	x = new(big.Rat).Set(x)
	x.Sub(x, big.NewRat(h.pos.x, 1))
	vx := big.NewRat(h.vel.x, 1)
	y = new(big.Rat).Set(y)
	y.Sub(y, big.NewRat(h.pos.y, 1))
	vy := big.NewRat(h.vel.y, 1)
	// fmt.Println("x, y =", x.FloatString(3), y.FloatString(3))
	// fmt.Println("vx, vy =", vx, vy)
	if vx.Sign() == 0 {
		// fmt.Println("vx.Sign() == 0")
		return x.Sign() == 0 && y.Sign() == vy.Sign()
	}
	if vy.Sign() == 0 {
		// fmt.Println("vy.Sign() == 0")
		return y.Sign() == 0 && x.Sign() == vx.Sign()
	}
	slope1 := y.Quo(y, x)
	slope2 := vy.Quo(vy, vx)
	// fmt.Println("slope1, slope2 =", slope1, slope2)
	if slope1.Cmp(slope2) != 0 {
		panic("should not happen")
	}
	return x.Sign() == vx.Sign()
}

func parseHail(line string) Hail {
	var nums [6]int64
	for i, numStr := range utils.Fields(line, "@, ") {
		n, _ := strconv.Atoi(numStr)
		nums[i] = int64(n)
	}
	return Hail{
		pos: Coords{x: nums[0], y: nums[1], z: nums[2]},
		vel: Coords{x: nums[3], y: nums[4], z: nums[5]},
	}
}

func appendHail(storm Storm, hail Hail) Storm {
	return append(storm, hail)
}

func Run() int {
	storm := utils.ProcessInput("day24.txt", nil, parseHail, appendHail)
	return storm.countHailsPathsIntersectingXY(200000000000000, 400000000000000)
	// storm := utils.ProcessInput("day24_test.txt", nil, parseHail, appendHail)
	// return storm.countHailsPathsIntersectingXY(7, 27)
}
