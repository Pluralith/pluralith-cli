package dblock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// Lock struct
type Lock struct {
	Id   int64
	Lock bool
}

// Populate new lock instance
func (L *Lock) GenerateLock() {
	// Set seed for rand
	rand.Seed(time.Now().UnixNano())

	L.Id = int64(rand.Intn(1000000))
	fmt.Println("  Process Id is: ", L.Id)
	L.Lock = true
}

// Get lock string
func (L *Lock) GetLockString() (string, error) {
	functionName := "GetLockString"

	lockBytes, marshalErr := json.MarshalIndent(L, "", " ")
	if marshalErr != nil {
		return "", fmt.Errorf("%v: %w", functionName, marshalErr)
	}

	return string(lockBytes), nil
}

// Set lock to given value
func (L *Lock) SetLock(value bool) (string, error) {
	L.Lock = value

	return L.GetLockString()
}

// Create lock instance to be used by dblock methods
var LockInstance = &Lock{}

// Generate Lock
// Before write -> Check lock
// 		Repeadely read lock file
// 		If locked:
// 			If processID current process -> perform write
// 			Else -> retry until unlocked (timeout with force-unlock)
// 		Else
//			Lock + write processID
// 			Perform write
