package storage

import (
	"net/url"
	"os"
	"path"
	"strings"
)

// URLToPath преобразует URL в путь для сохранения на локальном диске.
func URLToPath(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	p := u.Host
	// Добавляем путь, если он есть
	if u.Path != "" && u.Path != "/" {
		p = path.Join(p, u.Path)
	}

	// Если URL указывает на "директорию" (пустой путь, / на конце, или нет расширения),
	// то добавляем index.html.
	if u.Path == "" || strings.HasSuffix(u.Path, "/") || path.Ext(p) == "" {
		p = path.Join(p, "index.html")
	}

	return p, nil
}

// Save сохраняет данные по указанному локальному пути, создавая необходимые директории.
func Save(localPath string, data []byte) error {
	if err := os.MkdirAll(path.Dir(localPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(localPath, data, 0644)
}
