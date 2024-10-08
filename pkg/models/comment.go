package models

type PullRequestComment struct {
	// Message string
	Revision string
	Pull     PullRequest
}

func NewPullRequestComment(revision string, pull PullRequest) PullRequestComment {
	return PullRequestComment{
		// Message: msg,
		Revision: revision,
		Pull:     pull,
	}
}

type CommentResponse struct {
	Message  string
	Revision string
	Pull     PullRequest
}

func NewCommentResponse(msg string, request PullRequestComment) CommentResponse {
	return CommentResponse{
		Message:  msg,
		Revision: request.Revision,
		Pull:     request.Pull,
	}
}
