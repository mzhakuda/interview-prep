package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type Account struct {
	mu      sync.Mutex
	id      int
	name    string
	balance int
}

func getSleepDuration() time.Duration {
	rnd := rand.IntN(2) + 1
	return time.Duration(rnd)
}

// transfer переводит amount со счёта from на счёт to.
// Используй TryLock для предотвращения deadlock (без lock ordering).
func transfer(from, to *Account, amount int) {
	// TODO: реализуй
	// - Lock(from), TryLock(to)
	// - Если TryLock failed: отпусти from, backoff, retry
	// - Если TryLock ok: выполни перевод, отпусти оба

	// Note that while correct uses of TryLock do exist, they are rare,
	// and use of TryLock is often a sign of a deeper problem

	for {
		from.mu.Lock()
		if to.mu.TryLock() {
			defer from.mu.Unlock()
			defer to.mu.Unlock()

			from.balance -= amount
			to.balance += amount

			return
		} else {
			from.mu.Unlock()
			time.Sleep(getSleepDuration() * time.Second)
		}
	}
}

func main() {
	alice := &Account{id: 1, name: "Alice", balance: 1000}
	bob := &Account{id: 2, name: "Bob", balance: 1000}

	fmt.Printf("Before: Alice=%d, Bob=%d, Total=%d\n",
		alice.balance, bob.balance, alice.balance+bob.balance)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		transfer(alice, bob, 300)
	}()

	go func() {
		defer wg.Done()
		transfer(bob, alice, 200)
	}()

	wg.Wait()

	fmt.Printf("After: Alice=%d, Bob=%d, Total=%d\n",
		alice.balance, bob.balance, alice.balance+bob.balance)
}
