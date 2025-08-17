package main

import (
	"testing"
)

type Session struct {
	ID   int
	Name string
	Data [3024]byte // Simulating large data
}

func (s *Session) ToEntity() *Session {
	return &Session{ID: s.ID, Name: s.Name, Data: s.Data}
}

type UserPointer struct {
	Sessions []*Session
}

type UserValue struct {
	Sessions []Session
}

func generateSessions(n int) []*Session {
	sessions := make([]*Session, n)
	for i := 0; i < n; i++ {
		sessions[i] = &Session{ID: i, Name: "session"}
	}
	return sessions
}

func generateSessionValues(n int) []Session {
	sessions := make([]Session, n)
	for i := 0; i < n; i++ {
		sessions[i] = Session{ID: i, Name: "session"}
	}
	return sessions
}

func (u *UserPointer) ToEntity() []*Session {
	sessions := make([]Session, len(u.Sessions))
	for i, s := range u.Sessions {
		sessions[i] = *s.ToEntity()
	}
	result := make([]*Session, len(sessions))
	for i := range sessions {
		result[i] = &sessions[i]
	}
	return result
}

func (u *UserValue) ToEntity() []Session {
	sessions := make([]Session, len(u.Sessions))
	for i, s := range u.Sessions {
		sessions[i] = s
	}
	return sessions
}

// Benchmark helpers
func benchmarkPointer(b *testing.B, size int) {
	data := generateSessions(size)
	user := &UserPointer{Sessions: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.ToEntity()
	}
}

func benchmarkValue(b *testing.B, size int) {
	data := generateSessionValues(size)
	user := &UserValue{Sessions: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.ToEntity()
	}
}

// Benchmarks for 1,000 sessions
func BenchmarkPointerSliceToEntity_1k(b *testing.B) { benchmarkPointer(b, 1000) }
func BenchmarkValueSliceToEntity_1k(b *testing.B)   { benchmarkValue(b, 1000) }

// Benchmarks for 1,000,000 sessions
func BenchmarkPointerSliceToEntity_1M(b *testing.B) { benchmarkPointer(b, 1_000_000) }
func BenchmarkValueSliceToEntity_1M(b *testing.B)   { benchmarkValue(b, 1_000_000) }
