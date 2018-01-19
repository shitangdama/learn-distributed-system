单次执行once
线程池pool
互斥锁sync.Mutex lock unlock
### 读写锁rwmutex

// Lock 将 rw 设置为写锁定状态，禁止其他例程读取或写入。
func (rw *RWMutex) Lock()

// Unlock 解除 rw 的写锁定状态，如果 rw 未被写锁定，则该操作会引发 panic。
func (rw *RWMutex) Unlock()

// RLock 将 rw 设置为读锁定状态，禁止其他例程写入，但可以读取。
func (rw *RWMutex) RLock()

// Runlock 解除 rw 的读锁定状态，如果 rw 未被读锁顶，则该操作会引发 panic。
func (rw *RWMutex) RUnlock()

// RLocker 返回一个互斥锁，将 rw.RLock 和 rw.RUnlock 封装成了一个 Locker 接口。
func (rw *RWMutex) RLocker() Locker

### 条件等待cond
// 创建一个条件等待
func NewCond(l Locker) *Cond

// Broadcast 唤醒所有等待的 Wait，建议在“更改条件”时锁定 c.L，更改完毕再解锁。
func (c *Cond) Broadcast()

// Signal 唤醒一个等待的 Wait，建议在“更改条件”时锁定 c.L，更改完毕再解锁。
func (c *Cond) Signal()

// Wait 会解锁 c.L 并进入等待状态，在被唤醒时，会重新锁定 c.L
func (c *Cond) Wait()
### 组等待WaitGroup

WaitGroup 用于等待一组例程的结束。主例程在创建每个子例程的时候先调用 Add 增加等待计数，每个子例程在结束时调用 Done 减少例程计数。之后，主例程通过 Wait 方法开始等待，直到计数器归零才继续执行。

------------------------------

// 计数器增加 delta，delta 可以是负数。
func (wg *WaitGroup) Add(delta int)

// 计数器减少 1
func (wg *WaitGroup) Done()

// 等待直到计数器归零。如果计数器小于 0，则该操作会引发 panic。
func (wg *WaitGroup) Wait()