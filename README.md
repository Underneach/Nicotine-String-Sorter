# Nicotine-String-Sorter
Сортер ваших Url:Log:Pass строк на Go


![image](https://github.com/Underneach/Nicotine-String-Sorter/assets/137613889/df627566-40cf-43ec-b0d0-664c48749c7b)


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
