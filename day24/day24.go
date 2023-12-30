package day24

import (
	"advent/utils"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type Coords struct {
	x, y, z int64
}

type Hail struct {
	pos, vel Coords
}

func (hail Hail) dimX() HailDim {
	return HailDim{pos: hail.pos.x, vel: hail.vel.x}
}
func (hail Hail) dimY() HailDim {
	return HailDim{pos: hail.pos.y, vel: hail.vel.y}
}
func (hail Hail) dimZ() HailDim {
	return HailDim{pos: hail.pos.z, vel: hail.vel.z}
}

type HailDim struct {
	pos, vel int64
}

func (h1 HailDim) intersectTime(h2 HailDim) *big.Rat {
	if h1.vel == h2.vel {
		if h1.pos == h2.pos {
			return big.NewRat(0, 1)
		}
		return nil
	}
	return big.NewRat(h2.pos-h1.pos, h1.vel-h2.vel)
}

type Storm []Hail

func (storm Storm) checkCollisions(h1 Hail) {
	for _, h2 := range storm {
		times := [3]*big.Rat{
			h1.dimX().intersectTime(h2.dimX()),
			h1.dimY().intersectTime(h2.dimY()),
			h1.dimZ().intersectTime(h2.dimZ()),
		}
		for _, time := range times {
			if time == nil || time.Sign() < 0 {
				panic("time in the past")
			}
		}
		first := times[0]
		if first.Sign() == 0 {
			panic("nope")
		}
		for _, time := range times[1:] {
			if time.Sign() != 0 && time.Cmp(first) != 0 {
				panic("no intersection")
			}
		}
		fmt.Println("Collides with", h2, "at", first)
	}
}

type Matrix struct {
	coeff [][]*big.Rat
	rhs   []*big.Rat
}

func (m Matrix) String() string {
	var sb strings.Builder
	for i, row := range m.coeff {
		for _, elem := range row {
			fmt.Fprintf(&sb, "%v ", elem)
		}
		fmt.Fprintf(&sb, "= %v\n", m.rhs[i])
	}
	return sb.String()
}

func (m Matrix) solve() {
	for i, iRow := range m.coeff {
		div := new(big.Rat).Set(iRow[i])
		if div.Sign() == 0 {
			fmt.Println(m)
			panic("zero")
		}
		for _, elem := range iRow {
			elem.Quo(elem, div)
		}
		m.rhs[i].Quo(m.rhs[i], div)
		for j, jRow := range m.coeff {
			if i == j {
				continue
			}
			mul := new(big.Rat).Set(jRow[i])
			for k, elem := range jRow {
				elem.Sub(elem, new(big.Rat).Mul(iRow[k], mul))
			}
			m.rhs[j].Sub(m.rhs[j], new(big.Rat).Mul(m.rhs[i], mul))
		}
	}
}

// Kudos to @ash42 comment, havent figured out how to turn this into linear equations myself
// https://github.com/ash42/adventofcode/blob/95b412fe20da44002192e69d267733375241a9cd/adventofcode2023/src/nl/michielgraat/adventofcode2023/day24/Day24.java#L82-L123
func (storm Storm) findBullet() Hail {
	m := Matrix{
		coeff: make([][]*big.Rat, 6),
		rhs:   make([]*big.Rat, 6),
	}
	for i := 0; i < 3; i++ {
		h1 := storm[i]
		h2 := storm[i+1]
		x1, y1, z1 := h1.pos.x, h1.pos.y, h1.pos.z
		x2, y2, z2 := h2.pos.x, h2.pos.y, h2.pos.z
		vx1, vy1, vz1 := h1.vel.x, h1.vel.y, h1.vel.z
		vx2, vy2, vz2 := h2.vel.x, h2.vel.y, h2.vel.z
		m.coeff[i*2] = []*big.Rat{
			big.NewRat(vy2-vy1, 1),
			big.NewRat(vx1-vx2, 1),
			big.NewRat(0, 1),
			big.NewRat(y1-y2, 1),
			big.NewRat(x2-x1, 1),
			big.NewRat(0, 1),
		}
		m.rhs[i*2] = big.NewRat(-x1*vy1+y1*vx1+x2*vy2-y2*vx2, 1)
		m.coeff[i*2+1] = []*big.Rat{
			big.NewRat(vz2-vz1, 1),
			big.NewRat(0, 1),
			big.NewRat(vx1-vx2, 1),
			big.NewRat(z1-z2, 1),
			big.NewRat(0, 1),
			big.NewRat(x2-x1, 1),
		}
		m.rhs[i*2+1] = big.NewRat(-x1*vz1+z1*vx1+x2*vz2-z2*vx2, 1)
	}
	fmt.Println(m)
	m.solve()
	fmt.Println(m)

	var res [6]int64
	for i, br := range m.rhs {
		if br.Denom().Int64() != 1 {
			panic("not int")
		}
		res[i] = br.Num().Int64()
	}
	return Hail{
		pos: Coords{x: res[0], y: res[1], z: res[2]},
		vel: Coords{x: res[3], y: res[4], z: res[5]},
	}
}

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
	x = new(big.Rat).Set(x)
	x.Sub(x, big.NewRat(h.pos.x, 1))
	vx := big.NewRat(h.vel.x, 1)
	y = new(big.Rat).Set(y)
	y.Sub(y, big.NewRat(h.pos.y, 1))
	vy := big.NewRat(h.vel.y, 1)
	if vx.Sign() == 0 {
		return x.Sign() == 0 && y.Sign() == vy.Sign()
	}
	if vy.Sign() == 0 {
		return y.Sign() == 0 && x.Sign() == vx.Sign()
	}
	slope1 := y.Quo(y, x)
	slope2 := vy.Quo(vy, vx)
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
	// storm := utils.ProcessInput("day24_test.txt", nil, parseHail, appendHail)
	//storm.countHailsPathsIntersectingXY(200000000000000, 400000000000000)
	bullet := storm.findBullet()
	storm.checkCollisions(bullet)
	return int(bullet.pos.x + bullet.pos.y + bullet.pos.z)
}
