package event_listener

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventListener_Wait(t *testing.T) {
	event := NewEventListener()
	startTime := time.Now()
	shouldWait := true

	go func() {
		timer := time.NewTimer(time.Second / 2)
		select {
		case <-timer.C:
			shouldWait = false
			event.notify()
		}
	}()

	event.Wait(func() bool {
		return shouldWait
	})

	elapsed := time.Now().Sub(startTime)
	assert.GreaterOrEqual(t, elapsed, time.Second/2)
}

func TestEventListener_WaitTimeout_Timeout(t *testing.T) {
	event := NewEventListener()
	startTime := time.Now()

	didTimeOut := event.WaitTimeout(func() bool {
		return true
	}, time.Second)

	elapsed := time.Now().Sub(startTime)

	assert.True(t, didTimeOut)
	assert.Less(t, elapsed, time.Second*2)
}

func TestEventListener_WaitTimeout_No_Timeout(t *testing.T) {
	event := NewEventListener()
	startTime := time.Now()

	shouldWait := true

	go func() {
		timer := time.NewTimer(time.Second)
		select {
		case <-timer.C:
			shouldWait = false
			event.notify()
		}
	}()

	didTimeOut := event.WaitTimeout(func() bool {
		return shouldWait
	}, time.Second*2)

	elapsed := time.Now().Sub(startTime)

	assert.False(t, didTimeOut)
	assert.Less(t, elapsed, time.Second*2)
}

func TestEventListener_WaitUntil_Timeout(t *testing.T) {
	event := NewEventListener()
	startTime := time.Now()

	didTimeOut := event.WaitUntil(func() bool {
		return true
	}, time.Now().Add(time.Second))

	elapsed := time.Now().Sub(startTime)

	assert.True(t, didTimeOut)
	assert.Less(t, elapsed, time.Second*2)
}

func TestEventListener_WaitUntil_No_Timeout(t *testing.T) {
	event := NewEventListener()
	startTime := time.Now()
	shouldWait := true

	go func() {
		timer := time.NewTimer(time.Second)
		select {
		case <-timer.C:
			shouldWait = false
			event.notify()
		}
	}()

	didTimeOut := event.WaitUntil(func() bool {
		return shouldWait
	}, time.Now().Add(time.Second))

	elapsed := time.Now().Sub(startTime)

	assert.False(t, didTimeOut)
	assert.Less(t, elapsed, time.Second*2)
}
