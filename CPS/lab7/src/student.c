#include "student.h"
#include <stdio.h>

float get_average_score(const int scores[5])
{
    int sum = 0;
    for (int i = 0; i < 5; i++)
    {
        sum += scores[i];
    }
    return sum / 5.0;
}

void filter_students()
{
    FILE *in = fopen("students.txt", "r"),
         *out = fopen("filtered.txt", "w");

    student temp;
    while (fscanf(in, "%d %s %d %d %d %d %d",
            &temp.id, temp.surname, &temp.scores[0], &temp.scores[1], &temp.scores[2], &temp.scores[3], &temp.scores[4]) == 7)
    {
        float av_score = get_average_score(temp.scores);
        if (av_score > 8)
        {
            fprintf(out, "%d %s %f\n", temp.id, temp.surname, av_score);
        }
    }
    fclose(in);
    fclose(out);
}

void filter_numbers()
{
    FILE *in = fopen("numbers.txt", "r"),
        *out = fopen("filt_numbers.txt", "w");

    int temp;
    while (fscanf(in, "%d", &temp) == 1)
    {
        if (temp & 1 == 1) fprintf(out, "%d ", temp);
    }
    
    fclose(in);
    fclose(out);
}