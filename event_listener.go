package event_listener

import (
	"sync"
	"time"
)

type EventListener struct {
	condition sync.Cond
}

func NewEventListener() *EventListener {
	lock := &sync.Mutex{}
	return &EventListener{
		condition: sync.Cond{
			L: lock,
		},
	}
}

func (e *EventListener) Wait(shouldWait func() bool) {
	e.condition.L.Lock()
	for shouldWait() {
		e.condition.Wait()
	}
	e.condition.L.Unlock()
}

func (e *EventListener) WaitTimeout(shouldWait func() bool, timeout time.Duration) (didTimeOut bool) {
	timedOut := make(chan bool, 1)
	shouldWaitTimeout := func() bool {
		for {
			if !shouldWait() {
				return false
			}
			select {
			case <-timedOut:
				didTimeOut = true
				return false
			default:
				return true
			}
		}
	}
	go func() {
		timer := time.NewTimer(timeout)
		select {
		case <-timer.C:
			timedOut <- true
			e.notify()
		}
	}()
	e.Wait(shouldWaitTimeout)
	return didTimeOut
}

func (e *EventListener) WaitUntil(shouldWait func() bool, deadline time.Time) (didTimeOut bool) {
	timeout := deadline.Sub(time.Now())
	return e.WaitTimeout(shouldWait, timeout)
}

func (e *EventListener) notify() {
	e.condition.L.Lock()
	e.condition.Broadcast()
	e.condition.L.Unlock()
}
