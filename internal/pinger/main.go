package pinger

import (
	"context"
	"sync"
	"time"

	ping "github.com/go-ping/ping"
	"github.com/sirupsen/logrus"
)

type pingerIPStatus struct {
	ip string
	ok bool
}

func (s *pingerIPStatus) Ok() bool {
	return s.ok
}

type Pinger struct {
	logger *logrus.Entry

	ips         []string
	interval    time.Duration
	status      map[string]pingerIPStatus
	statusMutex sync.RWMutex
	statusChan  chan pingerIPStatus
	ok          bool

	stop     chan struct{}
	stopOnce sync.Once
}

func (p *Pinger) ping(addr string) {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	status := pingerIPStatus{ip: addr, ok: true}
	err = pinger.Run() // blocks
	if err != nil {
		status.ok = false
	}
	if pinger.Statistics().PacketLoss == 100 {
		status.ok = false
	}
	p.statusChan <- status
}

func (p *Pinger) check() {
	for _, ip := range p.ips {
		go p.ping(ip)
	}
}

func (p *Pinger) Start(ctx context.Context) {
	p.check()

	go p.process(ctx)
	go p.updater(ctx)
}

func (p *Pinger) process(ctx context.Context) {
	timer := time.NewTicker(p.interval)

	for {
		select {
		case <-ctx.Done():
			return

		case <-p.stop:
			return

		case <-timer.C:
			p.check()
		}
	}
}

func (p *Pinger) updater(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-p.stop:
			return

		case status := <-p.statusChan:
			p.statusMutex.Lock()
			p.status[status.ip] = status
			p.statusMutex.Unlock()
			p.updateStatus()
		}
	}
}

func (p *Pinger) Stop() {
	p.stopOnce.Do(func() {
		close(p.stop)
	})
}

func (p *Pinger) Ok() bool {
	return p.ok
}

func (p *Pinger) updateStatus() {
	defer p.statusMutex.RUnlock()

	fails := 0
	p.statusMutex.RLock()
	for _, s := range p.status {
		if !s.Ok() {
			fails++
		}
	}

	p.ok = fails != len(p.ips)
}

func NewPinger(logger *logrus.Entry, ips []string, interval int) *Pinger {
	return &Pinger{
		logger: logger.WithField("from", "pinger"),

		ips:         ips,
		interval:    time.Duration(interval) * time.Second,
		status:      make(map[string]pingerIPStatus, len(ips)),
		statusMutex: sync.RWMutex{},
		statusChan:  make(chan pingerIPStatus),

		stop:     make(chan struct{}),
		stopOnce: sync.Once{},
	}
}
