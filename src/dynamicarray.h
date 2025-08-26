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

#define da_memset(dst, c, n)                  \
do {                                          \
	for (size_t __I__ = 0; __I__ < (n); ++__I__)\
		(dst)[__I__] = (c);                       \
} while(0);                                   \

#define da_clean(da)                      \
do {                                      \
	da_reset((da));                         \
	da_memset((da)->buffer, 0, (da)->count);\
} while(0);                               \

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

#define da_safecstrlen(cstr, cstrlen, memmaxlen)                    \
do {                                                                \
	static_assert(sizeof((cstrlen))==sizeof(size_t));                 \
	(cstrlen) = 0;                                                    \
	if ((memmaxlen) != 0)                                             \
	{                                                                 \
		for(;(cstr)[(cstrlen)]!='\0'&&(cstrlen)<(memmaxlen);++(cstrlen))\
	}                                                                 \
	else                                                              \
	{                                                                 \
		while ((cstr)[(cstrlen)]!='\0') ++(cstrlen);                    \
	}                                                                 \
} while(0);                                                         \

#define da_cstrlen(cstr, cstrlen)      \
do {                                   \
	da_safecstrlen((cstr), (cstrlen), 0);\
} while(0);                            \

#define da_cstrcap(cstr) sizeof((cstr))

#define da_internal_append(dst, src, n, dstoffset, dstcapacity)\
do {                                                           \
  assert((dstcapacity) > (n) + (dstoffset));                   \
	da_fastbytecopy(&(dst)[(dstoffset)], (src), (n));            \
	(dstoffset) += (n);                                          \
} while(0);                                                    \

#define da_fastcstrappend(cstr, src, n, idxofnullbyte)\
do {                                                  \
	assert((idxofnullbyte) + (n) <= da_cstrcap((cstr)));\
	da_internal_append((cstr), (src), (n),     \
	__SL__, da_cstrcap((cstr)));               \
	(cstr)[__SL__] = '\0';                     \
} while(0);                                  \

#define da_safecstrappend(cstr, src, n)      \
do {                                         \
	size_t __SL__;                             \
	da_cstrlen((cstr), __SL__);                \
	assert(__SL__ + (n) <= da_cstrcap((cstr)));\
	da_internal_append((cstr), (src), (n),     \
	__SL__, da_cstrcap((cstr)));               \
	(cstr)[__SL__] = '\0';                     \
} while(0);                                  \

#define da_cstrappend(cstr, src, n, idxofnullbyte)\
do {                                              \
	if ((idxofnullbyte) < 0) {                      \
		da_safecstrappend((cstr), (src), (n));        \
	}                                               \
	else                                            \
	{                                               \
		da_fastcstrappend((cstr), (src), (n),         \
		(idxofnullbyte));                             \
	}                                               \
} while(0);                                       \

#define da_append(da, src, n)\
do {                         \
	da_internal_append(        \
	(da)->buffer, (src),       \
	(n), (da)->count,          \
	da_cap((da)));             \
} while(0);                  \

#define da_prepend(da, bytes, n)                                 \
do {                                                             \
	assert(da_cap((da)) > (n) + (da)->count);                      \
	da_safebytecopy(&(da)->buffer[(n)], (da)->buffer, (da)->count);\
	da_fastbytecopy((da)->buffer, (bytes), (n));                   \
	(da)->count += (n);                                            \
} while(0);                                                      \

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
// // output:
// 654321

#endif //_DYNAMICARRAY_H
