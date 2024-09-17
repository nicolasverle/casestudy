package extractor

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

type (
	urlExtractor struct {
		urlPaths []string
	}
)

func New(paths []string) *urlExtractor {
	return &urlExtractor{urlPaths: paths}
}

func (e *urlExtractor) ExtractAllLinks() ([]string, error) {
	eg := errgroup.Group{}

	parsedLinks := []string{}

	var parseNode func(path string, n *html.Node)

	parseNode = func(path string, n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if !strings.HasPrefix(attr.Val, "http") {
						parsedLinks = append(parsedLinks, fmt.Sprintf("%s%s", path, attr.Val))
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(path, c)
		}
	}

	for _, path := range e.urlPaths {
		eg.Go(func() error {
			resp, err := http.Get(path)
			if err != nil {
				return fmt.Errorf("error while querying url %s, %s", path, err.Error())
			}

			defer resp.Body.Close()

			nodes, err := html.Parse(resp.Body)
			if err != nil {
				return fmt.Errorf("error parsing html content of %s, %s", path, err.Error())
			}

			parseNode(path, nodes)

			return nil
		})

	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return sortAndDeduplicate(parsedLinks), nil
}

func sortAndDeduplicate(list []string) []string {
	slices.Sort(list)
	return slices.Compact(list)
}

func (e *urlExtractor) ToJSON(paths []string) ([]byte, error) {
	sorted := map[string][]string{}

	for _, p := range paths {
		parsedUrl, err := url.Parse(p)
		if err != nil {
			return nil, fmt.Errorf("error while parsing url %s for domain extraction, %s", parsedUrl, err.Error())
		}

		urlPrefix := fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)

		if _, ok := sorted[urlPrefix]; !ok {
			sorted[urlPrefix] = []string{}
		}

		sorted[urlPrefix] = append(sorted[urlPrefix], strings.ReplaceAll(p, urlPrefix, ""))
	}

	return json.MarshalIndent(&sorted, " ", "  ")
}

func init() {
	// not a good practice but for this case study, let's say it's ok...
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}
