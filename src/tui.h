#ifndef _TUI_H
#define _TUI_H

#include <termios.h>

#define settmodraw() do {                    \
  tcgetattr(STDIN_FILENO, &org);             \
	struct termios __t__ = org;                \
  __t__.c_lflag &= ~(ECHO | ICANON);         \
	tcsetattr(STDIN_FILENO, TCSAFLUSH, &__t__);\
} while(0);                                  \

#define settmodorg() do {                  \
	tcsetattr(STDIN_FILENO, TCSAFLUSH, &org);\
} while(0);                                \

#endif  //_TUI_H
