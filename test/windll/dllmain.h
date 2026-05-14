#ifndef __DLLMAIN_H_C520D85A4FF792FDFE52168146B8544D__
#define __DLLMAIN_H_C520D85A4FF792FDFE52168146B8544D__

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus*/


#if defined(WINLIB_DLL_IMPORT)
#define WINLIB_API  __declspec(dllimport)
#elif defined(WINLIB_DLL_EXPORT)
#define WINLIB_API __declspec(dllexport) 
#else
#define WINLIB_API
#endif

WINLIB_API int print_0();
WINLIB_API int print_1(int a0);
WINLIB_API int print_2(int a0, int a1);
WINLIB_API int print_3(int a0, int a1, int a2);
WINLIB_API int print_4(int a0, int a1, int a2, int a3);
WINLIB_API int print_5(int a0, int a1, int a2, int a3, int a4);
WINLIB_API int print_6(int a0, int a1, int a2, int a3, int a4, int a5);
WINLIB_API int print_7(int a0, int a1, int a2, int a3, int a4, int a5, int a6);
WINLIB_API int print_8(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7);
WINLIB_API int print_9(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8);
WINLIB_API int print_10(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9);
WINLIB_API int print_11(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10);
WINLIB_API int print_12(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11);
WINLIB_API int print_13(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11, int a12);
WINLIB_API int print_14(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11, int a12, int a13);
WINLIB_API int print_15(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11, int a12, int a13, int a14);
WINLIB_API int print_16(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11, int a12, int a13, int a14, int a15);
WINLIB_API int print_17(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11, int a12, int a13, int a14, int a15, int a16);
WINLIB_API int print_18(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11, int a12, int a13, int a14, int a15, int a16, int a17);


WINLIB_API char* charfunc_0();
WINLIB_API char* charfunc_1(int a0);
WINLIB_API char* charfunc_2(int a0, int a1);
WINLIB_API char* charfunc_3(int a0, int a1, int a2);
WINLIB_API char* charfunc_4(int a0, int a1, int a2, int a3);
WINLIB_API char* charfunc_5(int a0, int a1, int a2, int a3, int a4);
WINLIB_API char* charfunc_6(int a0, int a1, int a2, int a3, int a4, int a5);
WINLIB_API char* charfunc_7(int a0, int a1, int a2, int a3, int a4, int a5, int a6);
WINLIB_API char* charfunc_8(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7);
WINLIB_API char* charfunc_9(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8);
WINLIB_API char* charfunc_10(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9);
WINLIB_API char* charfunc_11(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10);
WINLIB_API char* charfunc_12(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11);
WINLIB_API char* charfunc_13(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11, int a12);
WINLIB_API char* charfunc_14(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11, int a12, int a13);
WINLIB_API char* charfunc_15(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11, int a12, int a13, int a14);
WINLIB_API char* charfunc_16(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11, int a12, int a13, int a14, int a15);
WINLIB_API char* charfunc_17(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11, int a12, int a13, int a14, int a15, int a16);
WINLIB_API char* charfunc_18(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9, int a10, int a11, int a12, int a13, int a14, int a15, int a16, int a17);


WINLIB_API char* strfunc_0();
WINLIB_API char* strfunc_1(char* a0);
WINLIB_API char* strfunc_2(char* a0, char* a1);
WINLIB_API char* strfunc_3(char* a0, char* a1, char* a2);
WINLIB_API char* strfunc_4(char* a0, char* a1, char* a2, char* a3);
WINLIB_API char* strfunc_5(char* a0, char* a1, char* a2, char* a3, char* a4);
WINLIB_API char* strfunc_6(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5);
WINLIB_API char* strfunc_7(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5, char* a6);
WINLIB_API char* strfunc_8(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5, char* a6, char* a7);
WINLIB_API char* strfunc_9(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5, char* a6, char* a7, char* a8);
WINLIB_API char* strfunc_10(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5, char* a6, char* a7, char* a8, char* a9);
WINLIB_API char* strfunc_11(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5, char* a6, char* a7, char* a8, char* a9, char* a10);
WINLIB_API char* strfunc_12(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5, char* a6, char* a7, char* a8, char* a9, char* a10, char* a11);
WINLIB_API char* strfunc_13(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5, char* a6, char* a7, char* a8, char* a9, char* a10, char* a11, char* a12);
WINLIB_API char* strfunc_14(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5, char* a6, char* a7, char* a8, char* a9, char* a10, char* a11, char* a12, char* a13);
WINLIB_API char* strfunc_15(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5, char* a6, char* a7, char* a8, char* a9, char* a10, char* a11, char* a12, char* a13, char* a14);
WINLIB_API char* strfunc_16(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5, char* a6, char* a7, char* a8, char* a9, char* a10, char* a11, char* a12, char* a13, char* a14, char* a15);
WINLIB_API char* strfunc_17(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5, char* a6, char* a7, char* a8, char* a9, char* a10, char* a11, char* a12, char* a13, char* a14, char* a15, char* a16);
WINLIB_API char* strfunc_18(char* a0, char* a1, char* a2, char* a3, char* a4, char* a5, char* a6, char* a7, char* a8, char* a9, char* a10, char* a11, char* a12, char* a13, char* a14, char* a15, char* a16, char* a17);


