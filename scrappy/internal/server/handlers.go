package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func companiesHandler(state *State) http.HandlerFunc {
	client := state.client

	return func(w http.ResponseWriter, r *http.Request) {
		// Only GET HTTP method allowed
		if r.Method != http.MethodGet {
			replyError(http.StatusMethodNotAllowed, w, r, ErrMethodNotSupported, "")
			return
		}

		// Get query parameters
		queryParams := r.URL.Query()
		query := queryParams.Get("q")
		phone := queryParams.Get("phone")

		// Search results
		results, err := client.SearchCompany(r.Context(), query, phone)
		if err != nil {
			replyError(http.StatusInternalServerError, w, r, err, "failed company search")
			return
		}

		if results.Total == 0 {
			replyError(http.StatusNotFound, w, r, ErrNotFound, "companies not found")
			return
		}

		replyJSONContent(http.StatusOK, w, r, results)
	}
}

// Helpers

func replyJSONContent(status int, w http.ResponseWriter, r *http.Request, content any) {
	// Serialize content as JSON
	body, err := json.Marshal(content)
	if err != nil {
		err = fmt.Errorf("failed to serialize content %w", err)
		replyError(http.StatusInternalServerError, w, r, err, "")
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(body)
	logErr(err)
}

func replyError(status int, w http.ResponseWriter, r *http.Request, err error, extraText string) {
	log.Printf("%s %s: Error: %d %s\n", r.URL, r.Method, status, err)

	message := http.StatusText(status)
	if extraText != "" {
		message = fmt.Sprintf("%s - %s", message, extraText)
	}

	http.Error(w, message, status)
}

func logErr(err error) {
	if err != nil {
		log.Printf("Error: %s\n", err)
	}
}
