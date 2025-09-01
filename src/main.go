package main

import (
		_ "github.com/glebarez/go-sqlite"
		"html/template"
    "database/sql"
		"net/http"
		"fmt"
)

const (
  PORT = ":8080"
	DBFILEPATH = "./file.db"
)

func check(err error) {
	if err != nil {
		panic(err);
	}
}

func check_rollback(err *error, tx *sql.Tx) {
	var e error;

	if *err != nil {
    e = tx.Rollback();
	  if e != nil { *err = e; }
	}
}

func main() {
	var err error;
	var db *sql.DB;
	var mux *http.ServeMux;
	var tmpl *template.Template;

	db, err = sql.Open("sqlite", DBFILEPATH);
	check(err);
	defer db.Close();

	err = initdb(db);
	check(err);

	mux = http.NewServeMux();

	tmpl = template.Must(template.ParseGlob("tmpl/*.html"));

	mux.HandleFunc("/add_workout_log/",
	func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed",
			http.StatusMethodNotAllowed);
			return
		}
    
		var wl Workout;

		if wl.FromFormValuesCheckValidity(r) {
			wl.TimeStamp();

			err = wl.InsertToDB(db);
      if  err != nil {
				fmt.Println("error wl.inserttodb - %v", err);
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther);
	});

	mux.HandleFunc("/add_exercise_name/",
	func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed);
			return
		}

		var e Exercise;

		if e.FromFormValuesCheckValidity(r) {
		  err = e.InsertToDB(db);
		  if err != nil {
		    fmt.Println("exercise.InsertToDB - %v", err);
		  }
		}

		http.Redirect(w, r, "/", http.StatusSeeOther);
	});

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed);
			return
		}

		var data struct{ 
		  Exercises []Exercise
		  Weights   []float32
		};

		data.Exercises, err = dbselectallexercises(db);
		if err != nil {
			fmt.Printf("error oh no - %s\n", err);
			return

		  http.Error(w, "Internal server error",
			http.StatusInternalServerError);
			return
		}

		data.Weights = []float32{ 0.0, 5.0, 8.3, 9.2, 12.5,
			11.5, 15, 16, 19.3, 18.5, 22.5, 27.5, 32,
			34, 38.5, 40.5, 45};

		err = tmpl.ExecuteTemplate(w, "index.html", data);
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error",
					http.StatusInternalServerError);
			return
		}
	});

	fmt.Printf("server running on port %s\n", PORT);
	err = http.ListenAndServe(PORT, mux);
	if (err != nil) {
		panic(err);
	}
}
