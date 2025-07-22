
var justString string

func someFunc() {
	v := createHugeString(1 << 10)
	// Копируем первые 100 символов (байт), чтобы отбросить ссылку на весь v
	justString = string([]byte(v[:100]))
}

func main() {
	someFunc()
}
