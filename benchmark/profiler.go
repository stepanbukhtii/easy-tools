package benchmark

import (
	"github.com/stepanbukhtii/easy-tools/stringx/pretty"
	"runtime"
	"time"
)

type Profiler struct {
	m1, m2    runtime.MemStats
	startTime time.Time
	elapsed   time.Duration
}

func (c *Profiler) Start() {
	runtime.ReadMemStats(&c.m1)
	c.startTime = time.Now()
}

func (c *Profiler) Finish() {
	c.elapsed = time.Since(c.startTime)
	runtime.ReadMemStats(&c.m2)
}

func (c *Profiler) Size() string {
	totalSize := int64(c.m2.TotalAlloc - c.m1.TotalAlloc)
	return pretty.BytesToSize(totalSize)
}

func (c *Profiler) Time() string {
	return c.elapsed.String()
}
