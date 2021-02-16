package gocron

import (
	"time"

	"github.com/gorhill/cronexpr"
)

type Task struct {
	job  func()
	expr *cronexpr.Expression
	from time.Time
	curr time.Time
}

type Cron struct {
	tasks   *SafeSlice
	channel chan func()
}

var cron Cron

func NewCron() *Cron {

	cron = Cron{
		tasks:   NewSafeSlice(0),
		channel: make(chan func()),
	}

	for i := 0; i < 100; i++ {
		go doWork(cron.channel)
	}

	go cron.scheduler(&cron)

	return &cron
}

func (c *Cron) scheduler(cron *Cron) {
	ticker := time.NewTicker(time.Second)

	for {

		select {
		case <-ticker.C:
			for _, t := range cron.tasks.slice {
				if t.curr.Sub(time.Now()) < time.Nanosecond {
					t.from = t.curr
					t.curr = t.expr.Next(t.from)
					cron.channel <- t.job
				}
			}

		}

	}
}

// MustRun Handles Errors for Run
func (c *Cron) MustRun(str string, job func()) *Task {
	task, err := c.Run(str, job)

	if err != nil {
		panic("Couldn't create the task")
	}

	return task
}

// Run a job (function) as a cron job
func (c *Cron) Run(str string, job func()) (*Task, error) {
	expr, err := cronexpr.Parse(str)
	if err != nil {
		panic(err)
	}

	from := time.Now()

	task := Task{
		job:  job,
		expr: expr,
		from: from,
		curr: expr.Next(from),
	}

	cron.tasks.Add(&task)

	return &task, nil
}

func doWork(ch <-chan func()) {
	for {
		job := <-ch
		job()
	}
}
