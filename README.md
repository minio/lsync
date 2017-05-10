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
func (lm *LRWMutex) GetLock(timeout time.Duration) bool
func (lm *LRWMutex) RLock()
func (lm *LRWMutex) GetRLock(timeout time.Duration) bool
func (lm *LRWMutex) Unlock()
func (lm *LRWMutex) RUnlock()
```
