#ifndef _FILEDATA_H
#define _FILEDATA_H

#ifndef _DYNAMICARRAY_H
#include "dynamicarray.h"
#endif  //_DYNAMICARRAY_H

#ifndef _ARENA_H
#define ARENA_IMPLEMENTATION
#include "arena.h"
#endif  //_ARENA_H

#ifndef _SPAN_H
#include "span.h"
#endif  // _SPAN_H

#define exercisenamemax   16
#define exercisenreps     (unsigned char)(0-1)
#define dbfileversion     1
#define dbfiledir         "/home/chr1swill/code/projects/old/c/workout-tracker/"
#define dbfilepath        dbfiledir"/file.db"

#define deffspansizeexercise_t(size)      \
typedef struct {                          \
	fspan ## size ## _t name;               \
	size_t                         duration;\
	size_t                         distance;\
	float                          weight;  \
	unsigned char                  nrep;    \
} exercise_t;

deffspansizeexercise_t(16);

typedef struct {
	size_t version;
	size_t count;
	exercise_t *items;
} filedata_t;

#define fd_isalphanum(c)                    \
							((c) >= 0x30 && (c) <= 0x39 &&\
							 (c) >= 0x41 && (c) <= 0x5a &&\
							 (c) >= 0x61 && (c) <= 0x7a ) \

#define fd_isspace(c) ((c) == ' ')

#define fd_ispunct(c)               \
			((c) >= 0x21 && (c) <= 0x2f &&\
			 (c) >= 0x3a && (c) <= 0x40 &&\
			 (c) >= 0x7b && (c) <= 0x7e &&\
			 (c) >= 0x5b && (c) <= 0x60 ) \

#define isvalidexcercisenamechar(c)\
							(fd_isalphanum((c))||\
							 fd_isspace((c))   ||\
							 fd_ispunct((c))   ) \

#endif  //_FILEDATA_H
