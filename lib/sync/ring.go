package sync

import (
	"sync"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGate(new(GoDialRing))
	RegisterEvalGate(new(GoRingLease))
}

type GoDialRing struct {
	Cap int64 `ko:"name=cap"`
}

type GoRingConn interface {
	GoRingConn()
	Lease() int64 // returns the sequential number of the lease
}

func (g *GoDialRing) Play(ctx *runtime.Context) GoRingConn {
	return newGoSyncRing(int(g.Cap))
}

type GoRingLease struct {
	Conn GoRingConn `ko:"name=conn,monadic"`
}

func (g *GoRingLease) Play(ctx *runtime.Context) int64 {
	return g.Conn.Lease()
}

// implementation

type goSyncRing struct {
	lock sync.Mutex
	next int64
	ring chan int64
}

func newGoSyncRing(cap int) *goSyncRing {
	gsr := &goSyncRing{
		next: 0,
		ring: make(chan int64, cap),
	}
	for i := 0; i < cap; i++ {
		gsr.refill()
	}
	return gsr
}

func (gsr *goSyncRing) GoRingConn() {}

func (gsr *goSyncRing) refill() {
	gsr.lock.Lock()
	defer gsr.lock.Unlock()
	gsr.ring <- gsr.next
	gsr.next++
}

func (gsr *goSyncRing) Lease() int64 {
	lease := <-gsr.ring
	gsr.refill()
	return lease
}
