package main

import (
	"context"
	"fmt"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage"
)

// This is a storage.Appender that accepts scraped metrics data.
type noOpStore struct{}

func (s *noOpStore) Add(l labels.Labels, t int64, v float64) (uint64, error) {
	fmt.Printf("store.Add %v %d %v\n", l.String(), t, v)
	return 0, nil
}

func (s *noOpStore) AddFast(ref uint64, t int64, v float64) error {
	fmt.Printf("store.AddFast %d %d %v\n", ref, t, v)
	return nil
}

func (s *noOpStore) Commit() error {
	fmt.Println("store.Commit")
	return nil
}

func (s *noOpStore) Rollback() error {
	fmt.Println("store.Rollback")
	return nil
}

func (s *noOpStore) Appender(context.Context) storage.Appender {
	return s
}
