package main

import (
		"html/template"
		"net/http"
		"fmt"
)

const (
  PORT = ":8080"
)

func main() {
	var err error;
	var mux *http.ServeMux;
	var tmpl *template.Template;

	mux = http.NewServeMux();

	tmpl = template.Must(template.New("index.html").
			ParseGlob("tmpl/*.html"));

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed);
		}

		err = tmpl.ExecuteTemplate(w, "index.html", nil);
		if err != nil {
			http.Error(w, "Internal server error",
			http.StatusInternalServerError);
		}
	});

	fmt.Printf("server running on port %s\n", PORT);
	err = http.ListenAndServe(PORT, mux);
	if (err != nil) {
		panic(err);
	}
}
