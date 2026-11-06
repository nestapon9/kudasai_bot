// estructuras XML (título, link, categorías, fecha, media:content)

package feed

import (
	"encoding/xml"
	"time"
)

type Post struct {
	Title     string
	URL       string
	PubDate   time.Time
	Category  string
	PhotoURL  string
	TotalPost int
	PostNum   int
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Item []struct {
			Title    string `xml:"title"`
			Link     string `xml:"link"`
			PubDate  string `xml:"pubDate"`
			Category string `xml:"category"`
			Content  struct {
				URL string `xml:"url,attr"`
			} `xml:"content"`
		} `xml:"item"`
	} `xml:"channel"`
}

type AtomFeed struct {
	Entry []struct {
		Title struct {
			Text string `xml:",chardata"`
		} `xml:"title"`
		Link struct {
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Published time.Time `xml:"published"`
		Category  struct {
			Term string `xml:"term,attr"`
		} `xml:"category"`
		Cmssection string `xml:"cmssection"`
	} `xml:"entry"`
}
