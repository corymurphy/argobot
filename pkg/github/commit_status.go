package github

type CommitState int

const (
	PendingCommitState CommitState = iota
	SuccessCommitState
	FailedCommitState
)

func (s CommitState) String() string {
	switch s {
	case PendingCommitState:
		return "pending"
	case SuccessCommitState:
		return "success"
	case FailedCommitState:
		return "failure"
	}
	return "failure"
}
