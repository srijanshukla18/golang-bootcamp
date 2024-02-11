package main

import (
	"fmt"
	"log"
)

// Book represents a book in the library
type Book struct {
	ID     int
	Title  string
	Author string
}

// Library manages a collection of books
type Library struct {
	books map[int]Book // Using a map for easy retrieval by ID
}

// NewLibrary initializes and returns a new instance of a Library
func NewLibrary() *Library {
	return &Library{
		books: make(map[int]Book),
	}
}

// AddBook is a method of Library. Methods are functions that execute in the context of a type.
// This method adds a book to the library.
func (l *Library) AddBook(book Book) {
	l.books[book.ID] = book
}

// GetBook is a method to retrieve a book by its ID. Demonstrates returning multiple values.
func (l *Library) GetBook(id int) (Book, bool) {
	book, exists := l.books[id]
	return book, exists
}

// ListBooks is a method that lists all books in the library. Demonstrates range over maps.
func (l *Library) ListBooks() []Book {
	var books []Book
	for _, book := range l.books {
		books = append(books, book)
	}
	return books
}

// recoverFromPanic demonstrates the use of defer, panic, and recover to handle unexpected errors.
func recoverFromPanic() {
	if r := recover(); r != nil {
		log.Println("Recovered from panic:", r)
	}
}

// Demonstrating functions vs. methods, panic and recovery, and complex data types.
func main() {
	defer recoverFromPanic() // Will recover from any panic in main

	library := NewLibrary()

	// Adding books to the library - Demonstrates structs and methods
	library.AddBook(Book{ID: 1, Title: "Go Programming", Author: "John Doe"})
	library.AddBook(Book{ID: 2, Title: "Advanced Go", Author: "Jane Smith"})

	// Retrieving a book - Demonstrates handling multiple return values
	book, exists := library.GetBook(1)
	if exists {
		fmt.Println("Retrieved book:", book.Title)
	} else {
		fmt.Println("Book not found")
	}

	// Listing all books - Demonstrates slices and iteration over them
	fmt.Println("Listing all books:")
	for _, book := range library.ListBooks() {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
	}

	// Demonstrating panic and recovery by intentionally causing a panic
	fmt.Println("Causing a panic after listing all books...")
	panic("simulated panic")
}
