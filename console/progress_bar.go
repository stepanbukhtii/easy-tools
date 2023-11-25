package console

import (
	"fmt"
	"time"
)

type ProgressBar struct {
	currentIndex int
	totalIndex   int
	lastTimeWork time.Time
}

func NewProgressBar(totalIndex int) ProgressBar {
	return ProgressBar{
		totalIndex:   totalIndex,
		lastTimeWork: time.Now(),
	}
}

func (p *ProgressBar) UpdateAndPrint() {
	percent := 100 * float64(p.currentIndex) / float64(p.totalIndex)
	remainingTime := time.Now().Sub(p.lastTimeWork) * time.Duration(p.totalIndex-(p.currentIndex))

	fmt.Printf("\rProgress %d%% %d/%d %s", int64(percent), p.currentIndex, p.totalIndex, remainingTime.String())

	p.lastTimeWork = time.Now()
	p.currentIndex++
}

func (p *ProgressBar) ClearConsole() {
	fmt.Printf("\r%c[2K", 27)
}
