# Nicotine-String-Sorter
Сортер ваших Url:Log:Pass строк на Go

## Что сортер умеет
+  Получение строк из файла или файлов в папке
+  Сохранение в виде Log:Pass или Url:Log:Pass
+  Сортировка по запросу в виде сайта (google.com) или ключевого слова в ссылке (google)
+  Многопоточная сортировка, запись в файлы и удаление дубликатов

## Стек
+  Многопоток - github.com/panjf2000/ants
+  Цветной вывод - github.com/fatih/color
+  Спеки проца - github.com/klauspost/cpuid
+  Получение кол-ва достпной оперы - github.com/pbnjay/memory
+  Определение кодировки файла - github.com/saintfish/chardet
+  Прогресс бар - github.com/schollz/progressbar
