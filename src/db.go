package main

import (
  "database/sql"
  "net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"fmt"
)

const (
  ExerciseNameLenMin = 1
  ExerciseNameLenMax = 32
  WNOTESMAXLEN       = 256
	SQLITE3_COL_MAX    = 2000
)

type User struct {
	Id             uint64
	Email          string
	HashedPassword string
	CreationDate   string
}

type Exercise struct {
	Id   uint64
	Name string
}

func (e *Exercise)FromFormValuesCheckValidity(r *http.Request) bool {
	var valid bool;

	valid = true;
	e.Name = r.FormValue("name");
	if len(e.Name) < ExerciseNameLenMin {
		valid = false;
	}

	if len(e.Name) > ExerciseNameLenMax {
		valid = false;
	}

	return valid;
}

func (e *Exercise)InsertToDB(db *sql.DB) error {
	var err error;
	var tx *sql.Tx;

	tx, err = db.Begin();
	check(err);

	_, err = tx.Exec(
			`insert into exercises(name) values (?);`, e.Name);
	check_rollback(&err, tx);
	if err != nil { return err };

	err = tx.Commit();
	check_rollback(&err, tx);
	if err != nil {
		return err;
	} else {
		return nil;
	}
}

type Workout struct {
	Id              uint64
	ExerciseId      uint64
	UserId          uint64
	Weight          float64
	DurationSeconds uint64
	NReps           uint64
	DistanceMetres  uint64
	Date            string
	Notes           string
}

func (w *Workout)FromFormValuesCheckValidity(r *http.Request) bool {
	var valid bool;

	valid = true;

	w.validasignexercise(&valid, r);
	w.validasignweight(&valid, r);
	w.validasignnreps(&valid, r);
	w.validasigndurationseconds(&valid, r);
	w.validasigndistancemetres(&valid, r);
	w.validasignnotes(&valid, r);

	return valid;
}

func (w *Workout)TimeStamp() {
	w.Date = time.Now().Format(time.DateTime);
}

func validasignintofloat64value(valid *bool, r *http.Request, wlvalue *float64, formvalue string) {
	var s string;
	var err error

	if s = r.FormValue(formvalue); s != "" {
		fmt.Printf("formvalue:%s=%s\n", formvalue, s);
		*wlvalue, err  = strconv.ParseFloat(s, 64);
		if err != nil {
			*valid = false;
		} else {
			fmt.Printf(
			"recieved value (%f) from formvalue (%s)\n",
			*wlvalue, formvalue);
		}
	} else {
		*valid = false;
	}
}

func validasignintouint64value(valid *bool, r *http.Request, wlvalue *uint64, formvalue string) {
	var s string;
	var err error

	if s = r.FormValue(formvalue); s != "" {
		fmt.Printf("formvalue:%s=%s\n", formvalue, s);
		*wlvalue, err  = strconv.ParseUint(s, 10, 64);
		if err != nil {
			*valid = false;
		} else {
			fmt.Printf(
			"recieved value (%d) from formvalue (%s)\n",
			*wlvalue, formvalue);
		}
	} else {
		*valid = false;
	}
}

func (w *Workout)validasignexercise(valid *bool, r *http.Request) {
	validasignintouint64value(valid, r, &w.ExerciseId, "exercise");
}

func (w *Workout)validasignweight(valid *bool, r *http.Request) {
	validasignintofloat64value(valid, r, &w.Weight, "weight");
}

func (w *Workout)validasignnreps(valid *bool, r *http.Request) {
	validasignintouint64value(valid, r, &w.NReps, "nreps");
}

func (w *Workout)validasigndurationseconds(valid *bool, r *http.Request) {
	validasignintouint64value(valid, r, &w.DurationSeconds, "durationseconds");
}

func (w *Workout)validasigndistancemetres(valid *bool, r *http.Request) {
	validasignintouint64value(valid, r, &w.DistanceMetres, "distancemetres");
}

func (w *Workout)validasignnotes(valid *bool, r *http.Request) {
	var s string;

	if s = r.FormValue("notes"); len(s) < WNOTESMAXLEN {
		w.Notes = s;
	} else {
		*valid = false;
	}
}

func(w *Workout)InsertToDB(db *sql.DB) error {
	var err error;
	var tx *sql.Tx;

	tx, err = db.Begin();
	check(err);

	check_rollback(&err, tx);
	if err != nil { return err };

	_, err = tx.Exec(
			`insert into workouts(
			exerciseid, userid, weight,
			durationseconds, nreps,
			distancemetres, date, notes)
			values (?, ?, ?, ?, ?, ?, ?);`,
			w.ExerciseId, w.UserId, w.Weight,
			w.DurationSeconds, w.NReps,
			w.DistanceMetres, w.Date, w.Notes);
	check_rollback(&err, tx);
	if err != nil { return err };

	err = tx.Commit();
	check_rollback(&err, tx);
	if err != nil {
		return err;
	} else {
		return nil;
	}
}

