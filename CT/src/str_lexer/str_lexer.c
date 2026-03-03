#include "str_lexer.h"
#include "string.h"
#include <ctype.h>

// Defining constants
const int transitions[STATE_COUNT][CHAR_TYPE_COUNT] = {
    { 0, 0, 0, 0, 0 },
    { 0, 0, 2, 0, 0 },
    { 0, 3, 6, 4, 3 },
    { 0, 3, 6, 4, 3 },
    { 0, 0, 5, 5, 5 },
    { 0, 3, 6, 4, 3 }, // A | B | C
    { 0, 0, 0, 0, 0 }
};
const bool is_final[STATE_COUNT] = {0, 0, 0, 0, 0, 0, 1};

bool is_letter(char c) 
{
    if (c == '\\' || c == '"')
        return false;
        
    if ((c >= 'a' && c <= 'z') || 
        (c >= 'A' && c <= 'Z') || 
        (c == ' ') || 
        (c >= '0' && c <= '9') ||
        (ispunct((unsigned char)c))) 
        return true;

    return false;
}

// Declared in str_lexer.h
char_type_t get_char_type(char c) 
{
    if (c == 'n' || c == 't')
        return CHAR_TYPE_ESCAPE_CHAR;

    if (is_letter(c))
        return CHAR_TYPE_LETTER;

    
    switch (c)
    {
    case '"':
        return CHAR_TYPE_QUOTES;
    case '\\':
        return CHAR_TYPE_SLASH;
    default:
        return CHAR_TYPE_UNKNOWN;
    }
}

bool check_string(char *s)
{
    int current_state = 1;
    for (int i = 0; i < strlen(s); i++)
    {
        current_state = transitions[current_state][get_char_type(s[i])];
    }
    
    return is_final[current_state];
}
