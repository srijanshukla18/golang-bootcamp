package main

import (
	"context" // Import the context package. Context is a mechanism to carry deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/book", handleBookRequest)
	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

// handleBookRequest is our HTTP handler function for book requests
func handleBookRequest(w http.ResponseWriter, r *http.Request) {
	// Here's where we first encounter 'context'. Think of a context as a small bag where you can
	// put some control instructions like "Hey, I can only wait for you to finish for 5 seconds!"
	// or "If I say stop, you need to stop what you're doing!". It's a way to manage the execution
	// of functions, particularly those that might take a long time, like network calls or database queries.
	//
	// We create a context with a timeout, which is like setting a stopwatch for 5 seconds.
	// If whatever we're doing (like fetching book details) doesn't finish in that time, the context
	// tells our code to stop waiting and move on. This prevents our server from getting stuck on requests that take too long.
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	// It's crucial to call 'cancel' when we're done with our context to clean up any resources.
	// 'defer' is a Go keyword that schedules this 'cancel' function to run as soon as our handleBookRequest function finishes.
	defer cancel()

	// Extract the book ID from the request URL query parameters.
	bookID := r.URL.Query().Get("id")
	if bookID == "" {
		http.Error(w, "missing book id", http.StatusBadRequest)
		return
	}

	// We then use this context when fetching book details. This way, all the operations
	// within 'fetchBookDetails' know they should try to finish before the context's deadline.
	bookDetails, err := fetchBookDetails(ctx, bookID)
	if err != nil {
		// If we didn't finish in time or there was another problem, we tell the client.
		// This could be because our code took too long (context's deadline exceeded) or
		// there was another error fetching the details.
		if err == context.DeadlineExceeded {
			http.Error(w, "request took too long", http.StatusGatewayTimeout)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// If all goes well, we respond with the book details.
	fmt.Fprintf(w, "Book Details: %+v\n", bookDetails)
}

// The remaining functions (fetchBookDetails, fetchFromDatabase, enrichWithExternalAPI)
// are designed to demonstrate how the context created above is passed through and used
// in potentially time-consuming operations, ensuring they adhere to the constraints (like timeouts)
// we've set up. This way, our application remains responsive, and resources are used efficiently.

// fetchBookDetails simulates fetching details for a book.
// It demonstrates how context can control the flow in a real-world scenario involving external calls.
func fetchBookDetails(ctx context.Context, bookID string) (string, error) {
	// First, attempt to fetch basic details from what we'll pretend is a database.
	details, err := fetchFromDatabase(ctx, bookID)
	if err != nil {
		return "", err // If there's an error (e.g., timeout), we stop and return the error.
	}

	// If successful, we then try to enrich the details by calling an external API,
	// again passing along our context to manage timeouts or cancellations.
	enrichedDetails, err := enrichWithExternalAPI(ctx, details)
	if err != nil {
		return "", err // Handle errors similarly, respecting the context's state.
	}

	return enrichedDetails, nil
}

// fetchFromDatabase simulates a database call.
// The key here is showing how a potentially long-running operation respects context.
func fetchFromDatabase(ctx context.Context, bookID string) (string, error) {
	// Simulate a delay to represent database latency.
	select {
	case <-time.After(2 * time.Second): // Wait for 2 seconds before "returning" data.
		return "Basic Book Details", nil
	case <-ctx.Done(): // If our context is cancelled or times out, stop waiting.
		// ctx.Err() tells us why the context was cancelled: a deadline exceeded or manual cancellation.
		return "", ctx.Err()
	}
}

// enrichWithExternalAPI simulates calling an external service to enrich the book details.
func enrichWithExternalAPI(ctx context.Context, details string) (string, error) {
	// Similarly, simulate a delay for the external API call.
	select {
	case <-time.After(2 * time.Second): // Pretend we're calling the API and it takes 2 seconds.
		return details + "; Enriched Details", nil
	case <-ctx.Done(): // Again, check if our context has been cancelled or timed out.
		return "", ctx.Err()
	}
}
