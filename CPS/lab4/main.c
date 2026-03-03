#include <stdio.h>
#include <string.h>
#include "include/lab4.h"


int main() 
{
    //1
    char text[200], word[20];
    printf("Enter text:\n");
    fgets(text, sizeof(text), stdin);
    text[strcspn(text, "\n")] = '\0';
    printf("Enter search word:\n");
    fgets(word, sizeof(word), stdin);
    word[strcspn(word, "\n")] = '\0';

    int count = find_matchs(text, word);
    printf("Count is %d\n", count);
    puts("");

    // 2
    int m, n;
    printf("Enter m, n:\n");
    scanf("%d %d", &m, &n);
    while (getchar() != '\n'); 

    char A[200], B[200];
    printf("Enter A:\n");
    fgets(A, sizeof(A), stdin);
    A[strcspn(A, "\n")] = '\0';

    printf("Enter B:\n");
    fgets(B, sizeof(B), stdin);
    B[strcspn(B, "\n")] = '\0';

    change_symbols(A, B, m, n);
    printf("New string is:\t%s\n", A);
    puts("");

    //3
    char shifr[200];
    printf("Enter shifr:\n");
    fgets(shifr, sizeof(shifr), stdin);
    shifr[strcspn(shifr, "\n")] = '\0';

    char res[200];
    deshifr(shifr, res);
    printf("Result is:\t%s\n", res);

    return 0;
}