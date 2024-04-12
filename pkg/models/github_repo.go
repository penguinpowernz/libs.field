package models

// GitHubRepo is a struct representing a GitHub repository (e.g. https://api.github.com/repos/gin-gonic/gin)
type GitHubRepo struct {
	ID              int    `json:"id"`                // 20904437
	NodeID          string `json:"node_id"`           // MDEwOlJlcG9zaXRvcnkyMDkwNDQzNw==
	Name            string `json:"name"`              // gin
	FullName        string `json:"full_name"`         // gin-gonic/gin
	HTMLURL         string `json:"html_url"`          // https://github.com/gin-gonic/gin
	Description     string `json:"description"`       // Gin is a HTTP web framework written in Go (Golang). It features a Martini-like API with much better performance -- up to 40 times faster. If you need smashing performance, get yourself some Gin.
	Fork            bool   `json:"fork"`              // false
	URL             string `json:"url"`               // https://api.github.com/repos/gin-gonic/gin
	EventsURL       string `json:"events_url"`        // https://api.github.com/repos/gin-gonic/gin/events
	TagsURL         string `json:"tags_url"`          // https://api.github.com/repos/gin-gonic/gin/tags
	ContributorsURL string `json:"contributors_url"`  // https://api.github.com/repos/gin-gonic/gin/contributors
	CommitsURL      string `json:"commits_url"`       // https://api.github.com/repos/gin-gonic/gin/commits{/sha}
	ReleasesURL     string `json:"releases_url"`      // https://api.github.com/repos/gin-gonic/gin/releases{/id}
	PushedAt        string `json:"pushed_at"`         // 2024-04-08T22:48:33Z
	StargazersCount int    `json:"stargazers_count"`  // 75253
	WatchersCount   int    `json:"watchers_count"`    // 75253
	Language        string `json:"language"`          // Go
	Archived        bool   `json:"archived"`          // false
	OpenIssuesCount int    `json:"open_issues_count"` // 766
	License         struct {
		Key    string `json:"key"`     // mit
		Name   string `json:"name"`    // MIT License
		SpdxID string `json:"spdx_id"` // MIT
		URL    string `json:"url"`     // https://api.github.com/licenses/mit
		NodeID string `json:"node_id"` // MDc6TGljZW5zZTEz
	} `json:"license"`
	Topics           []string `json:"topics"`            // [framework gin go middleware performance router server]
	Forks            int      `json:"forks"`             // 7809
	OpenIssues       int      `json:"open_issues"`       // 766
	Watchers         int      `json:"watchers"`          // 75253
	NetworkCount     int      `json:"network_count"`     // 7809
	SubscribersCount int      `json:"subscribers_count"` // 1358
}
