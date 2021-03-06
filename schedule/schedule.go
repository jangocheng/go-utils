package schedule

import (
	"log"
	"sync"
	"time"

	"github.com/robfig/cron"
)

type Schedule struct {
	Cron    *cron.Cron
	Stopped bool
	Wg      sync.WaitGroup
}

var (
	schedule *Schedule
	err      error
)

func init() {
	schedule, err = New()
	if err != nil {
		log.Fatalln(err)
	}
}

// New new cron
func New() (*Schedule, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	c := cron.NewWithLocation(loc)
	s := &Schedule{
		Cron:    c,
		Stopped: false,
		Wg:      sync.WaitGroup{},
	}

	s.Cron.Start()
	return s, nil
}

// AddCronFunc [docs](https://godoc.org/github.com/robfig/cron)
func AddCronFunc(spec string, fn func()) error {
	return schedule.AddCronFunc(spec, fn)
}

// AddLoopFunc can loop run task.
func AddLoopFunc(sleep time.Duration, fn func()) {
	schedule.AddLoopFunc(sleep, fn)
}

// Close close cron
func Close() {
	schedule.Close()
}

// AddCronFunc [docs](https://godoc.org/github.com/robfig/cron)
func (s *Schedule) AddCronFunc(spec string, fn func()) error {
	return s.Cron.AddFunc(spec, func() {
		s.Wg.Add(1)
		defer s.Wg.Done()

		fn()
	})
}

// AddLoopFunc can loop run task.
func (s *Schedule) AddLoopFunc(sleep time.Duration, fn func()) {
	go func() {
		for !s.Stopped {
			s.Wg.Add(1)
			fn()
			s.Wg.Done()

			time.Sleep(sleep)
		}
	}()
}

// Close close schedule
func (s *Schedule) Close() {
	s.Stopped = true
	s.Cron.Stop()
	s.Wg.Wait()
}
