В Go используется утиная типизация. Поэтому, если даже структура реализует несколько интерфейсов, то она может быть скастована к одному из этих интерфейсов.
Также в языке разделены данные и поведение кода. Интерфейс определяет поведение, а структура опрделеает данные.
# как устроены интерфейсы
```go
type iface struct {
    tab  *itab
    data unsafe.Pointer
}
```
`tab` - указатель на таблицу интерфейсов, которая хранит метаданные о типе и список методов, которые реализует структура.
`data` - указатель на данные, которые хранит структура.
`itable` создается "на лету" (late binding) при первом присвоении и затем кешируется, что обеспечивает эффективность.
# Пустой interface{}
Пустому интерфейсу удоволетворяет любой тип. Поскольку у пустого интерфейса нет никаких методов, то и `itable` для не него создавать не нужно.
