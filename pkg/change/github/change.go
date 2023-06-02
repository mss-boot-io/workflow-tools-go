/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/8/29 10:44:56
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/8/29 10:44:56
 */

package github

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-github/v44/github"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"

	"github.com/mss-boot-io/workflow-tools/pkg/change"
)

type Github struct {
	token string
	owner string
	repo  string
}

func (e *Github) SetAuth(auth interface{}) {
	switch auth := auth.(type) {
	case string:
		e.token = auth
	default:
		return
	}
}

// SetRepoURL e.g. WhiteMatrixTech/matrix-cloud-monorepo
func (e *Github) SetRepoURL(repo string) error {
	u, err := url.Parse(repo)
	if err != nil {
		log.Printf("SetRepoURL error, %s", err.Error())
		return err
	}
	paths := strings.Split(u.Path, "/")
	if len(paths) < 2 {
		return errors.New("repo url illegal")
	}
	e.owner = paths[0]
	e.repo = paths[1]
	return nil
}

// ChangeFiles get change files
func (e *Github) ChangeFiles(mark string) (*change.Files, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: e.token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	number, err1 := cast.ToIntE(mark)
	if err1 != nil {
		//push
		return e.getAllCommitFiles(ctx, client, mark)
	}
	// pr
	return e.getAllPRFiles(ctx, client, number)

}

func (e *Github) getAllPRFiles(ctx context.Context, client *github.Client, number int) (*change.Files, error) {
	var page int
	files := &change.Files{
		Added:    make([]string, 0),
		Modified: make([]string, 0),
		Deleted:  make([]string, 0),
		Renamed:  make([]string, 0),
	}
	for {
		list, resp, err := client.PullRequests.ListFiles(ctx, e.owner, e.repo, number, &github.ListOptions{
			Page: page,
		})
		if err != nil {
			return nil, err
		}

		filesCategory(list, files)

		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}
	return files, nil
}

func filesCategory(commitFiles []*github.CommitFile, files *change.Files) {
	if files == nil {
		files = &change.Files{
			Added:    make([]string, 0),
			Modified: make([]string, 0),
			Deleted:  make([]string, 0),
			Renamed:  make([]string, 0),
		}
	}
	for i := range commitFiles {
		switch strings.ToLower(commitFiles[i].GetStatus()) {
		case "added":
			files.Added = append(files.Added, commitFiles[i].GetFilename())
		case "modified":
			files.Modified = append(files.Modified, commitFiles[i].GetFilename())
		case "deleted":
			files.Deleted = append(files.Deleted, commitFiles[i].GetFilename())
		case "renamed":
			files.Renamed = append(files.Renamed, commitFiles[i].GetFilename())
		}
	}
}

func (e *Github) getAllCommitFiles(ctx context.Context, client *github.Client, sha string) (*change.Files, error) {
	files := &change.Files{
		Added:    make([]string, 0),
		Modified: make([]string, 0),
		Deleted:  make([]string, 0),
		Renamed:  make([]string, 0),
	}
	var page int
	for {
		repoCommit, resp, err := client.Repositories.GetCommit(ctx, e.owner, e.repo, sha,
			&github.ListOptions{
				Page: page,
			})
		if err != nil {
			return nil, err
		}
		filesCategory(repoCommit.Files, files)
		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}
	return files, nil
}
