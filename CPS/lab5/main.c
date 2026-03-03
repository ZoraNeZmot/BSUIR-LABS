#include <stdio.h>
#include "lab5.h"
#include <string.h>

int main()
{
    char word[MAX_LEN];
    printf("Enter word:\n");
    scanf("%s", word);
    while (getchar() != '\n'); 
    puts("Result:");
    print_word(word);
    puts("");

    char text[MAX_LEN];
    printf("Enter text:\n");
    fgets(text, sizeof(text), stdin);
    text[strcspn(text, "\n")] = '\0';
    print_sine_wave(text);
    puts("");
}   