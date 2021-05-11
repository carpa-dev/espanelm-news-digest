package providers

import (
	"context"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

type RssGetter interface {
	Get(context.Context, string) (*gofeed.Feed, error)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Elpais struct {
	RssGetter
	HttpClient
}

type Items []struct {
	gofeed.Item

	// TODO
	language string
}

func (e *Elpais) FetchPagesList() ([]Page, error) {
	ctx := context.TODO()

	feed, err := e.RSS(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]Page, 0, len(feed.Items))

	// for each item
	// TODO go routine
	for _, item := range feed.Items {
		request, err := http.NewRequest("GET", item.Link, nil)
		if err != nil {
			return nil, err
		}

		res, err := e.HttpClient.Do(request)
		// Maybe we don't want to fail when there's an error
		if err != nil {
			return nil, err
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}

		article := Page{}
		doc.Find(`link[rel="alternate"]`).Each(func(i int, s *goquery.Selection) {
			link := Link{}

			href, ok := s.Attr("href")
			if !ok {
				return
			}

			lang, ok := s.Attr("hreflang")
			if !ok {
				return
			}

			link.Lang = lang
			link.Url = href

			article.Links = append(article.Links, link)
		})

		// article link to something else other than itself
		if len(article.Links) > 1 {
			responses = append(responses, article)
		}

	}

	return responses, nil
}

// RSS scrapes the RSS feed
// validates that the language is supported
// and then return the feed
func (e *Elpais) RSS(ctx context.Context) (*gofeed.Feed, error) {
	feed, err := e.RssGetter.Get(ctx, "https://feeds.elpais.com/mrss-s/pages/ep/site/brasil.elpais.com/portada")
	if err != nil {
		return nil, err
	}

	// Validate it's a language we support
	switch feed.Language {
	case "pt-br", "es":
	default:
		return nil, ErrInvalidLanguage
	}

	return feed, nil
}
