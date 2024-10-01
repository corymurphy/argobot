package models

type PullRequestComment struct {
	// Message string
	Sha  string
	Pull PullRequest
}

func NewPullRequestComment(sha string, pull PullRequest) PullRequestComment {
	return PullRequestComment{
		// Message: msg,
		Sha:  sha,
		Pull: pull,
	}
}

type CommentResponse struct {
	Message string
	Sha     string
	Pull    PullRequest
}

func NewCommentResponse(msg string, request PullRequestComment) CommentResponse {
	return CommentResponse{
		Message: msg,
		Sha:     request.Sha,
		Pull:    request.Pull,
	}
}
