#ifndef _DYNAMICARRAY_H
#define _DYNAMICARRAY_H
#include <assert.h>

typedef struct {
	unsigned char buffer[2<<6];
	size_t count;
} dynamicarray_t;

#define da_cap(da) sizeof((da)->buffer)
#define da_len(da) (da)->count
#define da_reset(da) (da)->count = 0;

#define da_fastbytecopy(d, s, n) \
do {                             \
	size_t __I__=0;                \
	for(; __I__ < (n) ; ++__I__)   \
		(d)[__I__] = (s)[__I__];     \
} while(0);                      \

#define da_safebytecopy(d, s, n)   \
do {                               \
	unsigned char __T__[(n)];        \
	da_fastbytecopy(__T__, (s), (n));\
	da_fastbytecopy((d), __T__, (n));\
} while(0);                        \

#define da_append(da, bytes, n)                             \
do {                                                        \
	assert(da_cap((da)) > (n) + (da)->count);                 \
	da_fastbytecopy(&(da)->buffer[(da)->count], (bytes), (n));\
	(da)->count += (n);                                       \
} while(0);                                                 \

#define da_prepend(da, bytes, n)                                 \
do {                                                             \
	assert(da_cap((da)) > (n) + (da)->count);                      \
	da_safebytecopy(&(da)->buffer[(n)], (da)->buffer, (da)->count);\
	da_fastbytecopy((da)->buffer, (bytes), (n));                   \
	(da)->count += (n);                                            \
} while(0);                                                      \

#define da_indexof(da, str, strlen, result)\
do {                                       \
	assert((strlen) < da_cap((da)));         \
	ssize_t __T__ = 0;                       \
	size_t __I__ = 0;                        \
	for (;__I__ < (da)->count; ++__I__)      \


	
// EXAMPLE
//
// int main(void)
// {
// 	dynamicarray_t da = {0};
// 
// 	da_prepend(&da, "1", 1);
// 	da_prepend(&da, "2", 1);
// 	da_prepend(&da, "3", 1);
// 	da_prepend(&da, "4", 1);
// 	da_prepend(&da, "5", 1);
// 	da_prepend(&da, "6", 1);
// 
// 	write(STDOUT_FILENO, da.buffer, da_len(&da));
// 
// 	return 0;
// }

#endif //_DYNAMICARRAY_H
