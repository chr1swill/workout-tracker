package main

import (
		_ "github.com/glebarez/go-sqlite"
		"html/template"
    "database/sql"
		"net/http"
		"reflect"
		"errors"
		"fmt"
)

const (
  PORT = ":8080"
	DBFILEPATH = "./file.db"
)

type ExerciseType int;

const (
  ETCardio ExerciseType = iota
  ETStrength
)

type Exercise struct {
	Id   uint64
	Name string
	Type ExerciseType
}

type Weight struct {
	Id          uint64
	WeightInLbs float32
}

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
	);`);
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

func dbSelectAllByPosition[T any](db *sql.DB, table string) ([]T, error) {
	var result []T

	// confirm T is a struct
	var sample T
	t := reflect.TypeOf(sample)
	if t == nil || t.Kind() != reflect.Struct {
		return nil, errors.New("generic type T must be a struct")
	}

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s;", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// For each row, create a *T, build a slice of pointers to its fields, rows.Scan into them.
	for rows.Next() {
		// create pointer to a new T
		ptr := reflect.New(t)    // *T, Value
		val := ptr.Elem()        // T, Value

		// prepare []interface{} of field addresses
		num := t.NumField()
		scanArgs := make([]interface{}, num)
		for i := 0; i < num; i++ {
			field := val.Field(i)
			// use Addr() only for settable fields (exported). Unexported fields will panic.
			if !field.CanAddr() {
				return nil, fmt.Errorf("field %s is unexported; cannot scan into it", t.Field(i).Name)
			}
			scanArgs[i] = field.Addr().Interface()
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		// append the dereferenced struct value
		result = append(result, ptr.Elem().Interface().(T))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func dbselectall[T []interface{}](db *sql.DB, table string) ([]T, error) {
	var item T;
	var err error;
	var tx *sql.Tx;
	var result []T;
	var rows *sql.Rows;

	if t := reflect.TypeOf(item); t.Kind() != reflect.Struct {
		return nil, fmt.Errorf(
		  "error - generic fuction type was not a struct\n");
	}

	tx, err = db.Begin();
	if err != nil { return nil, err };
	
	rows, err = tx.Query(fmt.Sprintf("select * from %s;", table));
	check_rollback(&err, tx);
	if err != nil { return nil, err };
  defer rows.Close();
  
	for rows.Next() {
		//for i := 0; i < v.NumField(); i++ {
		//	item[vType.Field(i).Name] = v.Field(i).Interface();
		//}

    //fields := reflect.VisibleFields(reflect.TypeOf(items));
		//for i, field := range fields {
		//	i
		//}

		err = rows.Scan(&item);
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
		err = rows.Scan(&item);
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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed);
			return
		}

		var weights []Weight;
		var exercises []Exercise;

		exercises, err = dbSelectAllByPosition[Exercise](db, "exercises");
		if err != nil {
			print(err)
		  http.Error(w, "Internal server error",
			http.StatusInternalServerError);
			return
		}

		weights, err = dbSelectAllByPosition[Weight](db, "weights");
		if err != nil {
			print(err)
		  http.Error(w, "Internal server error",
			http.StatusInternalServerError);
			return
		}

		err = tmpl.ExecuteTemplate(w, "index.html",
				struct{
					Exercises []Exercise
					Weights   []Weight }{ 
					Exercises: exercises,
					Weights: weights });
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
