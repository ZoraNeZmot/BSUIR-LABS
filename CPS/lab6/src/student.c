#include "../inc/student.h"
#include <string.h>



int compare_students(const void *a, const void *b) 
{
    const student *s1 = (const student *)a;
    const student *s2 = (const student *)b;
    return strcmp(s1->surname, s2->surname);
}