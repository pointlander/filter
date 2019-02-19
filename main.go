package main

import (
	"math"
	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

const (
	// MemorySize is the size of the circular buffer holding the request times
	MemorySize = 256
	// Threshold is the ratio above which there is considered to be a spike in traffic
	Threshold = 2
)

func main() {
	rand.Seed(1)

	// stores the time of previous requests in a ring buffer
	index, prev, size, memory := 0, MemorySize-1, 0, make([]int, MemorySize)
	add := func(t int) (oldest, previous int) {
		memory[index], previous = t, memory[prev]
		index, prev = (index+1)%MemorySize, index
		if size < MemorySize {
			size++
		}
		oldest = memory[index]
		return
	}

	// low pass filter for the request rate; (rate / filtered rate) is returned
	// https://en.wikipedia.org/wiki/Low-pass_filter#Discrete-time_realization
	filter := 0.0
	lowpass := func(time int) float64 {
		oldest, previous := add(time)
		alpha := float64(time-previous) / float64(time-oldest)
		rate := float64(size) / float64(time-oldest)
		filter = alpha*rate + (1-alpha)*filter
		return rate / filter
	}

	lastSpike, lastRatio := 0, 0.0

	// simulator
	time, points, probabilities := 0, make(plotter.XYs, 0, 2*MemorySize*1024), make(plotter.XYs, 0, 2*MemorySize*1024)
	simulate := func(requests, delta int) (spikes int) {
		for i := 0; i < requests; i++ {
			time += rand.Intn(delta) + 1
			ratio := lowpass(time)
			if ratio > Threshold {
				lastSpike, lastRatio = time, ratio-1
				spikes++
			}
			points = append(points, plotter.XY{X: float64(time), Y: ratio})
			probability := 1 / (1 + lastRatio*math.Exp(.00001*float64(lastSpike-time)))
			probabilities = append(probabilities, plotter.XY{X: float64(time), Y: probability})
		}
		return
	}

	// a scenario is simulated below

	// base load
	if simulate(MemorySize*1024, 8) > 0 {
		panic("there should be no spike")
	}

	// burst of traffic
	if simulate(1024, 2) == 0 {
		panic("there should have been a spike")
	}

	// base load resumes
	simulate(MemorySize, 8)
	if simulate(MemorySize*1024, 8) > 0 {
		panic("there should be no spike")
	}

	// traffic stops for long time
	time += 1024

	// base load resumes
	if simulate(MemorySize*1024, 8) > 0 {
		panic("there should be no spike")
	}

	// a moderate burst of traffic
	if simulate(MemorySize*1024, 4) > 0 {
		panic("there should be no spike")
	}

	// base load resumes
	if simulate(MemorySize*1024, 8) > 0 {
		panic("there should be no spike")
	}

	// burst of traffic
	if simulate(1024, 2) == 0 {
		panic("there should have been a spike")
	}

	// base load resumes
	if simulate(MemorySize*1024, 8) > 0 {
		panic("there should be no spike")
	}

	// graph the ratio results
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "rate / filtered rate"
	p.X.Label.Text = "time"
	p.Y.Label.Text = "ratio"

	s, err := plotter.NewScatter(points)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Radius = vg.Length(1)
	s.GlyphStyle.Shape = draw.CircleGlyph{}
	p.Add(s)

	err = p.Save(8*vg.Inch, 8*vg.Inch, "points.png")
	if err != nil {
		panic(err)
	}

	// graph the probability results
	p, err = plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "don't drop probabiltiy"
	p.X.Label.Text = "time"
	p.Y.Label.Text = "probabiltiy"

	s, err = plotter.NewScatter(probabilities)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Radius = vg.Length(1)
	s.GlyphStyle.Shape = draw.CircleGlyph{}
	p.Add(s)

	err = p.Save(8*vg.Inch, 8*vg.Inch, "probabilities.png")
	if err != nil {
		panic(err)
	}
}
