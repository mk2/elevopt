package main

import "bytes"
import "time"
import "fmt"
import "math/rand"

type State string

const (
	MaxFloor                 = 14
	ElevatorCount            = 4
	RaisingProbability       = 10
	Moving             State = "MOVING"
	Stopped            State = "STOPPED"
	Arrived            State = "ARRIVED"
)

type Elev struct {
	pos   int
	dest  []bool
	state State
	num   int
}

func newElev() *Elev {
	e := new(Elev)
	e.dest = make([]bool, MaxFloor)
	e.state = Stopped
	return e
}

func (e *Elev) toChar() string {
	if e.state == Moving {
		return "M"
	} else if e.state == Stopped {
		return "S"
	} else if e.state == Arrived {
		return "A"
	}
	return "-"
}

func (e *Elev) tick(elevs []*Elev) State {
	e.state = Stopped
	for i := 0; i < MaxFloor; i++ {
		if e.dest[i] && i != e.pos {
			if i < e.pos {
				e.pos--
			} else {
				e.pos++
			}
			if i == e.pos {
				e.state = Arrived
				e.dest[i] = false
			} else {
				e.state = Moving
			}
		}
	}
	return e.state
}

func main() {
	// init elevators
	elevs := make([]*Elev, ElevatorCount)
	for i := 0; i < ElevatorCount; i++ {
		e := newElev()
		e.num = i
		elevs[i] = e
	}

	// init waiting peoples
	wps := make([]int, MaxFloor)

	sum := 0
	var outbuf bytes.Buffer
	for frame := 0; frame < 1000; frame++ {
		// gather waiting peoples
		for f := 0; f < MaxFloor; f++ {
			for _, e := range elevs {
				if rand.Int()%RaisingProbability == 0 {
					wps[f] += 1 + rand.Int()%5
					e.dest[f] = true
					sum += wps[f]
				}
			}
		}

		// check each elevator states
		for _, e := range elevs {
			if res := e.tick(elevs); res == Arrived || res == Stopped {
				wps[e.pos] = 0
			}
		}

		time.Sleep(1 * time.Millisecond)

		fmt.Println("")
		outbuf.Reset()
		outbuf.WriteString(fmt.Sprintf("-Frame:%d-Sum:%d---------------------\n", frame, sum))
		for f := 0; f < MaxFloor; f++ {
			outbuf.WriteString(fmt.Sprintf("%3dF", f))
			for _, e := range elevs {
				outbuf.WriteString("|")
				if e.pos == f {
					outbuf.WriteString(e.toChar())
				} else {
					if e.dest[f] {
						outbuf.WriteString("`")
					} else {
						outbuf.WriteString(" ")
					}
				}
			}
			outbuf.WriteString(fmt.Sprintf("|%3d", wps[f]))
			fmt.Println(outbuf.String())
			outbuf.Reset()
		}
		outbuf.Reset()
	}
	fmt.Printf("Total Waiting Peoples: %d\n", sum)
}
