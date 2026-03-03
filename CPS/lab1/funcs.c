#include "funcs.h"

void task13() 
{
    puts("TASK 13");
    double result = 1.0;
    
    for (double i = 0.1; i <= 10; i+=0.1) 
    {
        result *= 1.0 + sin(i);
    }
    printf("result: %.50e\n\n", result); 
}

void task17()
{
    puts("TASK 17");
    double result = 0;
    for (int i = 1; i <= 50; i++) 
    {
        result += (1.0/(i*i*i)); 
    }

    printf("result: %.50f\n\n", result);
}

void task25()
{
    puts("TASK 25");
    long result = 0;
    for (int i = 1; i <= 30; i++)
    {
        long a = i & 1 ? i : i/2;
        long b = i & 1 ? i*i : i*i*i;

        result += ((a - b)*(a - b));
    }
    
    printf("result: %ld\n\n", result);
}

void task26()
{
    puts("TASK 26");
    long n;
    puts("Enter the number:");
    scanf("%ld", &n);
    printf("divisors: ");
    for (int i = 1; i * i <= n; i++) 
    {
        if (n % i == 0) 
        {
            if (i*i != n)
            {
                printf("%ld %ld ", i, n / i);
            }
            else 
            {
                printf("%ld ", i);
            }

        }
    }
    printf("\n\n");
}

void task27()
{
    puts("TASK 27");
    long n;
    puts("Enter the number:");
    scanf("%ld", &n);
    printf("All Qs: ");
    for (int i = 1; i * i <= n; i++) 
    {
        if (n % (i*i) == 0 && n % (i*i*i) != 0) 
        {
            printf("%ld ", i);
        }
    }
    printf("\n\n");
}