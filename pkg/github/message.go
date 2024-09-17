package github

import (
	"crypto/hmac"
	"hash"
)

const (
	sha1Prefix            = "sha1"
	sha256Prefix          = "sha256"
	sha512Prefix          = "sha512"
	SHA1SignatureHeader   = "X-Hub-Signature"
	SHA256SignatureHeader = "X-Hub-Signature-256"
	EventTypeHeader       = "X-Github-Event"
	DeliveryIDHeader      = "X-Github-Delivery"
)

var (
	// eventTypeMapping maps webhooks types to their corresponding go-github struct types.
	eventTypeMapping = map[string]string{
		"branch_protection_rule":         "BranchProtectionRuleEvent",
		"check_run":                      "CheckRunEvent",
		"check_suite":                    "CheckSuiteEvent",
		"code_scanning_alert":            "CodeScanningAlertEvent",
		"commit_comment":                 "CommitCommentEvent",
		"content_reference":              "ContentReferenceEvent",
		"create":                         "CreateEvent",
		"delete":                         "DeleteEvent",
		"deploy_key":                     "DeployKeyEvent",
		"deployment":                     "DeploymentEvent",
		"deployment_status":              "DeploymentStatusEvent",
		"deployment_protection_rule":     "DeploymentProtectionRuleEvent",
		"discussion":                     "DiscussionEvent",
		"discussion_comment":             "DiscussionCommentEvent",
		"fork":                           "ForkEvent",
		"github_app_authorization":       "GitHubAppAuthorizationEvent",
		"gollum":                         "GollumEvent",
		"installation":                   "InstallationEvent",
		"installation_repositories":      "InstallationRepositoriesEvent",
		"installation_target":            "InstallationTargetEvent",
		"issue_comment":                  "IssueCommentEvent",
		"issues":                         "IssuesEvent",
		"label":                          "LabelEvent",
		"marketplace_purchase":           "MarketplacePurchaseEvent",
		"member":                         "MemberEvent",
		"membership":                     "MembershipEvent",
		"merge_group":                    "MergeGroupEvent",
		"meta":                           "MetaEvent",
		"milestone":                      "MilestoneEvent",
		"organization":                   "OrganizationEvent",
		"org_block":                      "OrgBlockEvent",
		"package":                        "PackageEvent",
		"page_build":                     "PageBuildEvent",
		"personal_access_token_request":  "PersonalAccessTokenRequestEvent",
		"ping":                           "PingEvent",
		"project":                        "ProjectEvent",
		"project_card":                   "ProjectCardEvent",
		"project_column":                 "ProjectColumnEvent",
		"public":                         "PublicEvent",
		"pull_request":                   "PullRequestEvent",
		"pull_request_review":            "PullRequestReviewEvent",
		"pull_request_review_comment":    "PullRequestReviewCommentEvent",
		"pull_request_review_thread":     "PullRequestReviewThreadEvent",
		"pull_request_target":            "PullRequestTargetEvent",
		"push":                           "PushEvent",
		"repository":                     "RepositoryEvent",
		"repository_dispatch":            "RepositoryDispatchEvent",
		"repository_import":              "RepositoryImportEvent",
		"repository_vulnerability_alert": "RepositoryVulnerabilityAlertEvent",
		"release":                        "ReleaseEvent",
		"secret_scanning_alert":          "SecretScanningAlertEvent",
		"security_advisory":              "SecurityAdvisoryEvent",
		"star":                           "StarEvent",
		"status":                         "StatusEvent",
		"team":                           "TeamEvent",
		"team_add":                       "TeamAddEvent",
		"user":                           "UserEvent",
		"watch":                          "WatchEvent",
		"workflow_dispatch":              "WorkflowDispatchEvent",
		"workflow_job":                   "WorkflowJobEvent",
		"workflow_run":                   "WorkflowRunEvent",
	}
)

func GenerateMAC(message, key []byte, hashFunc func() hash.Hash) []byte {
	mac := hmac.New(hashFunc, key)
	mac.Write(message)
	return mac.Sum(nil)
}
