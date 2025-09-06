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

var err error;
// just a runtime const lol...
var ADJUSTABLE_WEIGHT_SET_VALUES = []float32{
	0.0, 5.0, 8.3, 9.2, 12.5, 11.5, 15, 16, 19.3,
	18.5, 22.5, 27.5, 32, 34, 38.5, 40.5, 45};
		
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

type App struct {
	Mux  *http.ServeMux;
	DB   *sql.DB;
	Tmpl *template.Template;
}

func (a *App) remove_workout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed",
				http.StatusMethodNotAllowed);
		return;
	}

	var s string;
	var id uint64;
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

	tx, err  = a.DB.Begin();
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
}

func (a *App) add_workout_log(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed",
				http.StatusMethodNotAllowed);
		return
	}

	var wl Workout;

	if wl.FromFormValuesCheckValidity(r) {
		wl.TimeStamp();

		err = dbinsert[Workout](a.DB, wl);
		if  err != nil {
			fmt.Printf("error wl.inserttodb - %v\n", err);
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther);
} 

func (a *App) add_exercise_name(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed",
				http.StatusMethodNotAllowed);
		return
	}

	var e Exercise;

	if e.FromFormValuesCheckValidity(r) {
		err = dbinsert[Exercise](a.DB, e);
		if err != nil {
			fmt.Printf("dbinsert[Exercise] - %v\n", err);
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther);
}

func (a *App) all_workouts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed);
		return
	}

	var data struct {
		Workouts []Workout;
		Exercises []Exercise;
	}

	data.Workouts, err = dbselectall[Workout](a.DB, "workouts");
	if err != nil {
		fmt.Printf("dbselectallworkouts - %v\n", err);
		http.Error(w, "Internal server error",
				http.StatusInternalServerError);
		return
	}

	data.Exercises, err = dbselectall[Exercise](a.DB, "exercises");
	if err != nil {
		http.Error(w, "Internal server error",
				http.StatusInternalServerError);
		return
	}

	err = a.Tmpl.ExecuteTemplate(w, "all_workouts.html", data);
	if err != nil {
		fmt.Println(err)
			http.Error(w, "Internal server error",
					http.StatusInternalServerError);
		return
	}
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed);
		return;
	}
	var email, password string;

	email = r.FormValue("email");
	password = r.FormValue("password");

	// for now
	http.Redirect(w, r, "/", http.StatusSeeOther);
}

func (a *App) homepage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed);
		return
	}

	var data struct{ 
		Exercises []Exercise
			Weights   []float32
	};

	data.Exercises, err = dbselectall[Exercise](a.DB, "exercises");
	if err != nil {
		http.Error(w, "Internal server error",
				http.StatusInternalServerError);
		return
	}

	data.Weights = ADJUSTABLE_WEIGHT_SET_VALUES;

	err = a.Tmpl.ExecuteTemplate(w, "home.html", data);
	if err != nil {
		fmt.Println(err);
		http.Error(w, "Internal server error",
				http.StatusInternalServerError);
		return
	}
}

func (a *App) indexpage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println(err);
		http.Error(w, "Method not allowed",
				http.StatusMethodNotAllowed);
			return
	}

	err = a.Tmpl.ExecuteTemplate(w, "index.html", nil);
	if err != nil {
		fmt.Println(err);
		http.Error(w, "Internal server error",
				http.StatusInternalServerError);
		return
	}
}

func main() {
	var a App;

	a.DB, err = sql.Open("sqlite", DBFILEPATH);
	check(err);
	defer a.DB.Close();

	err = initdb(a.DB);
	check(err);

	a.Mux = http.NewServeMux();

	a.Tmpl = template.Must(template.ParseGlob("tmpl/*.html"));

	a.Mux.HandleFunc("/remove_workout/", a.remove_workout);
	a.Mux.HandleFunc("/add_workout_log/", a.add_workout_log);
	a.Mux.HandleFunc("/add_exercise_name/", a.add_exercise_name);
	a.Mux.HandleFunc("/all_workouts/", a.all_workouts);
	a.Mux.HandleFunc("/login/", a.login);
	a.Mux.HandleFunc("/home/", a.homepage);
	a.Mux.HandleFunc("/", a.indexpage);

	fmt.Printf("server running on port %s\n", PORT);
	err = http.ListenAndServe(PORT, a.Mux);
	if (err != nil) {
		panic(err);
	}
}
