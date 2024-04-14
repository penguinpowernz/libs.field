package parsers

import (
	"testing"

	"github.com/albrow/zoom"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewRepo(t *testing.T) {
	var saved zoom.Model
	var findCalled bool

	rp := repoParser{
		find:    func(string, zoom.Model) error { findCalled = true; return nil },
		save:    func(m zoom.Model) error { saved = m; return nil },
		publish: func(string, []byte) error { return nil },
		exists:  func(string) (bool, error) { return false, nil },
	}

	rp.parse([]byte(repoBody))

	assert.NotNil(t, saved)
	assert.False(t, findCalled)
	assert.IsType(t, &models.Lib{}, saved)
	assert.Equal(t, "gin-gonic/gin", saved.(*models.Lib).FullName)
	assert.NotNil(t, rp.lib)
}

func TestIgnoreForkedRepo(t *testing.T) {
	var called bool
	rp := repoParser{
		exists: func(string) (bool, error) { called = true; return false, nil },
	}

	rp.parse([]byte(`{"id": 20904437, "fork": true}`))
	assert.False(t, called)
}

func TestUpdateRepo(t *testing.T) {
	var findID string
	var saveCalled bool

	rp := repoParser{
		find: func(id string, m zoom.Model) error {
			findID = id
			*m.(*models.Lib) = models.Lib{FullName: "gin-gonic/gin"}
			return nil
		},
		save:    func(m zoom.Model) error { saveCalled = true; return nil },
		publish: func(string, []byte) error { return nil },
		exists:  func(string) (bool, error) { return true, nil },
	}

	rp.parse([]byte(repoBody))

	assert.Equal(t, "gin-gonic/gin", findID)
	assert.True(t, saveCalled)
	assert.IsType(t, &models.Lib{}, rp.lib)
	assert.NotNil(t, rp.lib)
}

func TestNewRepoPublishAll(t *testing.T) {
	var pubCount int

	rp := repoParser{
		find:    func(string, zoom.Model) error { return nil },
		save:    func(m zoom.Model) error { return nil },
		publish: func(string, []byte) error { pubCount++; return nil },
		exists:  func(string) (bool, error) { return false, nil },
	}

	rp.parse([]byte(repoBody))

	assert.Equal(t, 5, pubCount)
}

var repoBody = `
{
      "id": 20904437,
      "node_id": "MDEwOlJlcG9zaXRvcnkyMDkwNDQzNw==",
      "name": "gin",
      "full_name": "gin-gonic/gin",
      "private": false,
      "owner": {
        "login": "gin-gonic",
        "id": 7894478,
        "node_id": "MDEyOk9yZ2FuaXphdGlvbjc4OTQ0Nzg=",
        "avatar_url": "https://avatars.githubusercontent.com/u/7894478?v=4",
        "gravatar_id": "",
        "url": "https://api.github.com/users/gin-gonic",
        "html_url": "https://github.com/gin-gonic",
        "followers_url": "https://api.github.com/users/gin-gonic/followers",
        "following_url": "https://api.github.com/users/gin-gonic/following{/other_user}",
        "gists_url": "https://api.github.com/users/gin-gonic/gists{/gist_id}",
        "starred_url": "https://api.github.com/users/gin-gonic/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/gin-gonic/subscriptions",
        "organizations_url": "https://api.github.com/users/gin-gonic/orgs",
        "repos_url": "https://api.github.com/users/gin-gonic/repos",
        "events_url": "https://api.github.com/users/gin-gonic/events{/privacy}",
        "received_events_url": "https://api.github.com/users/gin-gonic/received_events",
        "type": "Organization",
        "site_admin": false
      },
      "html_url": "https://github.com/gin-gonic/gin",
      "description": "Gin is a HTTP web framework written in Go (Golang). It features a Martini-like API with much better performance -- up to 40 times faster. If you need smashing performance, get yourself some Gin.",
      "fork": false,
      "url": "https://api.github.com/repos/gin-gonic/gin",
      "forks_url": "https://api.github.com/repos/gin-gonic/gin/forks",
      "keys_url": "https://api.github.com/repos/gin-gonic/gin/keys{/key_id}",
      "collaborators_url": "https://api.github.com/repos/gin-gonic/gin/collaborators{/collaborator}",
      "teams_url": "https://api.github.com/repos/gin-gonic/gin/teams",
      "hooks_url": "https://api.github.com/repos/gin-gonic/gin/hooks",
      "issue_events_url": "https://api.github.com/repos/gin-gonic/gin/issues/events{/number}",
      "events_url": "https://api.github.com/repos/gin-gonic/gin/events",
      "assignees_url": "https://api.github.com/repos/gin-gonic/gin/assignees{/user}",
      "branches_url": "https://api.github.com/repos/gin-gonic/gin/branches{/branch}",
      "tags_url": "https://api.github.com/repos/gin-gonic/gin/tags",
      "blobs_url": "https://api.github.com/repos/gin-gonic/gin/git/blobs{/sha}",
      "git_tags_url": "https://api.github.com/repos/gin-gonic/gin/git/tags{/sha}",
      "git_refs_url": "https://api.github.com/repos/gin-gonic/gin/git/refs{/sha}",
      "trees_url": "https://api.github.com/repos/gin-gonic/gin/git/trees{/sha}",
      "statuses_url": "https://api.github.com/repos/gin-gonic/gin/statuses/{sha}",
      "languages_url": "https://api.github.com/repos/gin-gonic/gin/languages",
      "stargazers_url": "https://api.github.com/repos/gin-gonic/gin/stargazers",
      "contributors_url": "https://api.github.com/repos/gin-gonic/gin/contributors",
      "subscribers_url": "https://api.github.com/repos/gin-gonic/gin/subscribers",
      "subscription_url": "https://api.github.com/repos/gin-gonic/gin/subscription",
      "commits_url": "https://api.github.com/repos/gin-gonic/gin/commits{/sha}",
      "git_commits_url": "https://api.github.com/repos/gin-gonic/gin/git/commits{/sha}",
      "comments_url": "https://api.github.com/repos/gin-gonic/gin/comments{/number}",
      "issue_comment_url": "https://api.github.com/repos/gin-gonic/gin/issues/comments{/number}",
      "contents_url": "https://api.github.com/repos/gin-gonic/gin/contents/{+path}",
      "compare_url": "https://api.github.com/repos/gin-gonic/gin/compare/{base}...{head}",
      "merges_url": "https://api.github.com/repos/gin-gonic/gin/merges",
      "archive_url": "https://api.github.com/repos/gin-gonic/gin/{archive_format}{/ref}",
      "downloads_url": "https://api.github.com/repos/gin-gonic/gin/downloads",
      "issues_url": "https://api.github.com/repos/gin-gonic/gin/issues{/number}",
      "pulls_url": "https://api.github.com/repos/gin-gonic/gin/pulls{/number}",
      "milestones_url": "https://api.github.com/repos/gin-gonic/gin/milestones{/number}",
      "notifications_url": "https://api.github.com/repos/gin-gonic/gin/notifications{?since,all,participating}",
      "labels_url": "https://api.github.com/repos/gin-gonic/gin/labels{/name}",
      "releases_url": "https://api.github.com/repos/gin-gonic/gin/releases{/id}",
      "deployments_url": "https://api.github.com/repos/gin-gonic/gin/deployments",
      "created_at": "2014-06-16T23:57:25Z",
      "updated_at": "2024-04-11T12:08:17Z",
      "pushed_at": "2024-04-08T22:48:33Z",
      "git_url": "git://github.com/gin-gonic/gin.git",
      "ssh_url": "git@github.com:gin-gonic/gin.git",
      "clone_url": "https://github.com/gin-gonic/gin.git",
      "svn_url": "https://github.com/gin-gonic/gin",
      "homepage": "https://gin-gonic.com/",
      "size": 3223,
      "stargazers_count": 75253,
      "watchers_count": 75253,
      "language": "Go",
      "has_issues": true,
      "has_projects": false,
      "has_downloads": true,
      "has_wiki": false,
      "has_pages": false,
      "has_discussions": false,
      "forks_count": 7809,
      "mirror_url": null,
      "archived": false,
      "disabled": false,
      "open_issues_count": 766,
      "license": {
        "key": "mit",
        "name": "MIT License",
        "spdx_id": "MIT",
        "url": "https://api.github.com/licenses/mit",
        "node_id": "MDc6TGljZW5zZTEz"
      },
      "allow_forking": true,
      "is_template": false,
      "web_commit_signoff_required": false,
      "topics": [
        "framework",
        "gin",
        "go",
        "middleware",
        "performance",
        "router",
        "server"
      ],
      "visibility": "public",
      "forks": 7809,
      "open_issues": 766,
      "watchers": 75253,
      "default_branch": "master",
      "score": 1.0
    }
`
