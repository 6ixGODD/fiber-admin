package cron

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	cron   *cron.Cron
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}

func New(ctx context.Context) *Cron {
	c, cancel := context.WithCancel(ctx)
	return &Cron{
		cron:   cron.New(),
		ctx:    c,
		cancel: cancel,
	}
}

func (c *Cron) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cron.Start()
}

func (c *Cron) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cron.Stop()
	c.cancel()
}

func (c *Cron) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	return c.cron.AddFunc(spec, cmd)
}

func (c *Cron) AddJob(spec string, cmd cron.Job) (cron.EntryID, error) {
	return c.cron.AddJob(spec, cmd)
}

func (c *Cron) Remove(id cron.EntryID) {
	c.cron.Remove(id)
}

func (c *Cron) Entries() []cron.Entry {
	return c.cron.Entries()
}

func (c *Cron) Context() context.Context {
	return c.ctx
}
