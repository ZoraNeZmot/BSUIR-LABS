#include <stdio.h>
#include <stdlib.h>
#include <time.h>

int get_rand_range(int min, int max) 
{
    return min + (rand() % (max - min + 1));
}

int cmp_decr(const void *a, const void *b)
{
    return (*(int *)b - *(int *)a);
}

int cmp_incr(const void *a, const void *b)
{
    return (*(int *)a - *(int *)b);
}

int main() 
{
    srand((unsigned int)time(NULL));
    int n, m;
    printf("Enter sizes of matrix (n, m): ");
    scanf("%d %d", &n, &m);
    int matrix[n][m];

    for (int i = 0; i < n; i++) 
    {
        for (int j = 0; j < m; j++) 
        {
            matrix[i][j] = get_rand_range(0, 50);
        }
    }

    // Printing
    for (int i = 0; i < n; i++)
    {
        for (int j = 0; j < m; j++) 
        {
            printf("%4d", matrix[i][j]);
        }
        printf("\n");
    }

    // Sorting
    for (int i = 0; i < n; i++)
    {
        int even[m], odd[m];
        int e = 0, o = 0;

        for (int j = 0; j < m; j++)
        {
            if ((matrix[i][j] & 1) == 0) 
            {
                even[e++] = matrix[i][j];
            } else
            {
                odd[o++] = matrix[i][j];
            }
        }

            qsort(even, e, sizeof(int), cmp_decr);
            qsort(odd, o, sizeof(int), cmp_incr);
            e = o = 0;

            for (int j = 0; j < m; j++)
            {
                if ((matrix[i][j] & 1) == 0)
                {
                    matrix[i][j] = even[e++];
                } else 
                {
                    matrix[i][j] = odd[o++];
                }
            }
            
    }
    
    printf("\nSorted matrix:\n");
    for (int i = 0; i < n; i++) 
    {
        for (int j = 0; j < m; j++) 
        {
            printf("%4d", matrix[i][j]);
        }
        printf("\n");
    }

    puts("----------------------------------------\n");
    // ZAD 2
    printf("Enter sizes of matrix (n): ");
    scanf("%d", &n);
    int matrix2[n][n];
    for (int i = 0; i < n; i++) 
    {
        for (int j = 0; j < n; j++) 
        {
            matrix2[i][j] = get_rand_range(0, 200);
        }
    }
    for (int i = 0; i < n; i++)
    {
        for (int j = 0; j < n; j++) 
        {
            printf("%4d", matrix2[i][j]);
        }
        printf("\n");
    }
    int max_val = -1, max_i = -1, max_j = -1;
    for (int i = 0; i < n; i++)
    {
        if (max_val < matrix2[i][i])
        {
            max_i = i;
            max_j = i;

            max_val = matrix2[i][i];
        }
    }
    for (int i = 0; i < n; i++)
    {
        if (max_val < matrix2[i][n - 1 - i])
        {
            max_i = i;
            max_j = n - 1 - i;

            max_val = matrix2[i][n - 1 - i];
        }
    }
    int cross_i = n / 2;
    int cross_j = n / 2;

    int temp = matrix2[cross_i][cross_j];
    matrix2[cross_i][cross_j] = matrix2[max_i][max_j];
    matrix2[max_i][max_j] = temp;
    
    puts("\nAfter swaping:");
    for (int i = 0; i < n; i++)
    {
        for (int j = 0; j < n; j++) 
        {
            printf("%4d", matrix2[i][j]);
        }
        printf("\n");
    }

    return 0;
}