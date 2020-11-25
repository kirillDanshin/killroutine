package killroutine

import (
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	sup := newSupervisor()

	go sup.queueFunc(func() {
		println("hey")
		time.Sleep(10 * time.Second)
		panic("oops")
	}, 2*time.Second)

	time.Sleep(15 * time.Second)
}
