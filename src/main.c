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

		da_append(&da, "select name with 'j'/'k' or 'i' to insert:\n",
				sizeof "select name with 'j'/'k' or 'i' to insert:\n");

		n = read(dbfd, &fd, st.st_size);
		assert(n > 0);
		assert(n == st.st_size);
		assert(fd.version == dbfileversion);
		assert(fd.count == st.st_size - (sizeof(size_t)*2));

		for (size_t i = 0; i < fd.count; ++i)
		{
			if (i != 0)
			{
				da_append(&da, "  ", sizeof("  "));
			}
			else
			{
				da_append(&da, "> ", sizeof("> "));
			}

			da_append(&da, fd.items[i].name.items, exercisenamemax);
		}
	}

	n = write(STDOUT_FILENO, da.buffer, da.count);
	if ((size_t)n != da.count)
	{
		puts_then_goto_label(clean_return, "error - write");
	}
	da_reset(&da);

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
				if(fd.items[listposition].name.count == 0)
				{
					puts_then_goto_label(clean_return,
					"errror - name not inserted");
				}

				// lock in you chose for the exercise name 
				fd_writetodbfd(&fd, dbfd);
				as = as_choose_weight;
#       define str "select weight:\n  0lbs"
				da_append(&da, str, sizeof(str));
#       undef str
			}
			else if (c == '\b')
			{
				--fd.items[listposition].name.count;
			}
			else if (isvalidexcercisenamechar(c))
			{
				da_append(&fd.items[listposition].name, &c, sizeof(char));
			}
			else
			{
			}
		}
		else if (as == as_choose_weight)
		{
		}
		else
		{
		}

    //format the da.buffer to write to display again

		n = write(STDOUT_FILENO, da.buffer, da.count);
		if ((size_t)n != da.count)
		{
			puts_then_goto_label(clean_return, "error - write");
		}
		da_reset(&da);
	}

clean_return:
	at_free(a.buffer);
  settmodorg();
	close(dbfd);
	return 0;
}
