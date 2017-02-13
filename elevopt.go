package main

import "bytes"
import "time"
import "fmt"

import ui "github.com/gizak/termui"

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
}

func newElev() Elev {
	e := Elev{}
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

func (e *Elev) tick() State {
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
			} else {
				e.state = Moving
			}
		}
	}
	return e.state
}

func main() {
	// init ui
	err := ui.Init()
	if err != nil {
		panic(err)
	}

	defer ui.Close()

	// add q key handler
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	// init elevators
	elevs := make([]Elev, ElevatorCount)
	for i := 0; i < ElevatorCount; i++ {
		elevs[i] = newElev()
	}
	// init waiting peoples
	wps := make([]int, MaxFloor)

	ui.Loop()

	var outbuf bytes.Buffer
	frame := 0
	for frame = 0; frame < 100; frame++ {
		outbuf.Reset()
		// check each elevator states
		for _, e := range elevs {
			if res := e.tick(); res == Arrived || res == Stopped {
				wps[e.pos] = 0
			}
		}

		time.Sleep(500 * time.Millisecond)

		for f := 0; f < MaxFloor; f++ {
			outbuf.WriteString("-------------------------\n")
			for _, e := range elevs {
				if e.pos == f {
					outbuf.WriteString(e.toChar())
				} else {
					if e.dest[f] {
						outbuf.WriteString("`")
					} else {
						outbuf.WriteString(" ")
					}
				}
				outbuf.WriteString(fmt.Sprintf("|%d", wps[f]))
			}
			outbuf.WriteString("\n")
		}
		par := ui.NewPar(outbuf.String())
		par.Height = 100
		par.Width = 100
		ui.Render(par)
	}
}
