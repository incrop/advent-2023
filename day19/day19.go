package day19

import (
	"advent/utils"
	"regexp"
	"strconv"
)

type System struct {
	workflows map[string]Workflow
	parts     []Part
}

type Workflow struct {
	name  string
	rules []Rule
}

type Rule struct {
	condition Condition
	outcome   string
}

type Condition interface {
	matches(part Part) bool
	split(parts PartRange) [2]PartRange
}

type Always bool

func (always Always) matches(_ Part) bool {
	return bool(always)
}

func (always Always) split(parts PartRange) (split [2]PartRange) {
	split[1] = parts
	return
}

type Compare struct {
	category string
	op       string
	value    int
}

func (cmp Compare) matches(part Part) bool {
	value := part[cmp.category]
	switch cmp.op {
	case "<":
		return value < cmp.value
	case ">":
		return cmp.value < value
	default:
		panic("wat")
	}
}

func (cmp Compare) split(parts PartRange) (split [2]PartRange) {
	minMax := parts[cmp.category]
	min, max := minMax[0], minMax[1]
	switch {
	case cmp.op == "<" && max < cmp.value,
		cmp.op == ">" && min > cmp.value:
		split[0] = nil
		split[1] = parts
	case cmp.op == "<" && (cmp.value <= min || min == max),
		cmp.op == ">" && (cmp.value >= max || min == max):
		split[0] = parts
		split[1] = nil
	default:
		split[0] = make(PartRange)
		split[1] = make(PartRange)
		for category, minMax := range parts {
			if category != cmp.category {
				split[0][category] = minMax
				split[1][category] = minMax
			} else if cmp.op == "<" {
				split[0][category] = [2]int{cmp.value, max}
				split[1][category] = [2]int{min, cmp.value - 1}
			} else if cmp.op == ">" {
				split[0][category] = [2]int{min, cmp.value}
				split[1][category] = [2]int{cmp.value + 1, max}
			}
		}
	}
	return
}

type Part map[string]int

type Subset struct {
	workflowName string
	ruleIdx      int
	parts        PartRange
}

type PartRange map[string][2]int

func (parts PartRange) size() (product int) {
	product = 1
	for _, r := range parts {
		product *= r[1] - r[0] + 1
	}
	return
}

type WorkflowOrPart struct {
	workflow *Workflow
	part     Part
}

func (system System) isAccepted(part Part) bool {
	workflow := system.workflows["in"]
OUTER:
	for {
		for _, rule := range workflow.rules {
			if rule.condition.matches(part) {
				switch rule.outcome {
				case "A":
					return true
				case "R":
					return false
				default:
					workflow = system.workflows[rule.outcome]
					continue OUTER
				}
			}
		}
		panic("unreachable")
	}
}

func (system System) sumOfApprovedPartsCategories() (sum int) {
	for _, part := range system.parts {
		if system.isAccepted(part) {
			for _, value := range part {
				sum += value
			}
		}
	}
	return
}

func (system System) numberOfAcceptedCombinations() (total int) {
	subsets := []Subset{{
		workflowName: "in",
		ruleIdx:      0,
		parts: PartRange{
			"x": {1, 4000},
			"m": {1, 4000},
			"a": {1, 4000},
			"s": {1, 4000},
		},
	}}
	for len(subsets) > 0 {
		next := make([]Subset, 0)
		for _, subset := range subsets {
			workflow := system.workflows[subset.workflowName]
			rule := workflow.rules[subset.ruleIdx]
			subParts := rule.condition.split(subset.parts)
			if subParts[0] != nil {
				next = append(next, Subset{
					workflowName: subset.workflowName,
					ruleIdx:      subset.ruleIdx + 1,
					parts:        subParts[0],
				})
			}
			if subParts[1] != nil {
				switch rule.outcome {
				case "A":
					total += subParts[1].size()
				case "R":
				default:
					next = append(next, Subset{
						workflowName: rule.outcome,
						ruleIdx:      0,
						parts:        subParts[1],
					})
				}

			}
		}
		subsets = next
	}
	return
}

var workflowOrPartRe = regexp.MustCompile(`^(\w*)\{(.+)\}$`)
var compareRe = regexp.MustCompile(`^(\w+)([<>])(\d+)$`)

func parseWorkflow(name string, rulesStr string) (workflow *Workflow) {
	workflow = &Workflow{name: name}
	for _, ruleStr := range utils.Fields(rulesStr, ",") {
		condOut := utils.Fields(ruleStr, ":")
		if len(condOut) == 1 {
			workflow.rules = append(workflow.rules, Rule{Always(true), condOut[0]})
			continue
		}
		match := compareRe.FindStringSubmatch(condOut[0])
		value, _ := strconv.Atoi(match[3])
		workflow.rules = append(workflow.rules, Rule{Compare{match[1], match[2], value}, condOut[1]})
	}
	return
}

func parsePart(categoriesStr string) (part Part) {
	part = make(Part)
	for _, categoryStr := range utils.Fields(categoriesStr, ",") {
		keyValue := utils.Fields(categoryStr, "=")
		value, _ := strconv.Atoi(keyValue[1])
		part[keyValue[0]] = value
	}
	return
}

func parseLine(line string) (wop WorkflowOrPart) {
	if line == "" {
		return
	}
	match := workflowOrPartRe.FindStringSubmatch(line)
	if match[1] != "" {
		wop.workflow = parseWorkflow(match[1], match[2])
	} else {
		wop.part = parsePart(match[2])
	}
	return
}

func aggregate(system System, wop WorkflowOrPart) System {
	if wop.workflow != nil {
		if system.workflows == nil {
			system.workflows = make(map[string]Workflow)
		}
		system.workflows[wop.workflow.name] = *wop.workflow
	}
	if wop.part != nil {
		system.parts = append(system.parts, wop.part)
	}
	return system
}

func Run() int {
	system := utils.ProcessInput("day19.txt", System{}, parseLine, aggregate)
	return system.numberOfAcceptedCombinations()
}
