# lsync

Local syncing package with support for timeouts. This package offers both a `sync.Mutex` and `sync.RWMutex` compatible interface.

Additionally it provides `lsync.LFrequentAccess` which uses an atomic load and store of a consistently typed value. This can be usefull for shared data structures that are frequently read but infrequently updated (using an copy-on-write mechanism) without the need for protection with a regular mutex.

### Example of LRWMutex

```go
	// Create RWMutex compatible mutex 
	lrwm := NewLRWMutex()

	// Try to get lock within timeout 
	if !lrwm.GetLock(1000 * time.Millisecond) {
		fmt.Println("Timeout occured")
		return
	}
	
	// Acquired lock, do your stuff ...

	lrwm.Unlock() // Release lock
```

### Example of LFrequentAccess
````go
	type Map map[string]string

	// Create new LFrequentAccess for type Map
	freqaccess := NewLFrequentAccess(make(Map))

	cur := freqaccess.LockBeforeSet().(Map) // Lock in order to update
	mp := make(Map)                         // Create new Map
	for k, v := range cur {                 // Copy over old contents
		mp[k] = v
	}
	mp[key] = val                      // Add new value
	freqaccess.SetNewCopyAndUnlock(mp) // Exchange old version of map with new version

	mpReadOnly := freqaccess.ReadOnlyAccess().(Map) // Get read only access to Map
	fmt.Println(mpReadOnly[key])                    // Safe access with no further synchronization
````

## API

#### LMutex

```go
func (lm *LMutex) Lock()
func (lm *LMutex) GetLock(timeout time.Duration) (locked bool)
func (lm *LMutex) Unlock()
```

#### LRWMutex

```go
func (lm *LRWMutex) Lock()
func (lm *LRWMutex) GetLock(timeout time.Duration) (locked bool)
func (lm *LRWMutex) RLock()
func (lm *LRWMutex) GetRLock(timeout time.Duration) (locked bool)
func (lm *LRWMutex) Unlock()
func (lm *LRWMutex) RUnlock()
```

#### LFrequentAccess 
```go
func (lm *LFrequentAccess) ReadOnlyAccess() (constReadOnly interface{})
func (lm *LFrequentAccess) LockBeforeSet() (constCurVersion interface{})
func (lm *LFrequentAccess) SetNewCopyAndUnlock(newCopy interface{})
```

## Benchmarks
### LFrequentAccess

An `lsync.LFrequentAccess` provides an atomic load and store of a consistently typed value.

```
benchmark                           old ns/op     new ns/op     delta
BenchmarkLFrequentAccessMap-8       114           4.67          -95.90%
BenchmarkLFrequentAccessSlice-8     109           5.95          -94.54%
```

### sync.Mutex vs lsync.LMutex 

(with `defaultRetryUnit` and `defaultRetryCap` at 1 ms)
```
benchmark                          old ns/op     new ns/op     delta
BenchmarkMutex-8                   111           2649          +2286.49%
BenchmarkMutexSlack-8              120           1719          +1332.50%
BenchmarkMutexWork-8               133           2637          +1882.71%
BenchmarkMutexWorkSlack-8          137           1729          +1162.04%
```

(with `defaultRetryUnit` and `defaultRetryCap` at 100 ms)

```
benchmark                          old ns/op     new ns/op     delta
BenchmarkMutex-8                   111           2649          +2286.49%
BenchmarkMutexSlack-8              120           2478          +1965.00%
BenchmarkMutexWork-8               133           2547          +1815.04%
BenchmarkMutexWorkSlack-8          137           2683          +1858.39%
```
