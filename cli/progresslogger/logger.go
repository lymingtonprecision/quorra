package progresslogger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/vmware/govmomi/vim25/progress"
)

type ProgressLogger struct {
	Prefix string

	wg sync.WaitGroup

	sink chan chan progress.Report
	done chan struct{}
}

func NewProgressLogger(prefix string) *ProgressLogger {
	pl := &ProgressLogger{
		Prefix: prefix,
		sink:   make(chan chan progress.Report),
		done:   make(chan struct{}),
	}

	pl.wg.Add(1)
	go pl.waitForSinkCall()

	return pl
}

func (pl *ProgressLogger) Log(s string, replaceLine bool) (int, error) {
	if len(s) > 0 && replaceLine {
		os.Stdout.Write([]byte{'\r', 033, '[', 'K'})
		os.Stdout.Sync()
	}

	n, err := os.Stdout.Write([]byte(s))
	os.Stdout.Sync()
	return n, err
}

/* waitForSinkCall waits for a call to the ProgressLogger's progress
 * report sink channel and then passes processing over to
 * processReports.
 *
 * Prints any errors raised during processing or OK if the sink was
 * processed successfully.
 */
func (pl *ProgressLogger) waitForSinkCall() {
	var err error

	defer pl.wg.Done()

	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()

	for stop := false; !stop; {
		select {
		case ch := <-pl.sink:
			err = pl.processReports(tick, ch)
			stop = true
		case <-pl.done:
			stop = true
		case <-tick.C:
			pl.Log(pl.Prefix, true)
		}
	}

	if err != nil && err != io.EOF {
		pl.Log(fmt.Sprintf("%s Error: %s\n", pl.Prefix, err), true)
	} else {
		pl.Log(fmt.Sprintf("%s OK\n", pl.Prefix), true)
	}
}

/* Process progress reports from the given channel until, progressively
 * printing the progress completed per `tick`
 */
func (pl *ProgressLogger) processReports(tick *time.Ticker, ch <-chan progress.Report) error {
	var r progress.Report
	var ok bool
	var err error

	for ok = true; ok; {
		select {
		case r, ok = <-ch:
			if !ok {
				break
			}
			err = r.Error()
		case <-tick.C:
			if r != nil {
				line := fmt.Sprintf(" (%.0f%%", r.Percentage())
				detail := r.Detail()
				if detail != "" {
					line += fmt.Sprintf(", %s", detail)
				}
				line += ")"

				pl.Log(pl.Prefix+line, true)
			} else {
				pl.Log(pl.Prefix, true)
			}
		}
	}

	return err
}

func (pl *ProgressLogger) Sink() chan<- progress.Report {
	ch := make(chan progress.Report)
	pl.sink <- ch
	return ch
}

func (pl *ProgressLogger) Wait() {
	close(pl.done)
	pl.wg.Wait()
}
