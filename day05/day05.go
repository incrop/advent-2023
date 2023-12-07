package day05

import (
	"advent/utils"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	first int
	len   int
}

func (r Range) last() int {
	return r.first + r.len - 1
}

type RangeMapping struct {
	dstFirst int
	srcRange Range
}

type Mapping struct {
	from   string
	to     string
	ranges []RangeMapping
}

func (mapping Mapping) lookup(srcRange Range) []Range {
	rangeMappings := mapping.ranges
	sort.Slice(rangeMappings, func(i, j int) bool {
		return rangeMappings[i].srcRange.first < rangeMappings[j].srcRange.first
	})
	dstRanges := make([]Range, 0, 1)
	for _, rm := range rangeMappings {
		if rm.srcRange.last() < srcRange.first {
			continue
		}
		if srcRange.last() < rm.srcRange.first {
			break
		}
		if srcRange.first < rm.srcRange.first {
			offset := rm.srcRange.first - srcRange.first
			dstRanges = append(dstRanges, Range{srcRange.first, offset})
			srcRange = Range{srcRange.first + offset, srcRange.len - offset}
		}
		offest := srcRange.first - rm.srcRange.first
		if srcRange.last() <= rm.srcRange.last() {
			dstRanges = append(dstRanges, Range{rm.dstFirst + offest, srcRange.len})
			return dstRanges
		}
		insideLen := rm.srcRange.len - offest
		dstRanges = append(dstRanges, Range{rm.dstFirst + offest, insideLen})
		srcRange = Range{srcRange.first + insideLen, srcRange.len - insideLen}
	}
	dstRanges = append(dstRanges, srcRange)
	return dstRanges
}

func (mapping Mapping) lookupAll(srcs []Range) []Range {
	dstsNested := make([][]Range, len(srcs))
	totalLen := 0
	for i, src := range srcs {
		dsts := mapping.lookup(src)
		dstsNested[i] = dsts
		totalLen += len(dsts)
	}
	dstsFlat := make([]Range, 0, totalLen)
	for _, dsts := range dstsNested {
		dstsFlat = append(dstsFlat, dsts...)
	}
	return dstsFlat
}

type Mappings []Mapping

func (mappings Mappings) find(from string) Mapping {
	for _, mapping := range mappings {
		if mapping.from == from {
			return mapping
		}
	}
	panic("not found")
}

func (mappings Mappings) lookupChain(from string, to string, src Range) []Range {
	curr := from
	currRanges := []Range{src}
	for curr != to {
		mapping := mappings.find(curr)
		currRanges = mapping.lookupAll(currRanges)
		curr = mapping.to
	}
	return currRanges
}

type Almanac struct {
	seeds    []Range
	mappings Mappings
}

func (almanac Almanac) minLocation() int {

	minLocation := math.MaxInt
	for _, seedRange := range almanac.seeds {
		locationRanges := almanac.mappings.lookupChain("seed", "location", seedRange)
		for _, locationRange := range locationRanges {
			if locationRange.first < minLocation {
				minLocation = locationRange.first
			}
		}
	}
	return minLocation
}

const SEEDS_ARE_IN_RANGE_FORMAT = true

func (almanac *Almanac) parseSeeds(line string) bool {
	seedsStr, found := strings.CutPrefix(line, "seeds: ")
	if !found {
		return false
	}
	seedFields := strings.Fields(seedsStr)
	if SEEDS_ARE_IN_RANGE_FORMAT {
		almanac.seeds = make([]Range, len(seedFields)/2)
		for i := 0; i < len(seedFields); i += 2 {
			first, _ := strconv.Atoi(seedFields[i])
			len, _ := strconv.Atoi(seedFields[i+1])
			almanac.seeds[i/2] = Range{first, len}
		}
	} else {
		almanac.seeds = make([]Range, len(seedFields))
		for i, seedStr := range seedFields {
			seed, _ := strconv.Atoi(seedStr)
			almanac.seeds[i] = Range{seed, 1}
		}
	}
	return true
}

var mappingHeaderRe = regexp.MustCompile(`(\w+)-to-(\w+) map:`)

func (almanac *Almanac) parseMappingHeader(line string) bool {
	match := mappingHeaderRe.FindStringSubmatch(line)
	if match == nil {
		return false
	}
	almanac.mappings = append(almanac.mappings, Mapping{match[1], match[2], nil})
	return true
}

var mappingRangeRe = regexp.MustCompile(`(\d+) (\d+) (\d+)`)

func (almanac *Almanac) parseMappingRange(line string) bool {
	match := mappingRangeRe.FindStringSubmatch(line)
	if match == nil {
		return false
	}

	rangeNums := [3]int{}
	for i, str := range match[1:] {
		num, _ := strconv.Atoi(str)
		rangeNums[i] = num
	}
	ranges := &almanac.mappings[len(almanac.mappings)-1].ranges
	*ranges = append(*ranges, RangeMapping{rangeNums[0], Range{rangeNums[1], rangeNums[2]}})
	return true
}

func parseAlmanac(almanac Almanac, line string) Almanac {
	if line == "" {
		return almanac
	}
	if almanac.parseSeeds(line) {
		return almanac
	}
	if almanac.parseMappingHeader(line) {
		return almanac
	}
	if almanac.parseMappingRange(line) {
		return almanac
	}
	return almanac
}

func Run() int {
	almanac := utils.ProcessInput("day05.txt", Almanac{}, utils.Identity, parseAlmanac)
	return almanac.minLocation()
}
