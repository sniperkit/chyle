package decorators

// Config centralizes config needed for each decorator to being
// used by any third part package to make decorators work
type Config struct {
	GITHUBISSUE githubIssueConfig
	JIRAISSUE   jiraIssueConfig
	ENV         envConfig
}

// Features gives the informations if a decorator or several are defined
// and if so, which ones
type Features struct {
	ENABLED     bool
	JIRAISSUE   bool
	GITHUBISSUE bool
	ENV         bool
}