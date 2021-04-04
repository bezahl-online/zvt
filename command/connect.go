package command

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

var connecting bool = false

// Connect manages the connection to PT
func (p *PT) Connect() {
	if connecting {
		return
	}
	connecting = true
	var mutex *sync.RWMutex = (*p).lock
	if mutex == nil {
		mutex = &sync.RWMutex{}
	}
	var pause delay = delay{
		dur: 1 * time.Second,
	}

	if p.conn != nil {
		p.conn.Close()
	}

	// connect to PT via TCP/IP
	var err error = fmt.Errorf("no connection to PT")
	Logger.Debug("connecting", zap.String("url", os.Getenv("ZVT_URL")))
	for err != nil {
		err := p.Open()
		if err != nil {
			Logger.Error("error while connecting to PT",
				zap.Error(err))
			if pause.getSeconds() < 300 {
				pause.double()
			}
			pause.wait()
		} else {
			Logger.Debug("connection to PT established",
				zap.String("url", os.Getenv("ZVT_URL")),
			)
			break
		}
	}
	connecting = false
}

type delay struct {
	dur time.Duration
}

func (w *delay) getSeconds() int {
	return int((*w).dur.Seconds())
}

func (w *delay) wait() {
	time.Sleep(w.dur)
}

func (w *delay) double() {
	(*w).dur *= 2
}
