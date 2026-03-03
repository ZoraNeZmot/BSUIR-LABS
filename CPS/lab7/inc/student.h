#pragma once

#define MAX_LEN 50

typedef struct
{
    int id;
    char surname[MAX_LEN];
    int scores[5];
} student;

float get_average_score(const int[5]);

void filter_students();

void filter_numbers();