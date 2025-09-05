package main

import (
		_ "github.com/glebarez/go-sqlite"
		"html/template"
    "database/sql"
		"net/http"
		"strconv"
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

var ADJUSTABLE_WEIGHT_SET_VALUES = []float32{ 0.0, 5.0, 8.3, 9.2, 12.5,
			11.5, 15, 16, 19.3, 18.5, 22.5, 27.5, 32,
			34, 38.5, 40.5, 45}
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

	mux.HandleFunc("/remove_workout/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed",
			http.StatusMethodNotAllowed);
			return;
		}

		var s string;
    var id uint64;
    var err error;
		var tx *sql.Tx;

		s = r.FormValue("id");
		if  s == "" {
			http.Redirect(w, r, "/all_workouts/", http.StatusSeeOther);
			return;
		}

		id, err = strconv.ParseUint(s, 10, 64);
		if err != nil {
			fmt.Printf("strconv.ParseUint - %v\n", err);
			http.Redirect(w, r, "/all_workouts/", http.StatusSeeOther);
			return;
		}

		tx, err  = db.Begin();
		if err != nil {
			fmt.Printf("db.Begin() - %v\n", err);
			http.Redirect(w, r, "/all_workouts/", http.StatusSeeOther);
			return;
		}

		_, err = tx.Exec("delete from workouts where id = ?;", id);
		check_rollback(&err, tx); 
		if err != nil {
			fmt.Printf("tx.Exec - %v\n", err);
			http.Redirect(w, r, "/all_workouts/", http.StatusSeeOther);
			return;
		}

		err = tx.Commit();
		check_rollback(&err, tx); 
		if err != nil { fmt.Printf("tx.Commit() - %v\n", err); }

		http.Redirect(w, r, "/all_workouts/", http.StatusSeeOther);
		return;
	});

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
				fmt.Printf("error wl.inserttodb - %v\n", err);
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

	mux.HandleFunc("/all_workouts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed);
			return
		}

		var data struct {
			Workouts []Workout;
			Exercises []Exercise;
		}

		data.Workouts, err = dbselectall[Workout](db, "workouts");
		if err != nil {
			fmt.Printf("dbselectallworkouts - %v\n", err);
		  http.Error(w, "Internal server error",
			http.StatusInternalServerError);
			return
		}

		data.Exercises, err = dbselectallexercises(db);
		if err != nil {
		  http.Error(w, "Internal server error",
			http.StatusInternalServerError);
			return
		}

		err = tmpl.ExecuteTemplate(w, "all_workouts.html", data);
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error",
					http.StatusInternalServerError);
			return
		}
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
