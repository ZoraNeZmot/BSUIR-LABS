#include "lab5.h"
#include <string.h>
#include <math.h>
#include <stdio.h>

void print_word(char *word)
{
    int cur_len = 0, 
        len = strlen(word),
        i;

    int start_index = 0, end_index = len - 1;
    short is_end = 0;

    while (cur_len < len)
    {
        i = (is_end) ? end_index : start_index;
        if (word[i] == word[is_end ? i-1: i+1])
        {
            printf("%c", '_');
            cur_len += 2;
            if (is_end)
                end_index -= 2;
            else
                start_index += 2;

            is_end = !is_end;
        } else 
        {
            printf("%c", word[i]);   
            is_end ? end_index-- : start_index++;
            cur_len++;
        }

    }
}

void print_sine_wave(const char *text) 
{
    int len = strlen(text);
    for (int row = AMPLITUDE; row >= -AMPLITUDE; row--) 
    {
        for (int i = 0; i < len; i++) 
        {
            int y = (int)(AMPLITUDE * sin(FREQUENCY * i));
            if (y == row)
                printf("%c", text[i]);
            else
                printf(" ");
        }
        printf("\n");
    }
}