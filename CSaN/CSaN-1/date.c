#include "date.h"

int isLeapYear(int year) {
    return (year % 4 == 0 && year % 100 != 0) || (year % 400 == 0);
}

int getDaysInMonth(int month, int year) {
    int daysInMonth[] = {31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31};
    
    if (month == 2 && isLeapYear(year)) {
        return 29;
    }
    return daysInMonth[month - 1];
}

int isValidDate(Date date) {
    if (date.year < 1 || date.month < 1 || date.month > 12 || date.day < 1) {
        return 0;
    }
    
    int maxDays = getDaysInMonth(date.month, date.year);
    return date.day <= maxDays;
}


// Функция для парсинга строки с датой
int parseDate(const char* dateStr, Date* date) {
    if (strlen(dateStr) != 10) {
        return 0;
    }

    if (dateStr[2] != '.' || dateStr[5] != '.') {
        return 0;
    }
    
    for (int i = 0; i < 10; i++) {
        if (i == 2 || i == 5) continue;
        if (dateStr[i] < '0' || dateStr[i] > '9') {
            return 0;
        }
    }
    
    char dayStr[3] = {dateStr[0], dateStr[1], '\0'};
    char monthStr[3] = {dateStr[3], dateStr[4], '\0'};
    char yearStr[5] = {dateStr[6], dateStr[7], dateStr[8], dateStr[9], '\0'};
    
    date->day = atoi(dayStr);
    date->month = atoi(monthStr);
    date->year = atoi(yearStr);
    
    return isValidDate(*date);
}

long long dateToDays(Date date) {
    long long days = 0;
    
    for (int y = 1; y < date.year; y++) {
        days += isLeapYear(y) ? 366 : 365;
    }

    for (int m = 1; m < date.month; m++) {
        days += getDaysInMonth(m, date.year);
    }
    
    days += date.day;
    
    return days;
}

long long daysBetweenDates(const char* dateStr1, const char* dateStr2) {
    Date date1, date2;
    
    if (!parseDate(dateStr1, &date1) || !parseDate(dateStr2, &date2)) {
        printf("Ошибка: неверный формат даты. Используйте формат ДД.ММ.ГГГГ\n");
        return -1;
    }
    
    long long days1 = dateToDays(date1);
    long long days2 = dateToDays(date2);
    
    long long diff = days2 - days1;
    if (diff < 0) {
        diff = -diff;
    }
    
    return diff;
}