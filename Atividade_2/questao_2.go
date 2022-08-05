package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type MutexIntBuffer struct {
	mu       sync.Mutex
	capacity int
	buffer   []int
}

func NewMutexIntBuffer(capacity int) *MutexIntBuffer {
	new_mutex := new(MutexIntBuffer)
	new_mutex.capacity = capacity

	return new_mutex
}

func (s *MutexIntBuffer) push(n int) {
	s.mu.Lock()

	for len(s.buffer) == s.capacity {
		s.mu.Unlock()
		runtime.Gosched()
		s.mu.Lock()
	}

	s.buffer = append(s.buffer, n)
	s.mu.Unlock()
}

func (s *MutexIntBuffer) pop() int {
	s.mu.Lock()

	for len(s.buffer) == 0 {
		s.mu.Unlock()
		runtime.Gosched()
		s.mu.Lock()
	}
	ret := s.buffer[0]
	s.buffer = s.buffer[1:]
	s.mu.Unlock()
	return ret
}

func (s *MutexIntBuffer) checkAvailability() bool {
	s.mu.Lock()
	if len(s.buffer) == s.capacity {
		s.mu.Unlock()
		return false
	}
	s.mu.Unlock()
	return true
}

func (s *MutexIntBuffer) checkEmpty() bool {
	s.mu.Lock()
	if len(s.buffer) == 0 {
		s.mu.Unlock()
		return true
	}
	s.mu.Unlock()
	return false
}

func main() {
	sala_espera_capacity := 10
	fila_espera := NewMutexIntBuffer(sala_espera_capacity)

	go cortarCabelo(fila_espera)

	for i := 1; true; i++ {
		random_time := time.Duration(rand.Intn(5)) * time.Second
		time.Sleep(random_time)
		go desejoCortarCabelo(i, fila_espera)
	}

}

func desejoCortarCabelo(client_number int, fila_espera *MutexIntBuffer) {
	fmt.Println("Cliente " + strconv.Itoa(client_number) + " chega na barbearia.")
	space_available := fila_espera.checkAvailability()
	if space_available {
		fila_espera.push(client_number)
		fmt.Println("O cliente " + strconv.Itoa(client_number) + " entra na sala de espera.")
	} else {
		fmt.Println("Não há vaga na fila de espera para o cliente " + strconv.Itoa(client_number) + ", então ele vai embora da barbearia.")
	}
}

func cortarCabelo(fila_espera *MutexIntBuffer) {
	cutting_time := 2 * time.Second
	dormindo := true
	for true {
		sala_espera_vazia := fila_espera.checkEmpty()
		if sala_espera_vazia {
			dormindo = true
			fmt.Println("O barbeiro está dormindo pois não há clientes.")
		} else {
			cliente_atual := fila_espera.pop()

			if dormindo {
				fmt.Println("O cliente " + strconv.Itoa(cliente_atual) + " acorda o barbeiro.")
				dormindo = false
			}
			fmt.Println("O barbeiro está atendendo o cliente " + strconv.Itoa(cliente_atual) + ".")
			time.Sleep(cutting_time)
			fmt.Println("O barbeiro terminou de atender o cliente " + strconv.Itoa(cliente_atual) + ".")
		}
	}

}
