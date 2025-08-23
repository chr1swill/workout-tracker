#include <stdlib.h>
#include <unistd.h>
#include <sys/stat.h>

#define ROOT "/home/chr1swill/code/projects/old/c/workout-tracker"
#define SRC ROOT"/src"
#define BIN ROOT"/bin"

#define CC "gcc"
#define CFLAGS "-Wall -Wextra -ggdb"

int main()
{
	mkdir(BIN, -1);
	system(CC" "CFLAGS" -o "BIN"/main "SRC"/main.c");
	//system(BIN"/main");
	return 0;
}
