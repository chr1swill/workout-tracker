#include <fcntl.h>
#include <errno.h>
#include <stdio.h>
#include <assert.h>
#include <unistd.h>
#include <sys/stat.h>
#include "dynamicarray.h"
#define ARENA_IMPLEMENTATION
#include "arena.h"
#include "filedata.h"
#include "tui.h"

enum appstate_t {
  as_name_initial = 0x0,
  as_name_insert,
};

static size_t listposition = 0;
static struct termios org = {0};
static enum appstate_t as = as_name_initial;

#define puts_then_goto_label(label, msg)\
do {                                    \
  {                                     \
    puts((msg));                        \
		goto label;                         \
	}                                     \
} while(0);                             \

#define perror_then_goto_label(label, msg)\
do {                                      \
  {                                       \
    perror((msg));                        \
		goto label;                           \
	}                                       \
} while(0);                               \

#define ispunct(c)                  \ 
			((c) >= 0x21 && (c) <= 0x2f &&\
			 (c) >= 0x3a && (c) <= 0x40 &&\
			 (c) >= 0x7b && (c) <= 0x7e &&\
			 (c) >= 0x5b && (c) <= 0x60 ) \

#define isalphanum(c)                       \ 
							((c) >= 0x30 && (c) <= 0x39 &&\
							 (c) >= 0x41 && (c) <= 0x5a &&\
							 (c) >= 0x61 && (c) <= 0x7a ) \

#define isspace(c) ((c) == ' ')

#define isvalidexcercisenamechar(c)\
							(isalpanum((c))||    \
							 isspace((c))  ||    \
							 ispunct((c))   )    \
int main()
{
	dynamicarray_t da = {0};
	unsigned char c;
	struct stat st;
	filedata_t fd;
	arena_t a;
	ssize_t n;
	int dbfd;

	settmodraw();

	if (at_init(&a) == NULL)
	{
		puts_then_goto_label(clean_return, "errror - at_init");
	}

	dbfd = open(dbfilepath, O_RDWR|O_CREAT, S_IRUSR|S_IWUSR|S_IRGRP|S_IWGRP|S_IROTH);
	if (dbfd == -1)
	{
		perror_then_goto_label(clean_return, "error - open");
	}

	if (fstat(dbfd, &st) == -1)
	{
		puts_then_goto_label(clean_return, "errror - fstat");
	}
	
	// we knew the file that is support to be our db
	// is to small to be valid to we cant use its data
	// aka we have no data a can pretty much here as a
	// safty mesure empty the file
	if ((size_t)st.st_size < sizeof(filedata_t) && st.st_size != 0)
	{
		puts_then_goto_label(clean_return, "error - db file to small to be valid");
  }
  else if (st.st_size == 0)
	{
		as = as_name_insert;
		da_append(&da, "enter in your exercise name:\n",
				sizeof "enter in your exercise name:\n");
		da_append(&da, "> ", sizeof "> ");
	}
	else
	{
		assert((size_t)st.st_size < sizeof(filedata_t));
		// TODO: this is going to need to change with the
		// version with which could be made more flexible
		assert(st.st_size - (sizeof(size_t)*2) % sizeof(exercise_t));

		da_append(&da, "select name with 'j'/'k' or 'i' to insert:\n",
				sizeof "select name with 'j'/'k' or 'i' to insert:\n");

		n = read(dbfd, &fd, st.st_size);
		assert(n > 0);
		assert(n == st.st_size);
		assert(fd.version == dbfileversion);
		assert(fd.bytelength == st.st_size - (sizeof(size_t)*2));

		for (size_t i = 0; i < (fd.bytelength/sizeof(exercise_t)); ++i)
		{
			if (i != 0)
			{
				da_append(&da, "  ", sizeof("  "));
			}
			else
			{
				da_append(&da, "> ", sizeof("> "));
			}

			da_append(&da, fd.data[i].name, exercisenamemax);
		}
	}

	n = write(STDOUT_FILENO, da.buffer, da.count);
	if ((size_t)n != da.count)
	{
		puts_then_goto_label(clean_return, "error - write");
	}
	//da_reset(&da);

	while (1)
	{
		n = read(STDIN_FILENO, &c, sizeof(unsigned char));
		if (n != sizeof(unsigned char))
		{
		  puts_then_goto_label(clean_return, "errror - read");
		}

		if (as == as_name_initial)
		{
			if (c == '\n')
				goto clean_return;
			else if (c == 'i')
				as = as_name_insert;
		}
		else if (as == as_name_insert)
		{
			if (c == '\n' || c == 0x1b)
			{
				// lock in you chose for the exercise name 
			}
			else if (c == '\b')
			{
				// need to pop form cstr if needed
				//da_cstrlen();
				//--da.count;
			}
			else if (isvalidexcercisenamechar(c))
			{
				//as = as_name_insert;
				//
				// do something to listpostion
				//
				//da_append(&da, "enter in your exercise name:\n",
				//		sizeof "enter in your exercise name:\n");
				//da_append(&da, "> ", sizeof "> ");
				//da_append(&da, &c, sizeof char);
				//da.buffer[da.count-1];
				
				da_append(&da, &c, sizeof char);

				// typedef struct {
				// 	unsigned char name[exercisenamemax];
				// 	size_t        duration;
				// 	size_t        distance;
				// 	float         weight;
				// 	unsigned char nrep;
				// } exercise_t;

				size_t namelen;

				da_cstrlen();

				da_cstrappend(cstr, src, n, idxofnullbyte)
				da_cstrappend(fd.data[listpostion].name, &c, sizeof char);
				cstr, src, n, NT_IDX (da_cstrlen maybe), 
				"wht\0    "
				 
				"whtsrc  "

cstr[NT_IDX+n] = '\0';
				"whtsrc  "
			}
			else
			{
			}
		}

		n = write(STDOUT_FILENO, da.buffer, da.count);
		if ((size_t)n != da.count)
		{
			puts_then_goto_label(clean_return, "error - write");
		}
		//da_reset(&da);
	}

clean_return:
	at_free(a.buffer);
  settmodorg();
	close(dbfd);
	return 0;
}
