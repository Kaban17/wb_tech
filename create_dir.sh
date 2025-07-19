
#!/bin/bash

# Проверка аргументов
if [ $# -ne 2 ]; then
  echo "Usage: $0 <base_dir_name> <N>"
  echo "Example: $0 L42 5"
  exit 1
fi

BASEDIR="$1"
N="$2"

# Проверка, что N — число
if ! [[ "$N" =~ ^[0-9]+$ ]]; then
  echo "Error: N must be a positive integer."
  exit 2
fi

for ((i = 1; i <= N; i++)); do
  PROJECT_DIR="$BASEDIR/$i"
  mkdir -p "$PROJECT_DIR"
  cd "$PROJECT_DIR" || exit 3

  # Инициализация Go-модуля
  go mod init "wb_tech/${BASEDIR,,}_$i" > /dev/null 2>&1

  # Создание main.go
  cat > main.go <<EOF
package main

import "fmt"

func main() {
    fmt.Println("Hello, world from ${BASEDIR}/${i}!")
}
EOF

  cd - > /dev/null
done

echo "✅ Created $N Go projects in $BASEDIR/{1..$N}"
