#include <fcntl.h>
#include <errno.h>
#include <stdio.h>
#include <assert.h>
#include <unistd.h>
#include <termios.h>
#include <sys/stat.h>
#include "dynamicarray.h"
#define ARENA_IMPLEMENTATION
#include "arena.h"

static struct termios org = {0};

#define exercisenamemax   16
#define exercisenreps     (unsigned char)(0-1)
#define dbfileversion     1
#define dbfiledir         "/home/chr1swill/code/projects/old/c/workout-tracker/"
#define dbfilepath        dbfiledir"/file.db"

#define settmodraw() do {                    \
  tcgetattr(STDIN_FILENO, &org);             \
	struct termios __t__ = org;                \
  __t__.c_lflag &= ~(ECHO | ICANON);         \
	tcsetattr(STDIN_FILENO, TCSAFLUSH, &__t__);\
} while(0);                                  \

#define settmodorg() do {                  \
	tcsetattr(STDIN_FILENO, TCSAFLUSH, &org);\
} while(0);                                \

typedef struct {
	unsigned char name[exercisenamemax];
	size_t        duration;
	size_t        distance;
	float         weight;
	unsigned char nrep;
} exercise_t;

// [file format]
//                 /* bytelenth/sizeof(exercise_t)=n_exercise_t */
// [size_t version|size_t bytelenth|exercise_t data[]]
typedef struct {
	size_t version;
	size_t bytelength;
	exercise_t data[];
} filedata_t;

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

int main()
{
	dynamicarray_t da = {0};
	unsigned char c;
	struct stat st;
	arena_t a;
	ssize_t n;
	int dbfd;

	settmodraw();

	if (at_init(&a) == NULL)
	{
		puts_then_goto_label(clean_return, "errror - at_init");
	}

	dbfd = open(dbfilepath, O_RDWR | O_CREAT, S_IRUSR | S_IWUSR | S_IRGRP | S_IWGRP | S_IROTH);
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
		puts_then_goto_label(clean_return, "error - db file to small to be valid");

  if (st.st_size == 0)
	{
		da_append(&da, "press 'i' to insert exercise name:\n",
				sizeof "press 'i' to insert exercise name:\n");
		da_append(&da, "> \n", sizeof "> \n");
	}
	else
	{
		assert((size_t)st.st_size < sizeof(filedata_t));
		// TODO: this is going to need to change with the
		// version with which could be made more flexible
		assert(st.st_size - (sizeof(size_t)*2) % sizeof(exercise_t));

		da_append(&da, "select name with 'j'/'k' or 'i' to insert:\n",
				sizeof "select name with 'j'/'k' or 'i' to insert:\n");

		filedata_t fd;

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
	if (n < 1)
	{
		puts_then_goto_label(clean_return, "error - write");
	}

	while (1)
	{
		n = read(STDIN_FILENO, &c, sizeof(unsigned char));
		if (n < 1)
		{
		  puts_then_goto_label(clean_return, "errror - read");
		}

		//TODO: DO SOMETHING
		// ...
		if (c == '\n') goto clean_return;
	}

clean_return:
	at_free(a.buffer);
  settmodorg();
	close(dbfd);
	return 0;
}
