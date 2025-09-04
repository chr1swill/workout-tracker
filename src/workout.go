package main

import (
	"database/sql"
  "net/http"
	"strconv"
	"time"
	"fmt"
)

const (
  WNOTESMAXLEN        = 256
)

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
			exerciseid, weight,
			durationseconds, nreps,
			distancemetres, date, notes)
			values (?, ?, ?, ?, ?, ?, ?);`,
			w.ExerciseId, w.Weight,
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
