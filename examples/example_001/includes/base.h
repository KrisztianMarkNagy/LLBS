
#pragma once
#ifndef PROJ_BASE_H
#define PROJ_BASE_H

#ifndef _STDIO_H
#include <stdio.h>
#endif//_STDIO_H

#ifndef _UNISTDIO_H
#include <unistdio.h>
#endif//_UNISTDIO_H

#ifndef _GLIBCXX_STDLIB_H
#include <stdlib.h>
#endif//_GLIBCXX_STDLIB_H || _STDLIB_H

#ifndef _UNISTD_H
#include <unistd.h>
#endif//_UNISTD_H

#ifndef _SYS_WAIT_H
#include <sys/wait.h>
#endif//_SYS_WAIT_H

#ifndef _ERRNO_H
#include <errno.h>
#endif//_ERRNO_H

#ifndef __STDBOOL_H
#include <stdbool.h>
#endif//__STDBOOL_H

#ifndef _STRING_H
#include <string.h>
#endif//_STRING_H

#ifndef _STRINGS_H
#include <strings.h>
#endif//_STRINGS_H

#ifndef _CTYPE_H
#include <ctype.h>
#endif//_CTYPE_H

#ifndef _TIME_H
#include <time.h>
#endif//_TIME_H

#ifndef lua_h
#include <luajit-2.1/lua.h>
#endif//lua_h

#ifndef _GLIBCXX_MATH_H
#include <math.h>
#endif//_GLIBCXX_MATH_H

#ifndef _PTHREAD_H
#include <pthread.h>
#endif//_PTHREAD_H


// * Template
// #ifndef macro
// #include <>
// #endif//macro

// *********************
// #include <GL/...>
// #include <GLES2/...>
// #include <GLES3/...>
// #include <vulkan/vulkan.h>
// #include <vulkan/vulkan_core.h>
// #include <vulkan/vulkan_android.h>
// #include <vulkan/vulkan_screen.h>
// #include <SDL2/SDL_vulkan.h>
// #include <SDL2/SDL.h>
// #include <cairo/cairo.h>
// *********************

#endif//PROJ_BASE_H