package day15

import (
	"advent/utils"
	"fmt"
	"regexp"
	"strings"
)

type Key string

type Entry struct {
	key   Key
	value byte
}

type EntryList struct {
	entry Entry
	next  *EntryList
}

type HashMap [256]*EntryList

func (hashMap *HashMap) String() string {
	var sb strings.Builder
	for box, elem := range hashMap {
		if elem == nil {
			continue
		}
		fmt.Fprintf(&sb, "Box %d:", box)
		for ; elem != nil; elem = elem.next {
			fmt.Fprintf(&sb, " [%s %d]", elem.entry.key, elem.entry.value)
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (hashMap *HashMap) totalFocusingPower() (sum int) {
	for box, elem := range hashMap {
		for slot := 1; elem != nil; slot, elem = slot+1, elem.next {
			sum += (box + 1) * slot * int(elem.entry.value)
		}
	}
	return
}

type Operation interface {
	applyTo(hashMap *HashMap)
}

func (putEntry Entry) applyTo(hashMap *HashMap) {
	boxHash := putEntry.key.hash()
	elem := hashMap[boxHash]
	if elem == nil {
		hashMap[boxHash] = &EntryList{putEntry, nil}
		return
	}
	prev := elem
	for ; elem != nil; prev, elem = elem, elem.next {
		if elem.entry.key == putEntry.key {
			elem.entry.value = putEntry.value
			return
		}
	}
	prev.next = &EntryList{putEntry, nil}
}

func (deleteKey Key) applyTo(hashMap *HashMap) {
	boxHash := deleteKey.hash()
	elem := hashMap[boxHash]
	if elem == nil {
		return
	}
	if elem.entry.key == deleteKey {
		hashMap[boxHash] = elem.next
	}
	for ; elem.next != nil; elem = elem.next {
		if elem.next.entry.key == deleteKey {
			elem.next = elem.next.next
			return
		}
	}
}

func (step Key) hash() byte {
	value := 0
	for _, r := range step {
		value = ((value + int(r)) * 17) % 256
	}
	return byte(value)
}

func parseKeys(line string) []Key {
	fields := utils.Fields(line, ",")
	steps := make([]Key, len(fields))
	for i, field := range fields {
		steps[i] = Key(field)
	}
	return steps
}

func sumHash(sum int, steps []Key) int {
	for _, step := range steps {
		sum += int(step.hash())
	}
	return sum
}

var putRe = regexp.MustCompile(`^(\w+)=(\d)$`)
var deleteRe = regexp.MustCompile(`^(\w+)-$`)

func parseOperations(line string) []Operation {
	fields := utils.Fields(line, ",")
	operations := make([]Operation, len(fields))
	for i, field := range fields {
		if match := putRe.FindStringSubmatch(field); match != nil {
			key := Key(match[1])
			value := byte(match[2][0] - '0')
			operations[i] = Entry{key, value}
		} else if match := deleteRe.FindStringSubmatch(field); match != nil {
			operations[i] = Key(match[1])
		} else {
			panic("Unknown operation")
		}
	}
	return operations
}

func applyOperations(hashMap *HashMap, operations []Operation) *HashMap {
	for _, operation := range operations {
		operation.applyTo(hashMap)
	}
	return hashMap
}

func Run() int {
	// return utils.ProcessInput("day15.txt", 0, parseKeys, sumHash)
	hashMap := utils.ProcessInput("day15.txt", &HashMap{}, parseOperations, applyOperations)

	// hashMap := &HashMap{}
	// var line string
	// for n, _ := fmt.Scanln(&line); n > 0; n, _ = fmt.Scanln(&line) {
	// 	for _, op := range parseOperations(line) {
	// 		op.applyTo(hashMap)
	// 		fmt.Println(hashMap)
	// 	}
	// }

	fmt.Println(hashMap)
	return hashMap.totalFocusingPower()
}
