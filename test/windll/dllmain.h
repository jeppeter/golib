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


#ifdef __cplusplus
};
#endif /* __cplusplus*/

#endif /* __DLLMAIN_H_C520D85A4FF792FDFE52168146B8544D__ */
