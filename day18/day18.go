package day18

import (
	"advent/utils"
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

type Direction rune

const (
	Up    Direction = 'U'
	Left  Direction = 'L'
	Right Direction = 'R'
	Down  Direction = 'D'
)

func (dir Direction) offset(i, j, len int) (int, int) {
	switch dir {
	case Right:
		return i, j + len
	case Down:
		return i + len, j
	case Left:
		return i, j - len
	case Up:
		return i - len, j
	}
	return i, j
}

type Trench struct {
	dir Direction
	len int
}

type Lagoon []Trench

type Cut struct {
	j1, j2 int
}

func (cut Cut) len() int {
	return cut.j2 - cut.j1 + 1
}

type Section struct {
	i    int
	cuts []Cut
}

func (section Section) len() (len int) {
	for _, cut := range section.cuts {
		len += cut.len()
	}
	return
}

func (cs *Section) merge(bs Section) (countFilled int) {
	var result []Cut
	cuts := cs.cuts
	borders := bs.cuts
	for len(cuts) > 0 && len(borders) > 0 {
		cut := &cuts[0]
		border := &borders[0]
		switch {
		case *cut == *border:
			cuts = cuts[1:]
			borders = borders[1:]
			countFilled += cut.len()
		case cut.j2 < border.j1:
			cuts = cuts[1:]
			result = append(result, *cut)
			countFilled += cut.len()
		case border.j2 < cut.j1:
			borders = borders[1:]
			result = append(result, *border)
			countFilled += border.len()
		case cut.j2 == border.j1:
			cuts = cuts[1:]
			border.j1 = cut.j1
		case border.j2 == cut.j1:
			borders = borders[1:]
			cut.j1 = border.j1
		case cut.j1 == border.j1 && border.j2 < cut.j2:
			borders = borders[1:]
			countFilled += border.len() - 1
			cut.j1 = border.j2
		case cut.j2 == border.j2 && cut.j1 < border.j1:
			borders = borders[1:]
			countFilled += border.len() - 1
			cut.j2 = border.j1
		case cut.j1 < border.j1 && border.j2 < cut.j2:
			result = append(result, Cut{cut.j1, border.j1})
			borders = borders[1:]
			countFilled += border.j2 - cut.j1
			cut.j1 = border.j2
		default:
			panic("unreachable?")
		}
	}
	for _, cut := range cuts {
		result = append(result, cut)
		countFilled += cut.len()
	}
	for _, border := range borders {
		result = append(result, border)
		countFilled += border.len()
	}
	*cs = Section{i: bs.i + 1, cuts: result}
	return
}

type Plan []Section

func (lagoon Lagoon) plan() (plan Plan) {
	plan = nil
	sections := make(map[int]*Section)
	i, j := 0, 0
	for _, trench := range lagoon {
		i1, j1 := trench.dir.offset(i, j, trench.len)
		if i == i1 {
			border := Cut{j, j1}
			if j1 < j {
				border.j1, border.j2 = border.j2, border.j1
			}
			section := sections[i]
			if section == nil {
				section = &Section{i: i}
				sections[i] = section
			}
			section.cuts = append(section.cuts, border)
		}
		i, j = i1, j1
	}
	for _, section := range sections {
		plan = append(plan, *section)
	}
	sort.Slice(plan, func(i, j int) bool {
		return plan[i].i < plan[j].i
	})
	for _, section := range plan {
		sort.Slice(section.cuts, func(i, j int) bool {
			return section.cuts[i].j1 < section.cuts[j].j1
		})
	}
	return
}

func (plan Plan) countFilled() (count int) {
	curr := Section{i: plan[0].i}
	for _, next := range plan {
		count += (next.i - curr.i) * curr.len()
		count += curr.merge(next)
	}
	return
}

var trenchRe = regexp.MustCompile(`^([ULRD]) (\d+) \(#([0-9a-f]{6})\)$`)

func parseTrench(line string) (trench Trench) {
	match := trenchRe.FindStringSubmatch(line)
	if match == nil {
		panic("regex does not match")
	}
	trench.dir = Direction(match[1][0])
	trench.len, _ = strconv.Atoi(match[2])
	return
}

func parseTrenchFixed(line string) (trench Trench) {
	match := trenchRe.FindStringSubmatch(line)
	if match == nil {
		panic("regex does not match")
	}
	dir := match[3][5]
	switch dir {
	case '0':
		trench.dir = 'R'
	case '1':
		trench.dir = 'D'
	case '2':
		trench.dir = 'L'
	case '3':
		trench.dir = 'U'
	}
	_, _ = fmt.Sscanf(match[3][0:5], "%05x", &trench.len)
	return
}

func appendRows(lagoon Lagoon, trench Trench) Lagoon {
	return append(lagoon, trench)
}

func Run() int {
	lagoon := utils.ProcessInput("day18.txt", Lagoon{}, parseTrenchFixed, appendRows)
	return lagoon.plan().countFilled()
}
