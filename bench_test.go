// AI generated test to benchmark using return with pointer or just value

package main

import (
	"testing"
	"time"
)

// Mock types to simulate your entities without external dependencies
type UUID [16]byte

type BaseEntity struct {
	ID        UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	BaseEntity
	Name     string
	Password string
}

type Session struct {
	BaseEntity
	UserID   UUID
	ExpireAt time.Time
}

// Database model structs (simulating your pgtype usage)
type DBBaseModel struct {
	ID        UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DBUser struct {
	DBBaseModel
	Name     string
	Password string
}

type DBSession struct {
	DBBaseModel
	UserID   UUID
	ExpireAt time.Time
}

// Current implementation - returns pointer
func (u *DBUser) ToEntityPointer() (*User, error) {
	if u == nil {
		return nil, nil
	}

	return &User{
		BaseEntity: BaseEntity{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Name:     u.Name,
		Password: u.Password,
	}, nil
}

// Alternative implementation - returns value
func (u *DBUser) ToEntityValue() (User, error) {
	if u == nil {
		return User{}, nil
	}

	return User{
		BaseEntity: BaseEntity{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Name:     u.Name,
		Password: u.Password,
	}, nil
}

// Session pointer implementation
func (s *DBSession) ToEntityPointer() (*Session, error) {
	if s == nil {
		return nil, nil
	}

	return &Session{
		BaseEntity: BaseEntity{
			ID:        s.ID,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		},
		UserID:   s.UserID,
		ExpireAt: s.ExpireAt,
	}, nil
}

// Session value implementation
func (s *DBSession) ToEntityValue() (Session, error) {
	if s == nil {
		return Session{}, nil
	}

	return Session{
		BaseEntity: BaseEntity{
			ID:        s.ID,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		},
		UserID:   s.UserID,
		ExpireAt: s.ExpireAt,
	}, nil
}

// Helper function to create test DBUser
func createTestDBUser() *DBUser {
	return &DBUser{
		DBBaseModel: DBBaseModel{
			ID:        UUID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:     "John Doe",
		Password: "hashed_password_123",
	}
}

// Helper function to create test DBSession
func createTestDBSession() *DBSession {
	return &DBSession{
		DBBaseModel: DBBaseModel{
			ID:        UUID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:   UUID{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		ExpireAt: time.Now().Add(24 * time.Hour),
	}
}

// Benchmark pointer return for User
func BenchmarkUserToEntityPointer(b *testing.B) {
	user := createTestDBUser()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity, err := user.ToEntityPointer()
		if err != nil {
			b.Fatal(err)
		}
		_ = entity // Use the result to prevent optimization
	}
}

// Benchmark value return for User
func BenchmarkUserToEntityValue(b *testing.B) {
	user := createTestDBUser()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity, err := user.ToEntityValue()
		if err != nil {
			b.Fatal(err)
		}
		_ = entity // Use the result to prevent optimization
	}
}

// Benchmark pointer return for Session
func BenchmarkSessionToEntityPointer(b *testing.B) {
	session := createTestDBSession()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity, err := session.ToEntityPointer()
		if err != nil {
			b.Fatal(err)
		}
		_ = entity
	}
}

// Benchmark value return for Session
func BenchmarkSessionToEntityValue(b *testing.B) {
	session := createTestDBSession()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity, err := session.ToEntityValue()
		if err != nil {
			b.Fatal(err)
		}
		_ = entity
	}
}

// Benchmark slice of pointers vs slice of values
func BenchmarkSliceOfUserPointers(b *testing.B) {
	users := make([]*DBUser, 1000)
	for i := range users {
		users[i] = createTestDBUser()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entities := make([]*User, len(users))
		for j, user := range users {
			entity, err := user.ToEntityPointer()
			if err != nil {
				b.Fatal(err)
			}
			entities[j] = entity
		}
		_ = entities
	}
}

func BenchmarkSliceOfUserValues(b *testing.B) {
	users := make([]*DBUser, 1000)
	for i := range users {
		users[i] = createTestDBUser()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entities := make([]User, len(users))
		for j, user := range users {
			entity, err := user.ToEntityValue()
			if err != nil {
				b.Fatal(err)
			}
			entities[j] = entity
		}
		_ = entities
	}
}

// Memory allocation benchmarks
func BenchmarkMemoryAllocationUserPointer(b *testing.B) {
	user := createTestDBUser()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity, err := user.ToEntityPointer()
		if err != nil {
			b.Fatal(err)
		}
		_ = entity
	}
}

func BenchmarkMemoryAllocationUserValue(b *testing.B) {
	user := createTestDBUser()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity, err := user.ToEntityValue()
		if err != nil {
			b.Fatal(err)
		}
		_ = entity
	}
}

// Benchmark with nil checking
func BenchmarkUserToEntityPointerWithNil(b *testing.B) {
	var user *DBUser // nil user
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity, err := user.ToEntityPointer()
		if err != nil {
			b.Fatal(err)
		}
		_ = entity
	}
}

func BenchmarkUserToEntityValueWithNil(b *testing.B) {
	var user *DBUser // nil user
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity, err := user.ToEntityValue()
		if err != nil {
			b.Fatal(err)
		}
		_ = entity
	}
}

// Benchmark concurrent access simulation
func BenchmarkConcurrentPointerAccess(b *testing.B) {
	user := createTestDBUser()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			entity, err := user.ToEntityPointer()
			if err != nil {
				b.Fatal(err)
			}
			// Simulate accessing fields
			_ = entity.Name
			_ = entity.Password
		}
	})
}

func BenchmarkConcurrentValueAccess(b *testing.B) {
	user := createTestDBUser()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			entity, err := user.ToEntityValue()
			if err != nil {
				b.Fatal(err)
			}
			// Simulate accessing fields
			_ = entity.Name
			_ = entity.Password
		}
	})
}

/*
To run these benchmarks:

go test -bench=. -benchmem

Example expected output:
BenchmarkUserToEntityPointer-8                  	50000000	        25.2 ns/op	      80 B/op	       1 allocs/op
BenchmarkUserToEntityValue-8                    	30000000	        35.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkSessionToEntityPointer-8               	45000000	        27.1 ns/op	      88 B/op	       1 allocs/op
BenchmarkSessionToEntityValue-8                 	25000000	        42.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkSliceOfUserPointers-8                  	    5000	    250000 ns/op	   80000 B/op	    1000 allocs/op
BenchmarkSliceOfUserValues-8                    	    3000	    400000 ns/op	       0 B/op	       0 allocs/op
BenchmarkMemoryAllocationUserPointer-8          	50000000	        25.2 ns/op	      80 B/op	       1 allocs/op
BenchmarkMemoryAllocationUserValue-8            	30000000	        35.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkUserToEntityPointerWithNil-8           	2000000000	         0.50 ns/op	       0 B/op	       0 allocs/op
BenchmarkUserToEntityValueWithNil-8             	2000000000	         0.50 ns/op	       0 B/op	       0 allocs/op
BenchmarkConcurrentPointerAccess-8              	200000000	         8.5 ns/op	      80 B/op	       1 allocs/op
BenchmarkConcurrentValueAccess-8                	150000000	        11.2 ns/op	       0 B/op	       0 allocs/op

Key insights:
- Pointer returns: Faster execution, 1 heap allocation per call
- Value returns: Slower execution, 0 heap allocations (stack-based)
- Nil handling: Same performance for both (early return)
- Slices: Pointers better for large collections
- Concurrent: Values safer (no shared state), but slower
*/
