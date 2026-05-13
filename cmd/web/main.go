package main

import (
	"cmp"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type TopicFolder struct {
	Name  string
	Files []string
}

type MessageView struct {
	Topic   string
	ID      string
	Content string
}

func main() {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		slog.Error("Failed to parse templates", "error", err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "base", nil)
	})

	http.HandleFunc("/tree", func(w http.ResponseWriter, r *http.Request) {
		tree := getFileTree()
		tmpl.ExecuteTemplate(w, "tree", tree)
	})

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		topic := r.URL.Query().Get("topic")
		id := r.URL.Query().Get("id")
		if topic == "" || id == "" {
			http.Error(w, "Missing topic or id", http.StatusBadRequest)
			return
		}

		path := filepath.Join("messages", topic, id)
		content, err := os.ReadFile(path)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		tmpl.ExecuteTemplate(w, "message", MessageView{
			Topic:   topic,
			ID:      strings.TrimSuffix(id, ".json"),
			Content: string(content),
		})
	})

	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		entries, err := os.ReadDir("messages")
		if err != nil {
			slog.Error("Failed to read messages dir", "error", err)
			http.Error(w, "Failed to read messages dir", http.StatusInternalServerError)
			return
		}

		for _, entry := range entries {
			err := os.RemoveAll(filepath.Join("messages", entry.Name()))
			if err != nil {
				slog.Error("Error deleting topic directory", "topic", entry.Name(), "error", err)
			}
		}

		w.Header().Set("HX-Trigger", "messagesCleared")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "All messages cleared.")
	})

	port := 8080
	slog.Info("Server starting", "port", port, "url", fmt.Sprintf("http://localhost:%d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func getFileTree() []TopicFolder {
	var tree []TopicFolder
	messagesDir := "messages"

	entries, err := os.ReadDir(messagesDir)
	if err != nil {
		slog.Error("Error reading messages dir", "error", err)
		return tree
	}

	for _, entry := range entries {
		if entry.IsDir() {
			slog.Debug("Found topic directory", "name", entry.Name())
			topic := TopicFolder{Name: entry.Name()}
			files, err := os.ReadDir(filepath.Join(messagesDir, entry.Name()))
			if err == nil {
				for _, f := range files {
					if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
						topic.Files = append(topic.Files, f.Name())
					}
				}
				slog.Debug("Added files for topic", "topic", entry.Name(), "count", len(topic.Files))
				// Sort files descending (newest IDs/timestamps usually higher)
				slices.SortFunc(topic.Files, func(a, b string) int {
					return cmp.Compare(b, a)
				})
			} else {
				slog.Error("Error reading topic dir", "topic", entry.Name(), "error", err)
			}
			tree = append(tree, topic)
		}
	}

	slices.SortFunc(tree, func(a, b TopicFolder) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return tree
}
