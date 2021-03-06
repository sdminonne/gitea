// Copyright 2019 The Gitea Authors. All rights reserved.
// Copyright 2018 Jonas Franz. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"context"
	"io"
	"time"

	"code.gitea.io/gitea/modules/structs"
)

// AssetDownloader downloads an asset (attachment) for a release
type AssetDownloader interface {
	GetAsset(relTag string, relID, id int64) (io.ReadCloser, error)
}

// Downloader downloads the site repo informations
type Downloader interface {
	AssetDownloader
	SetContext(context.Context)
	GetRepoInfo() (*Repository, error)
	GetTopics() ([]string, error)
	GetMilestones() ([]*Milestone, error)
	GetReleases() ([]*Release, error)
	GetLabels() ([]*Label, error)
	GetIssues(page, perPage int) ([]*Issue, bool, error)
	GetComments(issueNumber int64) ([]*Comment, error)
	GetPullRequests(page, perPage int) ([]*PullRequest, bool, error)
	GetReviews(pullRequestNumber int64) ([]*Review, error)
}

// DownloaderFactory defines an interface to match a downloader implementation and create a downloader
type DownloaderFactory interface {
	New(ctx context.Context, opts MigrateOptions) (Downloader, error)
	GitServiceType() structs.GitServiceType
}

var (
	_ Downloader = &RetryDownloader{}
)

// RetryDownloader retry the downloads
type RetryDownloader struct {
	Downloader
	ctx        context.Context
	RetryTimes int // the total execute times
	RetryDelay int // time to delay seconds
}

// NewRetryDownloader creates a retry downloader
func NewRetryDownloader(ctx context.Context, downloader Downloader, retryTimes, retryDelay int) *RetryDownloader {
	return &RetryDownloader{
		Downloader: downloader,
		ctx:        ctx,
		RetryTimes: retryTimes,
		RetryDelay: retryDelay,
	}
}

// SetContext set context
func (d *RetryDownloader) SetContext(ctx context.Context) {
	d.ctx = ctx
	d.Downloader.SetContext(ctx)
}

// GetRepoInfo returns a repository information with retry
func (d *RetryDownloader) GetRepoInfo() (*Repository, error) {
	var (
		times = d.RetryTimes
		repo  *Repository
		err   error
	)
	for ; times > 0; times-- {
		if repo, err = d.Downloader.GetRepoInfo(); err == nil {
			return repo, nil
		}
		select {
		case <-d.ctx.Done():
			return nil, d.ctx.Err()
		case <-time.After(time.Second * time.Duration(d.RetryDelay)):
		}
	}
	return nil, err
}

// GetTopics returns a repository's topics with retry
func (d *RetryDownloader) GetTopics() ([]string, error) {
	var (
		times  = d.RetryTimes
		topics []string
		err    error
	)
	for ; times > 0; times-- {
		if topics, err = d.Downloader.GetTopics(); err == nil {
			return topics, nil
		}
		select {
		case <-d.ctx.Done():
			return nil, d.ctx.Err()
		case <-time.After(time.Second * time.Duration(d.RetryDelay)):
		}
	}
	return nil, err
}

// GetMilestones returns a repository's milestones with retry
func (d *RetryDownloader) GetMilestones() ([]*Milestone, error) {
	var (
		times      = d.RetryTimes
		milestones []*Milestone
		err        error
	)
	for ; times > 0; times-- {
		if milestones, err = d.Downloader.GetMilestones(); err == nil {
			return milestones, nil
		}
		select {
		case <-d.ctx.Done():
			return nil, d.ctx.Err()
		case <-time.After(time.Second * time.Duration(d.RetryDelay)):
		}
	}
	return nil, err
}

// GetReleases returns a repository's releases with retry
func (d *RetryDownloader) GetReleases() ([]*Release, error) {
	var (
		times    = d.RetryTimes
		releases []*Release
		err      error
	)
	for ; times > 0; times-- {
		if releases, err = d.Downloader.GetReleases(); err == nil {
			return releases, nil
		}
		select {
		case <-d.ctx.Done():
			return nil, d.ctx.Err()
		case <-time.After(time.Second * time.Duration(d.RetryDelay)):
		}
	}
	return nil, err
}

// GetLabels returns a repository's labels with retry
func (d *RetryDownloader) GetLabels() ([]*Label, error) {
	var (
		times  = d.RetryTimes
		labels []*Label
		err    error
	)
	for ; times > 0; times-- {
		if labels, err = d.Downloader.GetLabels(); err == nil {
			return labels, nil
		}
		select {
		case <-d.ctx.Done():
			return nil, d.ctx.Err()
		case <-time.After(time.Second * time.Duration(d.RetryDelay)):
		}
	}
	return nil, err
}

// GetIssues returns a repository's issues with retry
func (d *RetryDownloader) GetIssues(page, perPage int) ([]*Issue, bool, error) {
	var (
		times  = d.RetryTimes
		issues []*Issue
		isEnd  bool
		err    error
	)
	for ; times > 0; times-- {
		if issues, isEnd, err = d.Downloader.GetIssues(page, perPage); err == nil {
			return issues, isEnd, nil
		}
		select {
		case <-d.ctx.Done():
			return nil, false, d.ctx.Err()
		case <-time.After(time.Second * time.Duration(d.RetryDelay)):
		}
	}
	return nil, false, err
}

// GetComments returns a repository's comments with retry
func (d *RetryDownloader) GetComments(issueNumber int64) ([]*Comment, error) {
	var (
		times    = d.RetryTimes
		comments []*Comment
		err      error
	)
	for ; times > 0; times-- {
		if comments, err = d.Downloader.GetComments(issueNumber); err == nil {
			return comments, nil
		}
		select {
		case <-d.ctx.Done():
			return nil, d.ctx.Err()
		case <-time.After(time.Second * time.Duration(d.RetryDelay)):
		}
	}
	return nil, err
}

// GetPullRequests returns a repository's pull requests with retry
func (d *RetryDownloader) GetPullRequests(page, perPage int) ([]*PullRequest, bool, error) {
	var (
		times = d.RetryTimes
		prs   []*PullRequest
		err   error
		isEnd bool
	)
	for ; times > 0; times-- {
		if prs, isEnd, err = d.Downloader.GetPullRequests(page, perPage); err == nil {
			return prs, isEnd, nil
		}
		select {
		case <-d.ctx.Done():
			return nil, false, d.ctx.Err()
		case <-time.After(time.Second * time.Duration(d.RetryDelay)):
		}
	}
	return nil, false, err
}

// GetReviews returns pull requests reviews
func (d *RetryDownloader) GetReviews(pullRequestNumber int64) ([]*Review, error) {
	var (
		times   = d.RetryTimes
		reviews []*Review
		err     error
	)
	for ; times > 0; times-- {
		if reviews, err = d.Downloader.GetReviews(pullRequestNumber); err == nil {
			return reviews, nil
		}
		select {
		case <-d.ctx.Done():
			return nil, d.ctx.Err()
		case <-time.After(time.Second * time.Duration(d.RetryDelay)):
		}
	}
	return nil, err
}
