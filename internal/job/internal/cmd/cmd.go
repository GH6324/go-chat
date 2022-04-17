package cmd

import (
	"github.com/urfave/cli/v2"
	"go-chat/internal/job/internal/cmd/cron"
	"go-chat/internal/job/internal/cmd/other"
	"go-chat/internal/job/internal/cmd/queue"
	"go-chat/internal/pkg/cmdutil"
)

type Commands struct {
	CrontabCommand cron.Command
	QueueCommand   queue.Command
	OtherCommand   other.Command
}

func (c *Commands) Init() []*cli.Command {
	return cmdutil.ToSubCommand(c)
}
