#include <stdio.h>
#include <string.h>
#include "./str_lexer/str_lexer.h"

int main(int argc, char *argv[])
{
    printf("Enter string for check or __STOP__ to stop:\n");
    char s[256];
    while (true)
    {
        fgets(s, sizeof(s), stdin);
        s[strcspn(s, "\n")] = '\0';

        if (!strcmp(s, "__STOP__"))
        {
            return 0;
        }
        if (check_string(s)) 
        {
            printf("String: %s is VALID\n------------\n", s);
        } else 
        {
            printf("String: %s is inVALID\n------------\n", s);
        }
        printf("Substrings:\n");
        size_t len = strlen(s);
            for (size_t i = 0; i < len; i++)
            {
                size_t j = i + 1;
                bool is_valid = false;
                char sub[257];
                while (!is_valid && j <= len)
                {
                    strncpy(sub, s + i, j - i);
                    sub[j - i] = '\0';  
                    is_valid = check_string(sub);
                    j++;
                }
                if (is_valid && strcmp(s, sub)) 
                {
                    printf("VALID substring: %s\n", sub);
                    i = j - 2;
                }
            }
                printf("-----------\n");
            
    }
    return 0;
}