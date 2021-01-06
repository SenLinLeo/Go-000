package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	BucketOptNum = 3
	BucketSize   = 10
	BucketWidth  = time.Second
	UnlockState  = 0
	LockState    = 1
)

const (
	BucketOptSuc = iota
	BucketOptFail
	BucketOpt
)

type bucket struct {
	startTime int64
	optList   []int64
}

func NewBucket(startTime int64) *bucket {
	return &bucket{
		startTime: startTime,
		optList:   make([]int64, BucketOptNum),
	}
}

type rollingNumber struct {
	window     []*bucket
	size       int
	width      int64
	head       int
	tail       int
	lockValue  int32
	lockWindow sync.RWMutex
}

func NewRollingNumber(size int, width int64) *rollingNumber {
	if size <= 1 {
		size = BucketSize
	}

	if width <= 0 {
		width = int64(BucketWidth)
	}
	return &rollingNumber{
		window:    make([]*bucket, size),
		size:      size,
		width:     width,
		head:      0,
		tail:      0,
		lockValue: UnlockState,
	}
}

func (r *rollingNumber) reset() {
	r.lockWindow.RLock()
	defer r.lockWindow.RUnlock()
	r.window = make([]*bucket, r.size)
	r.head = 0
	r.tail = 0
}

func (r *rollingNumber) getLastBucket() *bucket {
	r.lockWindow.RLock()
	defer r.lockWindow.RUnlock()
	if r.head == r.tail {
		return nil
	}
	return r.window[(r.tail-1+r.size)%r.size]
}

func (r *rollingNumber) addBucket(bucket *bucket) {
	r.lockWindow.RLock()
	defer r.lockWindow.RUnlock()
	r.window[r.tail] = bucket
	if (r.tail+1)%r.size == r.head {
		r.head = (r.head + 1) % r.size
	}
	r.tail = (r.tail + 1) % r.size
}

func (r *rollingNumber) GetCurrentBucket() *bucket {
	currentTime := time.Now().UnixNano()

	last := r.getLastBucket()
	if last != nil && currentTime < last.startTime+r.width {
		return last
	}

	if atomic.CompareAndSwapInt32(&r.lockValue, UnlockState, LockState) {
		if r.getLastBucket() == nil {
			newBucket := NewBucket(currentTime)
			r.addBucket(newBucket)
			atomic.StoreInt32(&r.lockValue, UnlockState)
			return newBucket
		} else {
			for i := 0; i < r.size; i++ {
				last = r.getLastBucket()
				if currentTime < last.startTime+r.width {
					atomic.StoreInt32(&r.lockValue, UnlockState)
					return last
				} else if currentTime-(last.startTime+r.width) > int64(r.size)*r.width {
					r.reset()
					atomic.StoreInt32(&r.lockValue, UnlockState)
					return r.GetCurrentBucket()
				} else {
					newBucket := NewBucket(currentTime + r.width)
					r.addBucket(newBucket)
					atomic.StoreInt32(&r.lockValue, UnlockState)
				}
			}
			return r.getLastBucket()
		}
	} else {
		currentBucket := r.getLastBucket()
		if currentBucket != nil {
			return currentBucket
		} else {
			time.Sleep(5)
			return r.GetCurrentBucket()
		}
	}
}

func (r *rollingNumber) AddBulletCell(index int, num int64) {
	if index >= BucketOptNum {
		return
	}
	tmpBucket := r.GetCurrentBucket()
	atomic.AddInt64(&tmpBucket.optList[index], num)
}

func (r *rollingNumber) GetSumByIndex(index int) int64 {
	sum := int64(0)
	for _, v := range r.window {
		if v == nil {
			continue
		}
		sum += v.optList[index]
	}
	return sum
}

func (r *rollingNumber) GetSum() []int64 {
	sumList := make([]int64, BucketOptNum)
	for _, v := range r.window {
		if v == nil {
			continue
		}
		for i := 0; i < BucketOptNum; i++ {
			sumList[i] += v.optList[i]
		}
	}
	return sumList
}

func main() {
	rn := NewRollingNumber(10, int64(time.Second))
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			rn.AddBulletCell(BucketOptSuc, 1)
			rn.AddBulletCell(BucketOptFail, 2)
			rn.AddBulletCell(BucketOpt, 1)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			rn.AddBulletCell(BucketOptSuc, 1)
			rn.AddBulletCell(BucketOptFail, 1)
			rn.AddBulletCell(BucketOpt, 1)
		}
	}()
	wg.Wait()
	fmt.Println(rn.GetSum())
}
