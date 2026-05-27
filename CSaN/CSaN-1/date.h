#pragma once

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Структура для хранения даты
typedef struct {
    int day;
    int month;
    int year;
} Date;

// Функция для проверки високосного года
int isLeapYear(int year);

// Функция для получения количества дней в месяце
int getDaysInMonth(int month, int year);

// Функция для проверки корректности даты
int isValidDate(Date date);

// Функция для парсинга строки с датой
int parseDate(const char* dateStr, Date* date);

// Функция для преобразования даты в количество дней от некоторой начальной точки
long long dateToDays(Date date);

// Основная функция для вычисления разницы между датами в днях
long long daysBetweenDates(const char* dateStr1, const char* dateStr2);