package ch

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// Scrap scraps a page and returns the url to the comic
func Scrap(id int, sem chan int, result chan string, wg *sync.WaitGroup) {

	defer wg.Done()
	sem <- 1

	tmp, err := worker("http://explosm.net/comics/" + strconv.Itoa(id))
	<-sem
	if err != nil {
		log.Printf("error: %v", err)
	}

	fmt.Println(tmp)
}

// Random fetches link to a random Comic
func Random() (string, error) {
	return worker("http://explosm.net/comics/random")
}

// Download  a file
func Download(url string) (string, *http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return "", nil, fmt.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	fileName := strings.Split(path.Base(url), "?")[0]

	return fileName, resp, nil
}

func worker(u string) (string, error) {

	// fmt.Printf("processing, http://explosm.net/comics/%d\n", id)

	req, _ := http.NewRequest("GET", u, nil)

	client := &http.Client{}

	req.Header.Set("Host", "files.explosm.net")
	req.Header.Set("Referer", u)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error in downloading page: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	h, err := html.Parse(resp.Body)
	if err != nil {
		if err != nil {
			return "", fmt.Errorf("error in parsing: %v", err)
		}
	}

	comicLink := ""

	var f func(node *html.Node)

	f = func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "meta" {
			tmp := findComic(n)
			if tmp != "" {
				comicLink = tmp
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(h)

	if comicLink == "" {
		return "", fmt.Errorf("error: no image link found")
	}

	return comicLink, nil
}

func findComic(n *html.Node) string {

	for _, v := range n.Attr {
		if v.Val == "og:image" {
			for _, vi := range n.Attr {
				if vi.Key == "content" {
					return vi.Val
				}
			}

		}
	}
	return ""
}
