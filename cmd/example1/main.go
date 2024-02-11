// This code is all GPT-4 generated
// The code demonstrates a simple HTTP server in Go that fetches user data
// from a mock database and caches the results using a generic cache.
// It showcases the use of Go's concurrency features, generics, and interfaces
// to build a concurrent, type-safe, and reusable cache for any data type.
package main

import (
	"encoding/json" // Used for encoding and decoding JSON data.
	"fmt"           // Provides I/O formatting functions.
	"log"           // Used for logging messages to show errors or information.
	"net/http"      // Package for building HTTP servers and clients.
	"strconv"       // String conversion, used here for converting string IDs to integers.
	"sync"          // Synchronization primitives like WaitGroup and Mutex.
	"time"          // Provides functionality for measuring and displaying time.
)

// User defines a model for storing user data. Struct tags like `json:"id"` specify how fields should
// be encoded to or decoded from JSON, enabling automatic conversion between JSON data and Go structs.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DataStore is an interface that abstracts the operations on data. It's a contract that any type
// can implement, provided it has the GetUser method with the correct signature. This allows for
// flexibility in the data source, whether it's a mock database or a real one.
type DataStore interface {
	GetUser(id int) User
}

// MockDB simulates a database using a map for storing users. It implements the DataStore interface,
// making it interchangeable with any other DataStore implementation. This is useful for testing or
// when swapping out the data layer without changing the rest of the code.
type MockDB struct {
	users map[int]User
}

// GetUser simulates fetching a user by ID from the database. It includes a deliberate delay to
// mimic database latency. This method satisfies the DataStore interface requirement.
func (db *MockDB) GetUser(id int) User {
	time.Sleep(100 * time.Millisecond) // Simulate latency
	return db.users[id]
}

// Cache uses Go generics to create a type-safe, reusable cache for any data type. Generics allow
// creating data structures and functions that work with any type, avoiding the need for type assertions
// and increasing code reusability and safety.
type Cache[T any] struct {
	data map[int]T    // Stores cached data. The type T is generic.
	mu   sync.RWMutex // Protects access to the map, ensuring safe concurrent access.
}

// NewCache is a constructor function for Cache. It initializes the map and returns a new Cache instance.
// The use of [T any] indicates that NewCache is a generic function that can operate on any type.
func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		data: make(map[int]T),
		mu:   sync.RWMutex{},
	}
}

// Get retrieves a value from the cache. If the value is not present, it uses the provided fetcher
// function to retrieve the value from the data source, stores it in the cache, and then returns it.
// This method demonstrates how to use generics, mutexes for thread-safe operations, and higher-order
// functions (functions that take other functions as parameters) in Go.
func (c *Cache[T]) Get(id int, fetcher func(int) T) T {
	c.mu.RLock()
	if val, ok := c.data[id]; ok {
		c.mu.RUnlock()
		return val
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if val, ok := c.data[id]; ok { // Double-check locking pattern to prevent race conditions.
		return val
	}

	val := fetcher(id)
	c.data[id] = val
	return val
}

func main() {
	db := &MockDB{
		users: map[int]User{
			1: {ID: 1, Name: "John Doe"},
			2: {ID: 2, Name: "Jane Doe"},
		},
	}
	userCache := NewCache[User]()

	// Sets up HTTP route handlers. Here, we define a route "/user" that takes one or more user IDs
	// as query parameters, fetches their data (potentially from cache), and returns it as JSON.
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		ids, ok := r.URL.Query()["id"]
		if !ok || len(ids[0]) < 1 {
			http.Error(w, "Missing user id", http.StatusBadRequest)
			return
		}

		var wg sync.WaitGroup
		users := make([]User, len(ids))
		for i, idStr := range ids {
			wg.Add(1)
			go func(i int, idStr string) { // Launch a goroutine for each user ID to fetch data concurrently.
				defer wg.Done()
				id, err := strconv.Atoi(idStr) // Convert the ID from string to int.
				if err != nil {
					// Error handling omitted for brevity.
					return
				}
				user := userCache.Get(id, db.GetUser) // Get user data, using cache if available.
				users[i] = user
			}(i, idStr)
		}
		wg.Wait() // Wait for all goroutines to complete.

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users) // Encode and return the user data as JSON.
	})

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil)) // Start the HTTP server.
}

// Explanation Within Code Comments
// Structs and JSON: User struct models user data and includes tags for JSON serialization, allowing easy conversion between Go structs and JSON.
// Interfaces: DataStore interface abstracts data operations, allowing for different implementations (e.g., a real database vs. a mock database) without changing the consuming code.
// Generics: Cache[T any] demonstrates using Go generics to create a type-safe, reusable cache for any data type. This avoids the need for type assertions and allows more compile-time checks.
// Concurrency: The fetching of user data is done concurrently using goroutines and a sync.WaitGroup to synchronize the completion of all fetch operations. This showcases how Go makes concurrent programming straightforward and efficient.
// Mutexes: Cache uses a read-write mutex (sync.RWMutex) to safely access and modify the cached data from multiple goroutines, demonstrating how to handle concurrent read/write operations.
