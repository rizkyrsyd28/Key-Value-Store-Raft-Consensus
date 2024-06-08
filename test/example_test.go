package test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type BankAccount struct {
	RWMutex sync.RWMutex
	Balance int
}

func (b *BankAccount) Add(amount int) {
	b.RWMutex.Lock()
	b.Balance = b.Balance + amount
	b.RWMutex.Unlock()
}
func (b *BankAccount) Get() int {
	// b.RWMutex.RLock()
	val := b.Balance
	// b.RWMutex.RUnlock()
	return val
}

func TestExample(t *testing.T) {
	result := "A"
	if result != "A" {
		t.Fatalf("Seharusnya A")
	}
}
func TestExampleA(t *testing.T) {
	a := BankAccount{}

	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				a.Add(1)
				fmt.Println(a.Get())
			}
		}()
	}

	time.Sleep(5 * time.Second)
	fmt.Println("Final : ", a.Get())
}
