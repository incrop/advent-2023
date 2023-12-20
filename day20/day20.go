package day20

import (
	"advent/utils"
	"fmt"
	"regexp"
)

type Frequency int

const (
	None Frequency = -1
	Low  Frequency = 0
	High Frequency = 1
)

type Pulse struct {
	src  string
	dst  string
	freq Frequency
}

func (pulse Pulse) String() string {
	return fmt.Sprintf("%v -%v-> %v", pulse.src, pulse.freq, pulse.dst)
}

type Module interface {
	name() string
	dsts() []string
	connect(src string)
	apply(in Pulse) Frequency
	printMermaid()
}

type Modules map[string]Module

var prevs map[string]int = make(map[string]int)

func (modules Modules) pushButton(n int) (lowNum int, highNum int) {
	queue := make([]Pulse, 0)
	queue = append(queue, Pulse{"button", "broadcaster", Low})
	for len(queue) > 0 {
		pulse := queue[0]
		if pulse.dst == "vd" && pulse.freq == 1 {
			prev := prevs[pulse.src]
			fmt.Printf("%s: %d (+%d)\n", pulse.src, n, n-prev)
			prevs[pulse.src] = n
		}
		queue = queue[1:]
		switch pulse.freq {
		case Low:
			lowNum++
		case High:
			highNum++
		}
		module := modules[pulse.dst]
		outFreq := module.apply(pulse)
		if outFreq == None {
			continue
		}
		for _, dst := range module.dsts() {
			queue = append(queue, Pulse{module.name(), dst, outFreq})
		}
	}
	return
}

type Broadcaster []string

func (broadcaster Broadcaster) name() string {
	return "broadcaster"
}

func (broadcaster Broadcaster) dsts() []string {
	return broadcaster
}

func (broadcaster Broadcaster) connect(src string) {
	panic(src)
}

func (broadcaster Broadcaster) apply(in Pulse) Frequency {
	return in.freq
}

func (broadcaster Broadcaster) printMermaid() {
	for _, dst := range broadcaster {
		fmt.Printf("broadcaster:::LF --> %s\n", dst)
	}
}

type Sink string

func (sink Sink) name() string {
	return string(sink)
}

func (sink Sink) dsts() []string {
	return nil
}

func (sink Sink) connect(src string) {
}

func (sink Sink) apply(in Pulse) Frequency {
	return None
}

func (sink Sink) printMermaid() {
}

type FlipFlop struct {
	_name string
	state byte
	_dsts []string
}

func (flipFlop *FlipFlop) name() string {
	return flipFlop._name
}

func (flipFlop *FlipFlop) dsts() []string {
	return flipFlop._dsts
}

func (flipFlop *FlipFlop) connect(src string) {
}

func (flipFlop *FlipFlop) apply(in Pulse) Frequency {
	if in.freq == High {
		return None
	}
	flipFlop.state = 1 - flipFlop.state
	return Frequency(flipFlop.state)
}

func (flipFlop *FlipFlop) String() string {
	return fmt.Sprintf("&%v->%v", flipFlop.state, flipFlop._dsts)
}

func (flipFlop *FlipFlop) printMermaid() {
	f := "LF"
	if flipFlop.state == 1 {
		f = "HF"
	}
	for _, dst := range flipFlop._dsts {
		fmt.Printf("%s(%s):::%s --> %s\n", flipFlop._name, flipFlop._name, f, dst)
	}
}

type Conjunction struct {
	_name string
	srcs  map[string]Frequency
	_dsts []string
}

func (conjunction *Conjunction) name() string {
	return conjunction._name
}

func (conjunction *Conjunction) dsts() []string {
	return conjunction._dsts
}

func (conjunction *Conjunction) connect(src string) {
	conjunction.srcs[src] = Low
}

func (conjunction *Conjunction) apply(in Pulse) Frequency {
	conjunction.srcs[in.src] = in.freq
	for _, inFreq := range conjunction.srcs {
		if inFreq == Low {
			return High
		}
	}
	return Low
}

func (conjunction *Conjunction) String() string {
	return fmt.Sprintf("&%v->%v", conjunction.srcs, conjunction._dsts)
}

func (conjunction *Conjunction) printMermaid() {
	f := "HF"
	for _, inFreq := range conjunction.srcs {
		if inFreq == Low {
			f = "LF"
			break
		}
	}
	for _, dst := range conjunction._dsts {
		fmt.Printf("%s{%s}:::%s --> %s\n", conjunction._name, conjunction._name, f, dst)
	}
}

var moduleRe = regexp.MustCompile(`^([&%]?)(\w+) -> (\w+(?:, \w+)*)$`)

func parseModule(line string) Module {
	match := moduleRe.FindStringSubmatch(line)
	kind := match[1]
	name := match[2]
	dsts := utils.Fields(match[3], ", ")
	switch {
	case name == "broadcaster":
		return Broadcaster(dsts)
	case kind == "%":
		return &FlipFlop{name, 0, dsts}
	case kind == "&":
		return &Conjunction{name, make(map[string]Frequency), dsts}
	default:
		panic("todo")
	}
}

func appendModule(modules Modules, module Module) Modules {
	modules[module.name()] = module
	return modules
}

func (modules Modules) init() Modules {
	for src, srcModule := range modules {
		for _, dst := range srcModule.dsts() {
			dstModule, exist := modules[dst]
			if exist {
				dstModule.connect(src)
			} else {
				modules[dst] = Sink(dst)
			}
		}
	}
	return modules
}

func (modules Modules) printMermaid() {
	fmt.Println("flowchart TD")
	modules["broadcaster"].printMermaid()
	for name, module := range modules {
		if name == "broadcaster" {
			continue
		}
		module.printMermaid()
	}
	fmt.Println("classDef HF fill:#ffa400")
	fmt.Println("classDef LF fill:#7edf68")
}

func (modules Modules) productAfterPushingButton(times int) int {
	totalLow, totalHigh := 0, 0
	for i := 0; i < times; i++ {
		low, high := modules.pushButton(i)
		totalLow += low
		totalHigh += high
	}
	return totalHigh * totalLow
}

func Run() int {
	modules := utils.ProcessInput("day20.txt", Modules{}, parseModule, appendModule).init()
	// defer modules.printMermaid()
	return modules.productAfterPushingButton(100000)
}
