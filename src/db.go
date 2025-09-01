package main

import (
  "database/sql"
  "net/http"
)

const (
  ExerciseNameLenMin = 1
  ExerciseNameLenMax = 32
)

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
	ExerciseId      uint64
	Weight          float64
	DurationSeconds uint64
	NReps           uint64
	DistanceMetres  uint64
	Date            string
	Notes           string
}

func initdb(db *sql.DB) error {
	var err error;
	var tx *sql.Tx;

	tx, err = db.Begin();
	check(err);

	_, err = tx.Exec(`
	pragma foreign_keys = ON;

	create table if not exists exercises(
	  id integer primary key autoincrement,
	  name string unique not null);

  create table if not exists workouts(
	  exerciseid integer not null,
	  weightid integer not null,
		durationseconds integer not null,
	  nreps integer not null,
	  distancemetres integer not null,
	  date string not null,
	  notes string,
	  foreign key (exerciseid)
 	    references exercises (id)
	  foreign key (weightid)
	   references weights (id)
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

func dbselectallexercises(db *sql.DB) ([]Exercise, error) {
	var err error;
	var tx *sql.Tx;
	var item Exercise;
	var result []Exercise;
	var rows *sql.Rows;

	tx, err = db.Begin();
	if err != nil { return nil, err };
	
	rows, err = tx.Query("select * from exercises;");
	check_rollback(&err, tx);
	if err != nil { return nil, err };
  defer rows.Close();
  
	for rows.Next() {
		err = rows.Scan(&item.Id, &item.Name);
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

func intodbtransation(db *sql.DB, handler func() error) error {
	var err error;
	var tx *sql.Tx;

	tx, err = db.Begin();
	check(err);

	err = handler();
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

func dbinsertexercises(db *sql.DB, e *Exercise) error {
	var err error;
	var tx *sql.Tx;

	err = intodbtransation(db, func() error {
			_, err = tx.Exec(`
					insert into exercises(name) values (?);`, e.Name);
			return err;
	});
	return err;
}

func dbinsertworkoutlog(db *sql.DB, w *Workout) error {
	var err error;
	var tx *sql.Tx;

	err = intodbtransation(db, func() error {
			_, err = tx.Exec(`
					insert into workouts(exerciseid,
					weightid, date, n_reps, notes, durationseconds, distancemetres)
					values (?, ?, ?, ?, ?, ?, ?);`,
					w.ExerciseId, w.Weight, w.Date,
					w.NReps, w.Notes, w.DurationSeconds, w.DistanceMetres);
			return err;
	});
	return err;
}