// to generate callback with 0
typedef int (callbackfunc_0_func_t)();
 
// to generate callback with 1
typedef int (callbackfunc_1_func_t)(char* a0);
 
// to generate callback with 2
typedef int (callbackfunc_2_func_t)(char* a0,char* a1);
 
// to generate callback with 3
typedef int (callbackfunc_3_func_t)(char* a0,char* a1,char* a2);
 
// to generate callback with 4
typedef int (callbackfunc_4_func_t)(char* a0,char* a1,char* a2,char* a3);
 
// to generate callback with 5
typedef int (callbackfunc_5_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4);
 
// to generate callback with 6
typedef int (callbackfunc_6_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5);
 
// to generate callback with 7
typedef int (callbackfunc_7_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6);
 
// to generate callback with 8
typedef int (callbackfunc_8_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7);
 
// to generate callback with 9
typedef int (callbackfunc_9_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8);
 
// to generate callback with 10
typedef int (callbackfunc_10_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9);
 
// to generate callback with 11
typedef int (callbackfunc_11_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10);
 
// to generate callback with 12
typedef int (callbackfunc_12_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11);
 
// to generate callback with 13
typedef int (callbackfunc_13_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12);
 
// to generate callback with 14
typedef int (callbackfunc_14_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13);
 
// to generate callback with 15
typedef int (callbackfunc_15_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13,char* a14);
 
// to generate callback with 16
typedef int (callbackfunc_16_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13,char* a14,char* a15);
 
// to generate callback with 17
typedef int (callbackfunc_17_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13,char* a14,char* a15,char* a16);
 
// to generate callback with 18
typedef int (callbackfunc_18_func_t)(char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13,char* a14,char* a15,char* a16,char* a17);
 
 
WINLIB_API int callbackfunc_0(callbackfunc_0_func_t pfunc);
WINLIB_API int callbackfunc_1(callbackfunc_1_func_t pfunc,char* a0);
WINLIB_API int callbackfunc_2(callbackfunc_2_func_t pfunc,char* a0,char* a1);
WINLIB_API int callbackfunc_3(callbackfunc_3_func_t pfunc,char* a0,char* a1,char* a2);
WINLIB_API int callbackfunc_4(callbackfunc_4_func_t pfunc,char* a0,char* a1,char* a2,char* a3);
WINLIB_API int callbackfunc_5(callbackfunc_5_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4);
WINLIB_API int callbackfunc_6(callbackfunc_6_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5);
WINLIB_API int callbackfunc_7(callbackfunc_7_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6);
WINLIB_API int callbackfunc_8(callbackfunc_8_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7);
WINLIB_API int callbackfunc_9(callbackfunc_9_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8);
WINLIB_API int callbackfunc_10(callbackfunc_10_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9);
WINLIB_API int callbackfunc_11(callbackfunc_11_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10);
WINLIB_API int callbackfunc_12(callbackfunc_12_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11);
WINLIB_API int callbackfunc_13(callbackfunc_13_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12);
WINLIB_API int callbackfunc_14(callbackfunc_14_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13);
WINLIB_API int callbackfunc_15(callbackfunc_15_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13,char* a14);
WINLIB_API int callbackfunc_16(callbackfunc_16_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13,char* a14,char* a15);
WINLIB_API int callbackfunc_17(callbackfunc_17_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13,char* a14,char* a15,char* a16);
WINLIB_API int callbackfunc_18(callbackfunc_18_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13,char* a14,char* a15,char* a16,char* a17);


typedef int (*stkcallbackfunc_func_t)(int size,char** pargs);
WINLIB_API int stkcallbackfunc_0(stkcallbackfunc_func_t pfunc);
WINLIB_API int stkcallbackfunc_1(stkcallbackfunc_func_t pfunc,char* a0);
WINLIB_API int stkcallbackfunc_2(stkcallbackfunc_func_t pfunc,char* a0,char* a1);
WINLIB_API int stkcallbackfunc_3(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2);
WINLIB_API int stkcallbackfunc_4(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3);
WINLIB_API int stkcallbackfunc_5(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4);
WINLIB_API int stkcallbackfunc_6(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5);
WINLIB_API int stkcallbackfunc_7(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6);
WINLIB_API int stkcallbackfunc_8(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7);
WINLIB_API int stkcallbackfunc_9(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8);
WINLIB_API int stkcallbackfunc_10(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9);
WINLIB_API int stkcallbackfunc_11(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10);
WINLIB_API int stkcallbackfunc_12(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11);
WINLIB_API int stkcallbackfunc_13(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12);
WINLIB_API int stkcallbackfunc_14(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13);
WINLIB_API int stkcallbackfunc_15(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13,char* a14);
WINLIB_API int stkcallbackfunc_16(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13,char* a14,char* a15);
WINLIB_API int stkcallbackfunc_17(stkcallbackfunc_func_t pfunc,char* a0,char* a1,char* a2,char* a3,char* a4,char* a5,char* a6,char* a7,char* a8,char* a9,char* a10,char* a11,char* a12,char* a13,char* a14,char* a15,char* a16);


#ifdef __cplusplus
};
#endif /* __cplusplus*/

#endif /* __DLLMAIN_H_C520D85A4FF792FDFE52168146B8544D__ */
