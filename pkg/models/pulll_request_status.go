package models

import "time"

type PullRequestStatus struct {
	ApprovalStatus ApprovalStatus
	Mergeable      bool
}

type ApprovalStatus struct {
	IsApproved bool
	ApprovedBy string
	Date       time.Time
}
