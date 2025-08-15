# Утилита grep

Реализовать утилиту фильтрации текстового потока (аналог команды grep).

Программа должна читать входной поток (STDIN или файл) и выводить строки, соответствующие заданному шаблону (подстроке или регулярному выражению).

Необходимо поддерживать следующие флаги:

    -A N — после каждой найденной строки дополнительно вывести N строк после неё (контекст).

    -B N — вывести N строк до каждой найденной строки.

    -C N — вывести N строк контекста вокруг найденной строки (включает и до, и после; эквивалентно -A N -B N).

    -c — выводить только то количество строк, что совпадающих с шаблоном (т.е. вместо самих строк — число).

    -i — игнорировать регистр.

    -v — инв��ртировать фильтр: выводить строки, не содержащие шаблон.

    -F — воспринимать шаблон как фиксированную строку, а не регулярное выражение (т.е. выполнять точное совпадение подстроки).

    -n — выводить номер строки перед каждой найденной строкой.

Программа должна поддерживать сочетания флагов (например, -C 2 -n -i – 2 строки контекста, вывод номеров, без учета регистра и т.д.).

Результат работы должен максимально соответствовать поведению команды UNIX grep.

Обязательно учесть пограничные случаи (начало/конец файла для контекста, повторяющиеся совпадения и пр.).

Код должен быть чистым, отформатированным (gofmt), работать без ситуаций гонки и успешно проходить golint.

---

## Examples

First, let's create a sample file named `input.txt`:
```
Hello World
hello world
HELLO WORLD
Go is fun
Golang is great
This is a test file
Another line
Case matters
Another world
Final line
```

**1. Basic Search**
Search for the string "world" (case-sensitive).
```bash
go run main.go world input.txt
```
Output:
```
hello world
Another world
```

**2. Case-Insensitive Search (-i)**
Search for "world", ignoring case.
```bash
go run main.go -i world input.txt
```
Output:
```
Hello World
hello world
HELLO WORLD
Another world
```

**3. Context After (-A)**
Find "fun" and show the next 2 lines.
```bash
go run main.go -A 2 fun input.txt
```
Output:
```
Go is fun
Golang is great
This is a test file
```

**4. Inverted Search with Line Numbers (-v -n)**
Find all lines that *do not* contain "world" and show their line numbers.
```bash
go run main.go -v -n world input.txt
```
Output:
```
1:Hello World
3:HELLO WORLD
4:Go is fun
5:Golang is great
6:This is a test file
7:Another line
8:Case matters
10:Final line
```

**5. Count Matches (-c)**
Count the number of lines containing "world" (case-insensitive).
```bash
go run main.go -c -i world input.txt
```
Output:
```
4
```