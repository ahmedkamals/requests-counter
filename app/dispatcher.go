package app

import (
	"io"
	"log"
	"runtime"
	"sync"
	"time"
)

type (
	Dispatcher struct {
		requestChan    chan uint64
		responseChan   chan uint64
		oneSecTicker   *time.Ticker
		milliSecTicker *time.Ticker
		indexOneSec    uint8
		counter        [60]uint64
		saver          *Saver
	}
)

func NewDispatcher(saver *Saver) *Dispatcher {
	return &Dispatcher{
		requestChan:    make(chan uint64, 1000),
		responseChan:   make(chan uint64),
		oneSecTicker:   time.NewTicker(time.Second),
		milliSecTicker: time.NewTicker(time.Millisecond),
		counter:        [60]uint64{},
		saver:          saver,
	}
}

func (d *Dispatcher) Run() *Dispatcher {
	runtime.GOMAXPROCS(runtime.NumCPU())

	data, err := d.saver.Read()

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	recovery := NewEmptyRecovery()
	recovery.Unmarshal(data)

	mutex := sync.Mutex{}
	mutex.Lock()
	d.restoreCounter(recovery)
	mutex.Unlock()

	go d.dispatch()
	go d.heartbeat()

	return d
}

func (d *Dispatcher) dispatch() {
	mutex := sync.Mutex{}

	for request := range d.requestChan {
		mutex.Lock()
		d.counter[d.indexOneSec] += request
		d.updateResponse()
		mutex.Unlock()
	}
}

func (d *Dispatcher) restoreCounter(recovery *Recovery) {
	if recovery != nil {
		diff := uint64(time.Now().Unix()) - recovery.TimeStamp

		// Restore only the parts we are interested in.
		if diff < 60 {
			d.indexOneSec = recovery.Index
			rangeStart := uint64(d.indexOneSec)
			rangeEnd := (diff + rangeStart) % 60

			for i := uint64(0); i < 60; i++ {
				// Check if overlapping period
				// # means a cancelled period
				// [#####11111111111111#####]
				if rangeEnd < rangeStart {
					if i >= rangeStart || i <= rangeEnd {
						continue
					}
				} else {
					if i >= rangeStart && i <= rangeEnd {
						continue
					}
				}

				d.counter[i] = recovery.Data[i]
			}
		}
	}
}

func (d *Dispatcher) updateResponse() {
	count := uint64(0)

	for _, val := range d.counter {
		count += val
	}

	d.responseChan <- count
}

func (d *Dispatcher) heartbeat() {
	mutex := sync.Mutex{}

	for {
		select {
		// Resetting the value every second.
		case <-d.oneSecTicker.C:
			mutex.Lock()
			d.indexOneSec++
			d.indexOneSec %= 60
			d.counter[d.indexOneSec] = 0
			mutex.Unlock()

		case <-d.milliSecTicker.C:
			r := NewRecovery(d.indexOneSec, uint64(time.Now().Unix()), d.counter)
			data, err := r.Marshall()

			if err != nil {
				log.Fatal(err)
			}

			d.saver.Save(data)
		}
	}
}

func (d *Dispatcher) Stop() {
	d.saver.Close()
}
