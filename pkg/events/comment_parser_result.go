package events

type CommentParseResult struct {
	Command            *CommentCommand
	CommentResponse    string
	Ignore             bool
	ImmediateResponse  bool
	HasResponseComment bool
}
