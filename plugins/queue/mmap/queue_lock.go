// SiGG-Satellite-Network-SII  //

//go:build !windows

package mmap

import "sync/atomic"

func (q *Queue) lock(segmentID int64) {
	index := q.GetIndex(segmentID)
	q.lockByIndex(index)
}

func (q *Queue) unlock(segmentID int64) {
	index := q.GetIndex(segmentID)
	q.unlockByIndex(index)
}

func (q *Queue) lockByIndex(index int) {
	for !atomic.CompareAndSwapInt32(&q.locker[index], 0, 1) {
	}
}

func (q *Queue) unlockByIndex(index int) {
	for !atomic.CompareAndSwapInt32(&q.locker[index], 1, 0) {
	}
}
