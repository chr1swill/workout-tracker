#ifndef _SPAN_H
#define _SPAN_H

#ifndef _ASSERT_H
#include <assert.h>
#endif  //_ASSERT_H

#ifndef size_t
typedef unsigned long int size_t;
#endif

#define fspansizes \
	X(8192)          \
	X(4096)          \
	X(2048)          \
	X(1024)          \
	X(256)           \
	X(128)           \
	X(64)            \
	X(32)            \
	X(16)

#define X(size)                       \
	typedef struct {                    \
		char items[(size)-sizeof(size_t)];\
		size_t count;                     \
	} fspan ## size ## _t;

fspansizes
#undef X

#define mcopy(d, s, n)        \
do {                          \
	size_t __I__=0;             \
	for(; __I__ < (n) ; ++__I__)\
		(d)[__I__] = (s)[__I__];  \
} while(0);                   \

#define mmove(d, s, n)     \
do {                       \
	unsigned char __T__[(n)];\
	mcopy(__T__, (s), (n));  \
	mcopy((d), __T__, (n));  \
} while(0);                \

#define mappend(dst, src, n, dstoffset, dstcapacity)\
do {                                                \
  assert((dstcapacity) > (n) + (dstoffset));        \
	mcopy(&(dst)[(dstoffset)], (src), (n));           \
	(dstoffset) += (n);                               \
} while(0);                                         \

#define tmpspan_create(size)    \
	(span_t){ .capacity = (size)  \
		.items = (char [(size)]){0},\
		.count = 0 };               \

typedef struct {
	char *items;
	size_t count;
	size_t capacity;
} dspan_t;

#endif  //_SPAN_H
