package killroutine

import (
	"sync"
	"time"
)

type supervisor struct {
	mu       sync.RWMutex
	timerIDs map[int]*g
	nextID   int
}

func newSupervisor() *supervisor {
	return &supervisor{
		timerIDs: make(map[int]*g),
	}
}

func (s *supervisor) acquireTimerID(g *g) (id int) {
	s.mu.Lock()
	s.timerIDs[s.nextID] = g
	id = s.nextID
	s.nextID++
	s.mu.Unlock()
	return
}

func (s *supervisor) releaseTimerID(id int) {
	s.mu.Lock()
	delete(s.timerIDs, id)
	s.mu.Unlock()
}

func (s *supervisor) run(gpCh chan *g, done chan struct{}, f func()) {
	gpCh <- getg()

	f()
	done <- struct{}{}
}

func (s *supervisor) queueFunc(f func(), timeout time.Duration) {
	gpCh := make(chan *g)
	done := make(chan struct{})

	go s.run(gpCh, done, f)

	timerID := s.acquireTimerID(<-gpCh)

	select {
	case <-time.After(timeout):
		s.kill(timerID)
	case <-done:
	}
}

func (s *supervisor) kill(timerID int) {
	s.mu.RLock()
	gp := s.timerIDs[timerID]
	s.mu.RUnlock()

	go systemstack(func() {
		userG := ((*g)(getg())).m.curg
		if readgstatus(userG) == _Grunning {
			casgstatustyped(userG, _Grunning, _Gwaiting)
			(*g)(gp).waitreason = 8
		}

		runtime_suspendG(gp)
		goexit0(gp)
	})

	s.releaseTimerID(timerID)

	println("kill completed")
}
