#ifndef _FILEDATA_H
#define _FILEDATA_H

#ifndef _DYNAMICARRAY_H
#include "dynamicarray.h"
#endif  //_DYNAMICARRAY_H

#ifndef _ARENA_H
#define ARENA_IMPLEMENTATION
#include "arena.h"
#endif  //_ARENA_H

#define dbfileversion     1
#define dbfiledir         "/home/chr1swill/code/projects/old/c/workout-tracker/"
#define dbfilepath        dbfiledir"/file.db"

// [file format]
//                 /* bytelenth/sizeof(exercise_t)=n_exercise_t */
// [size_t version|size_t bytelenth|exercise_t data[]]
typedef struct {
	size_t version;
	size_t bytelength;
	exercise_t data[];
} filedata_t;

#define exercisenamemax   16
#define exercisenreps     (unsigned char)(0-1)

typedef struct {
	unsigned char name[exercisenamemax];
	size_t        duration;
	size_t        distance;
	float         weight;
	unsigned char nrep;
} exercise_t;

#endif  //_FILEDATA_H