func initdb(db *sql.DB) error {
	var err error;
	var tx *sql.Tx;

	tx, err = db.Begin();
	check(err);

	_, err = tx.Exec(`
	pragma foreign_keys = ON;

	create table if not exists users(
	  id integer primary key autoincrement,
		email string unique not null,
		hashedpassword string not null,
    creationdate string not null);

	create table if not exists exercises(
	  id integer primary key autoincrement,
	  name string unique not null);

  create table if not exists workouts(
	  id integer primary key autoincrement,
	  exerciseid integer not null,
		userid integer not null,
	  weight real not null,
		durationseconds integer not null,
	  nreps integer not null,
	  distancemetres integer not null,
	  date string not null,
	  notes string,
	  foreign key (exerciseid)
 	    references exercises (id)
		foreign key (userid)
		  references users (id)
	);`);
	check_rollback(&err, tx);
	if err != nil {
		return err;
	}

	err = tx.Commit()
	check_rollback(&err, tx);
	if err != nil { return err;
	} else {        return nil; }
}

func gettable[T any](rows *sql.Rows) ([]T, error) {
	var n_struct_members int;
	var _struct reflect.Value;
	var field reflect.Value;
	var members []interface{}
	var table []T;
	var member T;

	for rows.Next() {
		_struct = reflect.ValueOf(&member).Elem();
		n_struct_members = _struct.NumField();
		members = make([]interface{}, n_struct_members);

		for i := 0; i < n_struct_members; i++ {
			field = _struct.Field(i);
			members[i] = field.Addr().Interface();
		}

		if err := rows.Scan(members...); err != nil {
			return nil, err;
		}

		table = append(table, member);
	}

	return table, nil;
}

func dbinsert[T any](db *sql.DB, item T) error {
	var err error;
	var tx *sql.Tx;
	var n_fields int;
	var has_field_id bool;
	var sb strings.Builder;
	var _type reflect.Type;
	var members []interface{};
	var _struct reflect.Value;

	_struct = reflect.ValueOf(&item).Elem();
	n_fields = _struct.NumField();
	_type = _struct.Type();

	if _struct.Kind() != reflect.Struct {
		return fmt.Errorf("dbinsert - generic type was not of Kind() -> reflect.Struct");
	}

	if n_fields > SQLITE3_COL_MAX {
		return fmt.Errorf("dbinsert - struct number of fields exceeds the number of columns allowed in a sqlite3 table");
	}

  _, err = sb.WriteString("insert into ");
	if err != nil { return err; }

	_, err = sb.WriteString(fmt.Sprintf("%ss(",
				strings.ToLower(_type.Name())));
	if err != nil { return err; }

	has_field_id = false;
	for i := range n_fields {
		var fieldname string;

		fieldname = strings.ToLower(_type.Field(i).Name);
		if fieldname != "id" {
			members = append(members, _struct.Field(i).Addr().Interface());

			if i == n_fields - 1 {
				_, err = sb.WriteString(fmt.Sprintf("%s", fieldname));
				if err != nil { return err; }
			} else {
				_, err = sb.WriteString(fmt.Sprintf("%s,", fieldname));
				if err != nil { return err; }
			}
		} else {
			has_field_id = true;
		}
  }	

	_, err = sb.WriteString(") values (?");
	if err != nil { return err; }

	if has_field_id {
		_, err = sb.WriteString(strings.Repeat(",?", n_fields - 2));
		if err != nil { return err; }
	} else {
		_, err = sb.WriteString(strings.Repeat(",?", n_fields - 1));
		if err != nil { return err; }
	}

	_, err = sb.WriteString(");");
	if err != nil { return err; }

	tx, err = db.Begin();
	if err != nil { return err; }
	
	_, err = tx.Exec(sb.String(), members...);
	check_rollback(&err, tx);
	if err != nil { return err };

	err = tx.Commit();
	check_rollback(&err, tx);
	if err != nil {
		return err;
	} else {
		return nil;
	}
}

func dbselectall[T any](db *sql.DB, table string) ([]T, error) {
	var member T;
	var err error;
	var tx *sql.Tx;
	var members []interface{};
	var tabledata []T
	var rows *sql.Rows;

	tx, err = db.Begin();
	check(err);

	rows, err = tx.Query(fmt.Sprintf("select * from %s;", table));
	check_rollback(&err, tx);
	if err != nil { return nil, err; }
	defer rows.Close();
	
	for rows.Next() {
    structure := reflect.ValueOf(&member).Elem();
    n_structure_fields := structure.NumField();
    
    members = make([]interface{}, n_structure_fields);
    
    for i := 0; i < n_structure_fields; i++ {
      field := structure.Field(i);
      members[i] = field.Addr().Interface();
    }

    err := rows.Scan(members...);
	  check_rollback(&err, tx);
    if err != nil { return nil, err; }

		tabledata = append(tabledata, member);
  }

	check_rollback(&err, tx);
	if err != nil { return nil, err };

	err = tx.Commit()
	check_rollback(&err, tx);
	if err != nil {
		return nil, err;
	} else {
		return tabledata, nil;
	}
}
