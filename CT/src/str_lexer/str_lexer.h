#pragma once

#define STATE_COUNT 7

typedef enum { false = 0, true = 1 } bool;

typedef enum {
    CHAR_TYPE_UNKNOWN,       // 0
    CHAR_TYPE_LETTER,        // 1 [A-Za-z0-9!@#$%^&*()_-+=<>,.?/
    CHAR_TYPE_QUOTES,        // 2 "
    CHAR_TYPE_SLASH,         // 3 "slash"
    CHAR_TYPE_ESCAPE_CHAR,   // 4 n|t
    CHAR_TYPE_COUNT          // 5
} char_type_t;



extern const int transitions[STATE_COUNT][CHAR_TYPE_COUNT];


extern const bool is_final[STATE_COUNT];

char_type_t get_char_type(char c);

bool check_string(char *s);