package timer

import (
	"sync"
	"time"
)

// Timer 表示一个计时器
type Timer struct {
	ID       uint32
	Interval time.Duration
	Callback func()
	timer    *time.Timer
	running  bool
}

// TimerManager 管理多个计时器
type TimerManager struct {
	timers map[uint32]*Timer
	mu     sync.Mutex
	nextID uint32
}

// NewTimerManager 创建一个新的计时器管理器
func NewTimerManager() *TimerManager {
	return &TimerManager{
		timers: make(map[uint32]*Timer),
	}
}

// AddTimer 添加一个新的计时器
func (tm *TimerManager) AddTimer(interval time.Duration, callback func()) uint32 {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.nextID++
	id := tm.nextID

	timer := &Timer{
		ID:       id,
		Interval: interval,
		Callback: callback,
		running:  true,
	}

	tm.timers[id] = timer

	// 创建并启动计时器
	timer.timer = time.NewTimer(interval)
	go func() {
		for timer.running {
			<-timer.timer.C
			if timer.running {
				timer.Callback()
				timer.timer.Reset(interval)
			}
		}
	}()

	return id
}

// RemoveTimer 移除一个计时器
func (tm *TimerManager) RemoveTimer(id uint32) bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if timer, exists := tm.timers[id]; exists {
		timer.running = false
		timer.timer.Stop()
		delete(tm.timers, id)
		return true
	}
	return false
}

// Clear 清除所有计时器
func (tm *TimerManager) Clear() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, timer := range tm.timers {
		timer.running = false
		timer.timer.Stop()
	}
	tm.timers = make(map[uint32]*Timer)
}

func Test() {
	// 创建并使用计时器管理器
	timerMgr := NewTimerManager()
	timerID := timerMgr.AddTimer(time.Second*5, func() {

	})
	defer timerMgr.RemoveTimer(timerID)
}

// GetLocalTime 获取时间
func GetLocalTime() time.Time {
	return time.Now()
}
