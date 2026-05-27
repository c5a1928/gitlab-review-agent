package di

import (
	"github.com/samber/do/v2"

	"github.com/antlss/gitlab-review-agent/internal/config"
	"github.com/antlss/gitlab-review-agent/internal/domain"
	"github.com/antlss/gitlab-review-agent/internal/pkg/git"
	"github.com/antlss/gitlab-review-agent/internal/pkg/gitlab"
	"github.com/antlss/gitlab-review-agent/internal/pkg/queue"
	"github.com/antlss/gitlab-review-agent/internal/pkg/store"
)

var InfraPackage = do.Package(
	do.Lazy(provideStores),
	do.Lazy(provideRepositorySettingsStore),
	do.Lazy(provideReviewJobStore),
	do.Lazy(provideReplyJobStore),
	do.Lazy(provideFeedbackStore),
	do.Lazy(provideReviewRecordStore),
	do.Lazy(provideGitLabClient),
	do.Lazy(provideGitManager),
	do.Lazy(provideJobQueue),
)

func provideStores(i do.Injector) (*store.Stores, error) {
	cfg := do.MustInvoke[*config.Config](i)
	return store.New(cfg.Store)
}

func provideRepositorySettingsStore(i do.Injector) (domain.RepositorySettingsStore, error) {
	stores := do.MustInvoke[*store.Stores](i)
	return stores.RepoSettings, nil
}

func provideReviewJobStore(i do.Injector) (domain.ReviewJobStore, error) {
	stores := do.MustInvoke[*store.Stores](i)
	return stores.ReviewJobs, nil
}

func provideReplyJobStore(i do.Injector) (domain.ReplyJobStore, error) {
	stores := do.MustInvoke[*store.Stores](i)
	return stores.ReplyJobs, nil
}

func provideFeedbackStore(i do.Injector) (domain.FeedbackStore, error) {
	stores := do.MustInvoke[*store.Stores](i)
	return stores.Feedbacks, nil
}

func provideReviewRecordStore(i do.Injector) (domain.ReviewRecordStore, error) {
	stores := do.MustInvoke[*store.Stores](i)
	return stores.ReviewRecords, nil
}

func provideGitLabClient(i do.Injector) (domain.GitLabClient, error) {
	cfg := do.MustInvoke[*config.Config](i)
	return gitlab.NewClient(cfg.GitLab.BaseURL, cfg.GitLab.Token), nil
}

func provideGitManager(i do.Injector) (*git.Manager, error) {
	cfg := do.MustInvoke[*config.Config](i)
	return git.NewManager(cfg.Git.ReposDir, cfg.GitLab.BaseURL, cfg.GitLab.Token), nil
}

func provideJobQueue(_ do.Injector) (*queue.Queue, error) {
	return queue.New(), nil
}
