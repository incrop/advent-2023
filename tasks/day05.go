package tasks

import (
	"advent/utils"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type d05_Range struct {
	first int
	len   int
}

func (r d05_Range) last() int {
	return r.first + r.len - 1
}

type d05_RangeMapping struct {
	dstFirst int
	srcRange d05_Range
}

type d05_Mapping struct {
	from   string
	to     string
	ranges []d05_RangeMapping
}

func (mapping d05_Mapping) lookup(srcRange d05_Range) []d05_Range {
	rangeMappings := mapping.ranges
	sort.Slice(rangeMappings, func(i, j int) bool {
		return rangeMappings[i].srcRange.first < rangeMappings[j].srcRange.first
	})
	dstRanges := make([]d05_Range, 0, 1)
	for _, rm := range rangeMappings {
		if rm.srcRange.last() < srcRange.first {
			continue
		}
		if srcRange.last() < rm.srcRange.first {
			break
		}
		if srcRange.first < rm.srcRange.first {
			offset := rm.srcRange.first - srcRange.first
			dstRanges = append(dstRanges, d05_Range{srcRange.first, offset})
			srcRange = d05_Range{srcRange.first + offset, srcRange.len - offset}
		}
		offest := srcRange.first - rm.srcRange.first
		if srcRange.last() <= rm.srcRange.last() {
			dstRanges = append(dstRanges, d05_Range{rm.dstFirst + offest, srcRange.len})
			return dstRanges
		}
		insideLen := rm.srcRange.len - offest
		dstRanges = append(dstRanges, d05_Range{rm.dstFirst + offest, insideLen})
		srcRange = d05_Range{srcRange.first + insideLen, srcRange.len - insideLen}
	}
	dstRanges = append(dstRanges, srcRange)
	return dstRanges
}

func (mapping d05_Mapping) lookupAll(srcs []d05_Range) []d05_Range {
	dstsNested := make([][]d05_Range, len(srcs))
	totalLen := 0
	for i, src := range srcs {
		dsts := mapping.lookup(src)
		dstsNested[i] = dsts
		totalLen += len(dsts)
	}
	dstsFlat := make([]d05_Range, 0, totalLen)
	for _, dsts := range dstsNested {
		dstsFlat = append(dstsFlat, dsts...)
	}
	return dstsFlat
}

type d05_Mappings []d05_Mapping

func (mappings d05_Mappings) find(from string) d05_Mapping {
	for _, mapping := range mappings {
		if mapping.from == from {
			return mapping
		}
	}
	panic("not found")
}

func (mappings d05_Mappings) lookupChain(from string, to string, src d05_Range) []d05_Range {
	curr := from
	currRanges := []d05_Range{src}
	for curr != to {
		mapping := mappings.find(curr)
		currRanges = mapping.lookupAll(currRanges)
		curr = mapping.to
	}
	return currRanges
}

type d05_Almanac struct {
	seeds    []d05_Range
	mappings d05_Mappings
}

func (almanac d05_Almanac) minLocation() int {

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

func (almanac *d05_Almanac) parseSeeds(line string) bool {
	seedsStr, found := strings.CutPrefix(line, "seeds: ")
	if !found {
		return false
	}
	seedFields := strings.Fields(seedsStr)
	if SEEDS_ARE_IN_RANGE_FORMAT {
		almanac.seeds = make([]d05_Range, len(seedFields)/2)
		for i := 0; i < len(seedFields); i += 2 {
			first, _ := strconv.Atoi(seedFields[i])
			len, _ := strconv.Atoi(seedFields[i+1])
			almanac.seeds[i/2] = d05_Range{first, len}
		}
	} else {
		almanac.seeds = make([]d05_Range, len(seedFields))
		for i, seedStr := range seedFields {
			seed, _ := strconv.Atoi(seedStr)
			almanac.seeds[i] = d05_Range{seed, 1}
		}
	}
	return true
}

var mappingHeaderRe = regexp.MustCompile(`(\w+)-to-(\w+) map:`)

func (almanac *d05_Almanac) parseMappingHeader(line string) bool {
	match := mappingHeaderRe.FindStringSubmatch(line)
	if match == nil {
		return false
	}
	almanac.mappings = append(almanac.mappings, d05_Mapping{match[1], match[2], nil})
	return true
}

var mappingRangeRe = regexp.MustCompile(`(\d+) (\d+) (\d+)`)

func (almanac *d05_Almanac) parseMappingRange(line string) bool {
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
	*ranges = append(*ranges, d05_RangeMapping{rangeNums[0], d05_Range{rangeNums[1], rangeNums[2]}})
	return true
}

func d05_parseAlmanac(almanac d05_Almanac, line string) d05_Almanac {
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

func Day05() int {
	almanac := utils.ProcessInput("day05.txt", d05_Almanac{}, utils.Identity, d05_parseAlmanac)
	return almanac.minLocation()
}
