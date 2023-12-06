package arb_update

import (
	"sync"
)

// uses sync.Map or sync lock to make maps are safe for concurrent use
// var mutex = &sync.RWMutex{}

type Job struct {
	Entry       string
	Content     interface{}
}

type Result struct {
	WorkId      int
	KeepEntry   string
	KeepContent interface{}
}

type WorkerPool struct {
	Wg          *sync.WaitGroup
	Jobs        chan Job
	Results     chan Result
	Done        chan bool
	MaxWorkers  int
}

func NewWorkerPool(num int) *WorkerPool {
	wg := new(sync.WaitGroup)
	jobs := make(chan Job, num)
	results := make(chan Result, num)
	done := make(chan bool)
	return &WorkerPool {
		Wg:         wg,
		Jobs:       jobs,
		Results:    results,
		Done:       done,
		MaxWorkers: num,
	}
}

func (wp *WorkerPool) Allocate(templateEntries map[string]interface{}) {
	for k, v := range templateEntries {
		j := Job {
			Entry: k,
			Content: v,
		}
		wp.Jobs <- j
	}
	close(wp.Jobs)
}

func (wp *WorkerPool) Run(entries map[string]interface{}) {
	for n := 0; n < wp.MaxWorkers; n++ {
		wp.Wg.Add(1)
		go func(id int) {
			w := &Worker {
				Wg:      wp.Wg,
				Jobs:    wp.Jobs,
				Results: wp.Results,
				Entries: entries,
				Id:      id + 1,
			}
			w.Do()
		}(n)
	}
	wp.Wg.Wait()
	close(wp.Results)
}

func (wp *WorkerPool) ReadData() *sync.Map {
	entries := new(sync.Map)
	for {
		select {
		case val, ok := <- wp.Results:
			if !ok {
				wp.Done <- true
				return entries // break loop
			} else {
				//mutex.Lock()
				entries.Store(val.KeepEntry, val.KeepContent)
				// mutex.Unlock()
			}
		default:
		}
	}
}

type Worker struct {
	Wg       *sync.WaitGroup
	Jobs     <-chan Job
	Results  chan<- Result
	Entries  map[string]interface{} // locale entries
	Id       int
}

func (w *Worker) Do() {
	for j := range w.Jobs {
		check := CompareEntries(w.Entries, j.Entry)
		res := Result {
			WorkId:    w.Id,
			KeepEntry: j.Entry,
		}
		if check {
			res.KeepContent = w.Entries[j.Entry]
		} else {
			res.KeepContent = j.Content
		}
		w.Results <- res
	}
	w.Wg.Done()
}
