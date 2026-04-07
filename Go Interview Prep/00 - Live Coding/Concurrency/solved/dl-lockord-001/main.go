package main

import (
	"fmt"
	"sync"
)

type Account struct {
	mu      sync.Mutex
	id      int
	name    string
	balance int
}

// transfer переводит amount со счёта from на счёт to.
// ⚠️ Этот код содержит deadlock. Найди и исправь.

// Горутины работают конкурентно, из за этого может случиться такое, что они оба локнут успеют первые
// и из за этого возникает дедлок, вот пример:

// Горутина 1: transfer(alice, bob)     Горутина 2: transfer(bob, alice)
// ─────────────────────────────────     ─────────────────────────────────
// 1. lock(alice.mu) ✅                  2. lock(bob.mu) ✅
// 3. lock(bob.mu) — BLOCKED ❌          4. lock(alice.mu) — BLOCKED ❌
//    (bob.mu держит горутина 2)            (alice.mu держит горутина 1)

// Lock ordering: всегда захватываем мьютексы в порядке возрастания id,
// чтобы исключить deadlock при встречных переводах.

func transfer(from, to *Account, amount int) {
	if from.id < to.id {
		from.mu.Lock()
		defer from.mu.Unlock()

		to.mu.Lock()
		defer to.mu.Unlock()
	} else {
		to.mu.Lock()
		defer to.mu.Unlock()

		from.mu.Lock()
		defer from.mu.Unlock()
	}

	if from.balance >= amount {
		from.balance -= amount
		to.balance += amount
		fmt.Printf("Transfer: %s -> %s: %d\n", from.name, to.name, amount)
	} else {
		fmt.Printf("Transfer: %s -> %s: insufficient funds\n", from.name, to.name)
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
		transfer(alice, bob, 300) // Alice -> Bob
	}()

	go func() {
		defer wg.Done()
		transfer(bob, alice, 200) // Bob -> Alice
	}()

	wg.Wait()

	fmt.Printf("After: Alice=%d, Bob=%d, Total=%d\n",
		alice.balance, bob.balance, alice.balance+bob.balance)
}
