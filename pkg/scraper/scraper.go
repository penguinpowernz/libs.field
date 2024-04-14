package scraper

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

var defaultInterval = 5 * time.Second

type Scraper struct {
	nc   *nats.Conn
	next <-chan time.Time
	cl   *http.Client
}

func New(nc *nats.Conn) *Scraper {
	return &Scraper{
		nc:   nc,
		next: time.After(time.Second),
		cl: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (s *Scraper) Run(ctx context.Context) error {
	in := make(chan *nats.Msg, 100)
	sub, err := s.nc.QueueSubscribeSyncWithChan("urls", "scraper", in)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			sub.Unsubscribe()
			return nil

		case msg := <-in:
			<-s.next
			dur := s.request(string(msg.Data))
			s.next = time.After(dur)

		}
	}
}

func (s *Scraper) request(url string) time.Duration {
	log.Printf("requesting url %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error creating request for url %s: %s", url, err)
		// s.queue = append(s.queue, url) // put the url back in the queue
		return defaultInterval
	}

	res, err := s.cl.Do(req)
	if err != nil {
		log.Printf("error making request for url %s: %s", url, err)
		// s.queue = append(s.queue, url) // put the url back in the queue
		return defaultInterval
	}

	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		log.Printf("got 200 for url %s", url)

		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("error reading response body for url %s: %s", url, err)
			// s.queue = append(s.queue, url) // put the url back in the queue
			return defaultInterval
		}

		if err := s.nc.Publish(subjFromURL(url), data); err != nil {
			log.Printf("error publishing response for url %s: %s", url, err)
			// s.queue = append(s.queue, url) // put the url back in the queue
			return defaultInterval
		}

		log.Printf("published %s response for url %s", subjFromURL(url), url)

	case 404:
		log.Printf("got 404 for url %s", url)
		return defaultInterval

	case 403, 429:
		log.Printf("got %d for url %s", res.StatusCode, url)
		// s.queue = append(s.queue, url) // put the url back in the queue

		if res.Header.Get("x-ratelimit-remaining") == "0" {
			epochS := res.Header.Get("x-ratelimit-reset")
			epoch, err := strconv.Atoi(epochS)
			if err != nil {
				log.Printf("error parsing x-ratelimit-reset header for url %s: %s", url, err)
				return time.Minute
			}

			log.Println("ratelimit reset in", time.Until(time.Unix(int64(epoch), 0)))
			return time.Until(time.Unix(int64(epoch), 0))
		}

		log.Println("couldn't determine ratelimit, sleeping for a minute")
		return time.Minute
	}

	return defaultInterval
}

func subjFromURL(url string) string {
	parts := strings.Split(url, "/")
	last := parts[len(parts)-1]

	if strings.Contains(url, "/search/repositories?") {
		return "search"
	}

	switch last {
	case "tags":
		return "tags"
	case "releases":
		return "releases"
	case "commits":
		return "commits"
	case "contributors":
		return "contributors"
	}

	if len(parts) == 6 {
		return "repos"
	}

	log.Printf("unknown url %s", url)
	return "unknown"
}
