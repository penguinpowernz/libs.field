package scraper

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type Scraper struct {
	nc    *nats.Conn
	next  <-chan time.Time
	queue []string
	cl    *http.Client
	qmu   *sync.Mutex
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
		qmu: &sync.Mutex{},
	}
}

func (s *Scraper) Run(ctx context.Context) error {
	in := make(chan *nats.Msg)
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
			s.qmu.Lock()
			s.queue = append(s.queue, string(msg.Data))
			s.qmu.Unlock()

		case <-s.next:
			dur := s.request()
			s.next = time.After(dur)

		}
	}
}

func (s *Scraper) request() time.Duration {
	s.qmu.Lock()
	defer s.qmu.Unlock()

	if len(s.queue) == 0 {
		return time.Second
	}

	url := s.queue[0]
	if len(s.queue) > 1 {
		s.queue = s.queue[1:]
	}
	s.queue = []string{}

	log.Printf("Scraper: requesting url %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Scraper: error creating request for url %s: %s", url, err)
		s.queue = append(s.queue, url) // put the url back in the queue
		return time.Second
	}

	res, err := s.cl.Do(req)
	if err != nil {
		log.Printf("Scraper: error making request for url %s: %s", url, err)
		s.queue = append(s.queue, url) // put the url back in the queue
		return time.Second
	}

	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		log.Printf("Scraper: got 200 for url %s", url)

		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Scraper: error reading response body for url %s: %s", url, err)
			s.queue = append(s.queue, url) // put the url back in the queue
			return time.Second
		}

		if err := s.nc.Publish(subjFromURL(url), data); err != nil {
			log.Printf("Scraper: error publishing response for url %s: %s", url, err)
			s.queue = append(s.queue, url) // put the url back in the queue
			return time.Second
		}

	case 404:
		log.Printf("Scraper: got 404 for url %s", url)
		return time.Second

	case 403, 429:
		log.Printf("Scraper: got %d for url %s", res.StatusCode, url)
		s.queue = append(s.queue, url) // put the url back in the queue

		if res.Header.Get("x-ratelimit-remaining") == "0" {
			epochS := res.Header.Get("x-ratelimit-reset")
			epoch, err := strconv.Atoi(epochS)
			if err != nil {
				log.Printf("Scraper: error parsing x-ratelimit-reset header for url %s: %s", url, err)
				return time.Minute
			}

			return time.Until(time.Unix(int64(epoch), 0))
		}

		return time.Minute
	}

	return time.Second
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
	}

	if len(parts) == 6 {
		return "repos"
	}

	log.Printf("Scraper: unknown url %s", url)
	return "unknown"
}
