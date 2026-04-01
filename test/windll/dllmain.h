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


#ifdef __cplusplus
};
#endif /* __cplusplus*/

#endif /* __DLLMAIN_H_C520D85A4FF792FDFE52168146B8544D__ */
