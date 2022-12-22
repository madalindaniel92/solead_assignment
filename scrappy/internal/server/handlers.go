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

		// Get query parameter
		query := r.URL.Query()
		queryParam := query.Get("q")

		if queryParam == "" {
			description := "query parameter is missing"
			err := fmt.Errorf("%w: %s", ErrInvalidRequest, description)
			replyError(http.StatusUnprocessableEntity, w, r, err, description)
			return
		}

		results, err := client.SearchCompany(r.Context(), queryParam)
		if err != nil {
			replyError(http.StatusInternalServerError, w, r, err, "failed company search")
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
