package seeders

import (
	"sync"
)

type Seeder struct {
	seeders []func() error
	mu      sync.Mutex
}

func NewSeeder() *Seeder {
	return &Seeder{}
}

func (sm *Seeder) AddSeeder(seeder func() error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.seeders = append(sm.seeders, seeder)
}

func (sm *Seeder) RunAllSeeders() error {
	for _, seeder := range sm.seeders {
		if err := seeder(); err != nil {
			return err
		}
	}
	return nil
}
