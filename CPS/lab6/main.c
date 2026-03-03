#include "inc/student.h"
#include <stdio.h>
#include <stdlib.h>

const student students[] = {
    { "Bbram",  true,  8.5f },
    { "Petrov",  false, 7.2f },
    { "Sidorov", true,  6.1f },
    { "Abram",  true,  8.5f },
    { "Smirnova",false, 8.1f },
    { "Kuznetsov",true,  7.9f }
};

int main()
{
    student filtered[len(students)];
    int count = 0;

    for (size_t i = 0; i < len(students); i++)
    {
        if (students[i].average_mark > 7.0f)
        {   
            filtered[count++] = students[i];
        }
    }

    qsort(filtered, count, sizeof(student), compare_students);

    for (size_t i = 0; i < count; i++)
    {
        printf("surname: %s, mark: %.1f\n", filtered[i].surname, filtered[i].average_mark);
    }
    
    return 0;
}