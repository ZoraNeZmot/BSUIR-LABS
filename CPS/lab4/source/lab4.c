#include <string.h>
#include <stdlib.h>

int char_compare(const void *a, const void *b)
{
    return (*(const char *)a - *(const char *)b);
}


int find_matchs(const char *text, const char *word) {
    int count = 0;
    char temp[200];

    strncpy(temp, text, sizeof(temp) - 1);
    temp[sizeof(temp) - 1] = '\0';  

    char *token = strtok(temp, " ");
    while (token != NULL) {
        if (strcmp(token, word) == 0) {
            count++;
        }
        token = strtok(NULL, " ");
    }

    return count;
}


void change_symbols(char *A, char *B, int m, int n)
{
    for (int i = 0, j = m; i < n  && i < strlen(B) && j < strlen(A); i++, j++)
    {
        A[j] = B[i];
    }
}

void deshifr(const char *shifr, char *res)
{
    size_t j = 0;
    for (size_t i = 0; i < strlen(shifr); i++)
    {
        if (shifr[i] >= 'A' && shifr[i] <= 'I')
            res[j++] = shifr[i];
        if (shifr[i] >= 'a' && shifr[i] <= 'i')
            res[j++] = shifr[i] - 32;
    }
    res[j] = '\0';

    qsort(res, strlen(res), sizeof(char), char_compare);
}