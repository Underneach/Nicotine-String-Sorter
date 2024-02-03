# Nicotine-String-Sorter
Сортер ваших Url:Log:Pass строк на Go


![image](https://raw.githubusercontent.com/Underneach/Nicotine-String-Sorter/String-Sorter-regexp/image_1.png)
![image](https://raw.githubusercontent.com/Underneach/Nicotine-String-Sorter/String-Sorter-regexp/image_2.png)


## Что умеет сортер

    Получение строк из файла или файлов в папке
    Сохранение в виде Log:Pass или Url:Log:Pass
    Сортировка по запросу в виде сайта (google.com) или ключевого слова в ссылке (google)
    Многопоточная сортировка и одновременная запись в файлы с пропуском повторов строк - чтение базы любого размера

## Что умеет клинер

    Чистка базы любого размера - строки обрабатываются сразу при чтении, без загрузки списком в оперативную память
    Чистка нескольких баз по отдельности или всех баз в один файл
    Удаление невалид строк (A-z / 0-9 / Специмволы | 10-256 символов | UNKNOWN
    Удаление дублей реализовано через xxh3 хеш



## Стек
+  Многопоток - github.com/panjf2000/ants
+  Цветной вывод - github.com/fatih/color
+  Спеки проца - github.com/klauspost/cpuid
+  Получение кол-ва достпной оперы - github.com/pbnjay/memory
+  Определение кодировки файла - github.com/saintfish/chardet
+  Прогресс бар - github.com/schollz/progressbar
