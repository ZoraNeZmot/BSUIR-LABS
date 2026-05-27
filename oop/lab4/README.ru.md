# Лабораторная 4 — Плагины (иерархия)

Развивает лабораторную 3. Хост-приложение поставляется с семью
встроенными классами транспорта и при запуске сканирует папку
`plugins/`, подхватывая любые `*.json` файлы, описывающие новые
классы.

## Почему JSON-дескрипторы вместо нативных Go-плагинов

Стандартный пакет Go `plugin` работает только под Linux/macOS. На
Windows нативные плагины во время выполнения недоступны, поэтому
выбран кросс-платформенный подход на JSON. Он закрывает все
требования задания:

* новый модуль действительно расширяет иерархию (добавляет класс с
  собственным набором полей);
* все обобщённые функции (Marshal/Unmarshal, диалог редактирования,
  список) сразу его «понимают»;
* загрузка действительно **динамическая**: положили файл в
  `plugins/` и либо запустили программу, либо нажали
  «Reload plugins» — без пересборки и без правок кода.

## Формат плагина

```json
{
  "typeName": "Bicycle",
  "category": "Land",
  "summary": "[%s] %s %s (%s)",
  "summaryFields": ["TypeName", "Manufacturer", "Model", "Year"],
  "fields": [
    {"name": "ID",            "label": "Identifier",      "kind": "string"},
    {"name": "Manufacturer",  "label": "Manufacturer",    "kind": "string"},
    {"name": "Model",         "label": "Model",           "kind": "string"},
    {"name": "Year",          "label": "Year",            "kind": "int",   "default": "2024"},
    {"name": "GearCount",     "label": "Number of gears", "kind": "int",   "default": "21"},
    {"name": "FrameMaterial", "label": "Frame material",  "kind": "string"},
    {"name": "IsElectric",    "label": "Electric",        "kind": "bool"}
  ]
}
```

Поддерживаемые значения `kind`: `string`, `int`, `float`, `bool`.

В комплект включены два примера:

* `plugins/bicycle.json` — добавляет класс `Bicycle` в наземную ветку.
* `plugins/submarine.json` — добавляет класс `Submarine` в водную
  ветку.

## Сборка и запуск

```
cd lab4
go mod tidy
go run ./cmd/app
go run ./cmd/app -plugins=other/dir   # альтернативная папка плагинов
```

Готовая сборка: `bin/lab4/lab4.exe` (рядом лежит папка `plugins/`).

Кнопка «Add plugin file…» позволяет загрузить дескриптор, лежащий
вне выбранной по умолчанию папки.
