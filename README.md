# lsync

Local syncing package with support for timeouts.

## Usage

```go
	// Try to get lock within timeout 
	if !lrwm.GetLock(1000 * time.Millisecond) {
		fmt.Println("Timeout occured")
		return
	}
	
	// Acquired lock, do your stuff ...

	lrwm.Unlock() // Release lock
```

## API

```go
func (lm *LRWMutex) Lock()
```
```go
func (lm *LRWMutex) GetLock(timeout time.Duration) (locked bool)
```
```go
func (lm *LRWMutex) RLock()
```
```go
func (lm *LRWMutex) GetRLock(timeout time.Duration) (locked bool)
```
```go
func (lm *LRWMutex) Unlock()
```
```go
func (lm *LRWMutex) RUnlock()
```

### Frequently Accessed Shared Values

An `lsync.LFrequentAccess` provides an atomic load and store of a consistently typed value.

```
benchmark                           old ns/op     new ns/op     delta
BenchmarkLFrequentAccessMap-8       114           4.67          -95.90%
BenchmarkLFrequentAccessSlice-8     109           5.95          -94.54%
```


