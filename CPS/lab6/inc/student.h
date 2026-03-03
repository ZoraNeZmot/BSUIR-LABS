#pragma once
#define true 1
#define false 0
#define len(arr) (sizeof(arr) / sizeof((arr)[0]))




typedef struct student
{
    char surname[30];
    short is_dormitory;
    float average_mark;
} student;


int compare_students(const void *a, const void *b);