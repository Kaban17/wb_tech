# Программы выведет
```bash
[3 2 3]
```
# Срезы
В `Go` срезы представлюят из себя указатель на данные, длину среза и ёмоксть(capacity).
При передаче среза в функцию, эти данные просто копируются. То есть мы видим
размер, данные, емкость. Когда мы внутри функции изменяем срез, то они меняются в копии среза внутри функции.
Именно поэтому измение элемента по индексу в данном коде повлияло на исходный массив, а добавление новых элементов не отразилось никак. Для этого в го даже функции `append` требует присваиваня нового слайса в исходный.
