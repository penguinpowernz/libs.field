package models

import "time"

type Lib struct {
	Name                    string    `json:"name" redis:"Name"`
	FullName                string    `json:"full_name" redis:"FullName"`
	CurrentTag              string    `json:"current_tag" redis:"CurrentTag"`
	TaggedAt                string    `json:"tagged_at" redis:"TaggedAt" zoom:"index"`
	ReleaseTag              string    `json:"release_tag" redis:"ReleaseTag"`
	ReleasedAt              string    `json:"released_at" redis:"ReleasedAt" zoom:"index"`
	Description             string    `json:"description" redis:"Description"`
	PushedAt                string    `json:"pushed_at" redis:"PushedAt" zoom:"index"`
	UpdatedAt               string    `json:"updated_at" redis:"UpdatedAt" zoom:"index"`
	Stargazers              int       `json:"stargazers" redis:"Stargazers" zoom:"index"`
	PushesPerday            int       `json:"pushes_per_day" redis:"PushesPerday" zoom:"index"`
	StargazersChange        int       `json:"stargazers_change" redis:"StargazersChange" zoom:"index"`
	OpenIssuesCount         int       `json:"open_issues_count" redis:"OpenIssuesCount"`
	URL                     string    `json:"url" redis:"URL"`
	Language                string    `json:"language" redis:"Language"`
	License                 string    `json:"license" redis:"License"`
	APIURL                  string    `json:"api_url" redis:"APIURL"`
	IsApp                   bool      `json:"is_app" redis:"IsApp"`
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

func NewLibFromRepo(repo GitHubRepo) *Lib {
	return &Lib{
		PushedAt:         repo.PushedAt,
		Stargazers:       repo.StargazersCount,
		PushesPerday:     0,
		StargazersChange: 0,
		OpenIssuesCount:  repo.OpenIssuesCount,
		Description:      repo.Description,
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
	lib.Description = repo.Description
	lib.License = repo.License.SpdxID
	lib.UpdatedTime = time.Now()
	lib.UpdatedAt = time.Now().Format(time.RFC3339)
	lib.Name = repo.Name
	lib.FullName = repo.FullName
}
