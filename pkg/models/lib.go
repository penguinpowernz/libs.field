package models

import "time"

type Lib struct {
	Name                  string    `json:"name"`
	FullName              string    `json:"full_name"`
	CurrentTag            string    `json:"current_tag"`
	TaggedAt              string    `json:"tagged_at"`
	ReleaseTag            string    `json:"release_tag"`
	ReleasedAt            string    `json:"released_at"`
	PushedAt              string    `json:"pushed_at"`
	Stargazers            int       `json:"stargazers"`
	PushesPerday          int       `json:"pushes_per_day"`
	StargazersChange      int       `json:"stargazers_change"`
	OpenIssuesCount       int       `json:"open_issues_count"`
	URL                   string    `json:"url"`
	Language              string    `json:"language"`
	License               string    `json:"license"`
	APIURL                string    `json:"api_url"`
	IsApp                 bool      `json:"is_app"`
	UpdatedAt               string    `json:"updated_at" redis:"UpdatedAt" zoom:"index"`
	UpdatedTime             time.Time `json:"updated_time" redis:"UpdatedTime"`
	ReleasesCheckedTime     time.Time `json:"releases_checked_time" redis:"ReleasesCheckedTime"`
	TagsCheckedTime         time.Time `json:"tags_checked_time" redis:"TagsCheckedTime"`
	ContributorsCheckedTime time.Time `json:"contributors_checked_time" redis:"ContributorsCheckedTime"`
	CommitsCheckedTime      time.Time `json:"commits_checked_time" redis:"CommitsCheckedTime"`
}

func (lib *Lib) ModelID() string {
	return lib.FullName
}

func (lib *Lib) SetModelID(id string) {
	lib.FullName = id
}

func NewLibFromRepo(repo GitHubRepo) Lib {
	return Lib{
		PushedAt:         repo.PushedAt,
		Stargazers:       repo.StargazersCount,
		PushesPerday:     0,
		StargazersChange: 0,
		OpenIssuesCount:  repo.OpenIssuesCount,
		URL:              repo.HTMLURL,
		License:          repo.License.SpdxID,
		Language:         repo.Language,
		Name:             repo.Name,
		FullName:         repo.FullName,
		UpdatedAt:        time.Now().Format(time.RFC3339),
		UpdatedTime:      time.Now(),
		APIURL:           repo.URL,
	}
}

func (lib *Lib) UpdateFromRepo(repo GitHubRepo) {
	lib.PushedAt = repo.PushedAt
	lib.StargazersChange = repo.StargazersCount - lib.Stargazers
	lib.Stargazers = repo.StargazersCount
	lib.OpenIssuesCount = repo.OpenIssuesCount
	lib.License = repo.License.SpdxID
	lib.UpdatedTime = time.Now()
	lib.UpdatedAt = time.Now().Format(time.RFC3339)
	lib.Name = repo.Name
	lib.FullName = repo.FullName
}
