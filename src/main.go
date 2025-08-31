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

type AnyStruct interface {
	IsStruct()
}

type ExerciseType int;

const (
  ETCardio ExerciseType = iota
  ETStrength
)

const (
  ExerciseNameLenMin = 1
  ExerciseNameLenMax = 32
)

type Exercise struct {
	Id   uint64
	Name string
	Type ExerciseType
}

func (Exercise) IsStruct() {}

func (e *Exercise)fromFormValuesCheckValidity() bool {
	if len(e.Name) < ExerciseNameLenMin {
		return false;
  }

	if len(e.Name) > ExerciseNameLenMax {
		return false;
	}

	return true;
}

type Weight struct {
	Id          uint64
	WeightInLbs float32
}

func (Weight) IsStruct() {}

type WorkoutLog struct {
	Id         uint64
	ExerciseId uint64
	WeightId   uint64
	NReps      uint64
	Date       string
}

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

func initdb(db *sql.DB) error {
	var err error;
	var tx *sql.Tx;
	var stmt *sql.Stmt;
	var default_weight_values []float32;
	tx, err = db.Begin();
	check(err);

	_, err = tx.Exec(`
	pragma foreign_keys = ON;

	drop table weights;
	drop table workout_log;
	drop table exercises;

	create table if not exists exercises(
	  id integer primary key autoincrement,
	  type integer not null,
	  name string unique not null);

	create table if not exists weights(
	  id integer primary key autoincrement,
	  weight_in_lbs real unique);

  create table if not exists workout_log(
	  id integer primary key autoincrement,
	  exercise_id integer not null,
	  weight_id integer not null,
	  date string not null,
	  n_reps integer not null,
	  foreign key (exercise_id)
 	    references exercises (id)
	  foreign key (weight_id)
	   references weights (id)
	);

	insert into exercises(type, name) values (0, "bicep curl");
	`);
	check_rollback(&err, tx);
	if err != nil {
		return err;
	}

	stmt, err = tx.Prepare("insert into weights(weight_in_lbs) values (?);");
	defer stmt.Close();
	check_rollback(&err, tx);
	if err != nil {
		return err;
	}

	default_weight_values = []float32 { 0.0,
	5.0, 8.3, 9.2, 12.5, 11.5, 15, 16, 19.3,
	18.5, 22.5, 27.5, 32, 34, 38.5, 40.5, 45};

	for _, v := range default_weight_values {
		_, err  = stmt.Exec(v);
	  check_rollback(&err, tx);
	  if err != nil {
	  	return err;
	  }
	}

	err = tx.Commit()
	check_rollback(&err, tx);
	if err != nil { return err;
	} else {        return nil; }
}

func dbselectallexercises(db *sql.DB) ([]Exercise, error) {
	var item Exercise;
	var err error;
	var tx *sql.Tx;
	var result []Exercise;
	var rows *sql.Rows;

	tx, err = db.Begin();
	if err != nil { return nil, err };
	
	rows, err = tx.Query("select * from exercises;");
	check_rollback(&err, tx);
	if err != nil { return nil, err };
  defer rows.Close();
  
	for rows.Next() {
		var id    uint64
		var name  string
		var etype ExerciseType

		err = rows.Scan(&id, &etype, &name);
		check_rollback(&err, tx);
		if err != nil { return nil, err };

		item.Id = id;
		item.Type = etype;
		item.Name = name;

		result = append(result, Exercise{ Id: id, Type: etype, Name: name });
	}

	err = tx.Commit()
	check_rollback(&err, tx);
	if err != nil {
		return nil, err;
	} else {
		return result, nil;
	}
}

func dbselectallweights(db *sql.DB) ([]Weight, error) {
	var item Weight;
	var err error;
	var tx *sql.Tx;
	var result []Weight;
	var rows *sql.Rows;

	tx, err = db.Begin();
	if err != nil { return nil, err };
	
	rows, err = tx.Query("select * from weights;");
	check_rollback(&err, tx);
	if err != nil { return nil, err };
  defer rows.Close();
  
	for rows.Next() {
		err = rows.Scan(&item.Id, &item.WeightInLbs);
		check_rollback(&err, tx);
		if err != nil { return nil, err };

		result = append(result, item);
	}

	err = tx.Commit()
	check_rollback(&err, tx);
	if err != nil {
		return nil, err;
	} else {
		return result, nil;
	}
}

func dbinsertexercises(db *sql.DB, e *Exercise) error {
	var err error;
	var tx *sql.Tx;

	tx, err = db.Begin();
	check(err);
	
	_, err = tx.Exec(`
			insert into exercises(type, name) values (?, ?);`,
			e.Type, e.Name);
	check_rollback(&err, tx);
	if err != nil { return err };

	err = tx.Commit()
	check_rollback(&err, tx);
	if err != nil {
		return err;
	} else {
		return nil;
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

	mux.HandleFunc("/add_exercise_name/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed);
			return
		}

		var valid bool;
		var exercise Exercise;
		
		valid = true;
		if s := r.FormValue("type"); s != "" {
			if s == "strength" {
				exercise.Type = ETStrength;
			} else if s == "cardio" {
				exercise.Type = ETCardio;
			} else {
				valid = false;
			}
		} else {
			valid = false;
		}

		if s := r.FormValue("name"); s != "" {
			exercise.Name = s;
		} else {
			valid = false;
		}
		
		if valid {
			// add new exercise to db
			err = dbinsertexercises(db, &exercise);
			if err != nil {
				fmt.Println(err);
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
		  Weights   []Weight
		};

		data.Exercises, err = dbselectallexercises(db);
		if err != nil {
			fmt.Printf("error oh no - %s\n", err);
			return

		  http.Error(w, "Internal server error",
			http.StatusInternalServerError);
			return
		}

		data.Weights, err = dbselectallweights(db);
		if err != nil {
			fmt.Println(err)
		  http.Error(w, "Internal server error",
			http.StatusInternalServerError);
			return
		}

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
