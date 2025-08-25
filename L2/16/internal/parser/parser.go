package parser

import (
	"bytes"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"wb_tech/l2_16/internal/storage"
)

// Result содержит результат парсинга: измененный HTML и найденные новые ссылки.
type Result struct {
	ModifiedHTML []byte
	NewLinks     []string
}

// ProcessHTML парсит HTML, заменяет ссылки на локальные и возвращает новые ссылки для скачивания.
func ProcessHTML(body io.Reader, baseURL, baseHost string) (*Result, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	newLinks := findAndProcessLinks(doc, baseURL, baseHost)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, err
	}

	return &Result{
		ModifiedHTML: buf.Bytes(),
		NewLinks:     newLinks,
	}, nil
}

// findAndProcessLinks рекурсивно обходит DOM, модифицирует ссылки и собирает новые.
func findAndProcessLinks(n *html.Node, baseURL, baseHost string) []string {
	var newLinks []string

	if n.Type == html.ElementNode {
		attrsToUpdate := []string{}
		if n.Data == "a" || n.Data == "link" {
			attrsToUpdate = append(attrsToUpdate, "href")
		} else if n.Data == "script" || n.Data == "img" || n.Data == "source" {
			attrsToUpdate = append(attrsToUpdate, "src")
		}

		for i, attr := range n.Attr {
			isAttrToUpdate := false
			for _, a := range attrsToUpdate {
				if attr.Key == a {
					isAttrToUpdate = true
					break
				}
			}

			if isAttrToUpdate {
				absURL, err := resolveURL(baseURL, attr.Val)
				if err != nil {
					continue
				}

				parsedAbsURL, err := url.Parse(absURL)
				if err != nil || parsedAbsURL.Host != baseHost {
					continue
				}

				localPath, err := storage.URLToPath(absURL)
				if err != nil {
					continue
				}
				
				n.Attr[i].Val = strings.TrimPrefix(localPath, baseHost+"/")
				newLinks = append(newLinks, absURL)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		newLinks = append(newLinks, findAndProcessLinks(c, baseURL, baseHost)...)
	}

	return newLinks
}

// resolveURL преобразует относительную ссылку в абсолютную.
func resolveURL(baseURL, relURL string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	rel, err := url.Parse(relURL)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(rel).String(), nil
}
