package crawler

import (
	"net/http"
	"strings"
	"time"

	"v2ex-tui/internal/model"

	"github.com/PuerkitoBio/goquery"
)

const (
	BaseURL = "https://www.v2ex.com"
)

// Crawler manages fetching data from V2EX.
type Crawler struct {
	client *http.Client
}

// New creates a new Crawler with a default HTTP client.
func New() *Crawler {
	// http.ProxyFromEnvironment will automatically use the proxy
	// specified by the HTTP_PROXY or HTTPS_PROXY environment variables.
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	return &Crawler{
		client: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}
}

// FetchTopics fetches the list of topics from the V2EX homepage.
func (c *Crawler) FetchTopics() ([]model.Topic, error) {
	resp, err := c.client.Get(BaseURL + "/?tab=all")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var topics []model.Topic
	doc.Find(".cell.item").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".item_title a").Text()
		url, _ := s.Find(".item_title a").Attr("href")
		author := s.Find("strong a").First().Text()
		comments := s.Find(".count_livid").Text()
		if comments == "" {
			comments = "0"
		}

		// More robustly parse the time string.
		var timeStr string
		parts := strings.Split(s.Find(".topic_info").Text(), "•")
		if len(parts) > 2 {
			timeStr = strings.TrimSpace(parts[2])
		}

		if !strings.HasPrefix(url, "http") {
			url = BaseURL + url
		}

		topics = append(topics, model.Topic{
			Title:    title,
			Author:   author,
			Comments: comments,
			Time:     timeStr,
			URL:      url,
		})
	})

	return topics, nil
}

// FetchTopicDetail fetches the details of a single topic, including its replies.
func (c *Crawler) FetchTopicDetail(detailURL string) (*model.Topic, error) {
	resp, err := c.client.Get(detailURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	topic := &model.Topic{
		Title:   doc.Find("h1").Text(),
		Author:  doc.Find(".header small a").First().Text(),
		Time:    doc.Find(".header small span").First().AttrOr("title", ""),
		Content: doc.Find(".topic_content").Text(),
		URL:     detailURL,
	}

	var replies []model.Reply
	doc.Find(".cell[id^='r_']").Each(func(i int, s *goquery.Selection) {
		content := s.Find(".reply_content").Text()
		replyTo := ""
		// Extract reply-to user more reliably.
		if strings.HasPrefix(strings.TrimSpace(content), "@") {
			parts := strings.SplitN(content, " ", 2)
			if len(parts) > 0 {
				replyTo = parts[0][1:]
			}
		}

		reply := model.Reply{
			Author:  s.Find("strong a").Text(),
			Time:    s.Find(".ago").Text(),
			Content: content,
			Number:  s.Find(".no").Text(),
			ReplyTo: replyTo,
		}
		replies = append(replies, reply)
	})

	// This logic for calculating reply counts is inefficient (O(n^2))
	// but for a typical forum thread, it's acceptable.
	// A more optimized approach would use a map to store counts.
	for i := range replies {
		count := 0
		for _, r := range replies {
			if r.ReplyTo == replies[i].Author {
				count++
			}
		}
		replies[i].ReplyCount = count
	}

	topic.Replies = replies
	return topic, nil
}