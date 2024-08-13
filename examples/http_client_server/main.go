package main

import (
	"github.com/NicolasBissig/gotel"
	"github.com/NicolasBissig/gotel/gotelhttp"
	"io"
	"net/http"
)

func main() {
	_, err := gotel.Setup()
	if err != nil {
		panic(err)
	}

	gotelhttp.InstrumentDefaultClient()

	gotelhttp.HandleFunc("GET /blogposts/{id}", blogpost)

	http.ListenAndServe(":8080", nil)
}

func blogpost(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL.
	id := r.PathValue("id")

	req, err := gotelhttp.NewRequest(r.Context(), http.MethodGet, "https://jsonplaceholder.typicode.com/posts/"+id, nil)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	resp, err := gotelhttp.Do(req)
	if err != nil {
		http.Error(w, "failed to fetch blog post", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Fetch the blog post with the given ID.
	blogPost, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to fetch blog post", http.StatusInternalServerError)
		return
	}

	// Write the blog post to the response.
	w.Write([]byte(blogPost))
}
