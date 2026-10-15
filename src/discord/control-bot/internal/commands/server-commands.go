package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
	"github.com/cRotermund/gameserver-manager/src/pkg/controlplaneclient"
)

const ServerStatusPollTimeoutSeconds = 30

func ServerStatusCommand(gsm *controlplaneclient.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.RoutingContext) error {
		if len(ctx.Command.Args) < 1 {
			return ctx.Client.Reply(
				ctx.Context,
				"Please include the name of the server, with syntax !serverstatus [server name]. (use !listservers to see your options)")
		}

		serverName := ctx.Command.Args[0]

		summary, err := summaryFromName(ctx.Context, serverName, gsm)
		if err != nil {
			return err
		}

		if summary == nil {
			return ctx.Client.Reply(
				ctx.Context,
				fmt.Sprintf(
					"I could not find a server '%s', are you sure you spelled that right?",
					serverName))
		}

		return ctx.Client.Reply(
			ctx.Context,
			fmt.Sprintf("Server: '%s' is: %s", serverName, summary.Status))
	}
}

func StartServerCommand(gsm *controlplaneclient.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.RoutingContext) error {
		return handleStateChangeCommand("start", ctx, gsm)
	}
}

func StopServerCommand(gsm *controlplaneclient.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.RoutingContext) error {
		return handleStateChangeCommand("stop", ctx, gsm)
	}
}

func ListServersCommand(gsm *controlplaneclient.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.RoutingContext) error {
		summaryArr, err := gsm.ListServers(ctx.Context)
		if err != nil {
			return err
		}

		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("There are %v servers available", len(summaryArr)))
		sb.WriteString("\n")
		sb.WriteString("--\n")

		for idx, summary := range summaryArr {
			sb.WriteString(summary.Name)

			if idx < len(summaryArr)-1 {
				sb.WriteString("\n")
			}
		}

		if err := ctx.Client.Reply(
			ctx.Context,
			sb.String()); err != nil {
			return err
		}

		return nil
	}
}

func handleStateChangeCommand(statecmd string, ctx *messaging.RoutingContext, gsm *controlplaneclient.GsmApiClient) error {
	if len(ctx.Command.Args) < 1 {
		return ctx.Client.Reply(
			ctx.Context,
			"Please include the server name, with syntax !serverstatus [server name]. (use !listservers to see your options)")
	}

	serverName := ctx.Command.Args[0]

	summary, err := summaryFromName(ctx.Context, serverName, gsm)
	if err != nil {
		return err
	}

	if summary == nil {
		return ctx.Client.Reply(
			ctx.Context,
			fmt.Sprintf(
				"I could not find a server '%s', are you sure you spelled that right?",
				serverName))
	}

	if statecmd == "start" {
		_, err = gsm.StartServer(ctx.Context, summary.ServerID)
	} else {
		_, err = gsm.StopServer(ctx.Context, summary.ServerID)
	}

	if err != nil {
		return err
	}

	if err := ctx.Client.Reply(
		ctx.Context,
		fmt.Sprintf(
			"'%s' request issued, polling the server for appropriate status...",
			statecmd)); err != nil {
		return err
	}

	desiredState := "running"
	if statecmd == "stop" {
		desiredState = "stopped"
	}

	achieved, lastStatus, err := pollForStatus(ctx.Context, summary.ServerID, desiredState, gsm)
	if err != nil {
		return err
	}

	if !achieved {
		if err := ctx.Client.Reply(
			ctx.Context,
			fmt.Sprintf(
				"Timed out waiting for the desired status (%s) on server '%s'.  Last seen status: %s",
				desiredState,
				serverName,
				lastStatus)); err != nil {
			return err
		}

		return nil
	}

	return ctx.Client.Reply(
		ctx.Context,
		fmt.Sprintf("Succes, server: '%s' is: %s",
			serverName,
			lastStatus))
}

func summaryFromName(ctx context.Context, name string, gsm *controlplaneclient.GsmApiClient) (summary *controlplaneclient.ServerSummary, err error) {
	summaryArr, err := gsm.ListServers(ctx)
	if err != nil {
		return nil, err
	}

	for _, summary := range summaryArr {
		if strings.EqualFold(summary.Name, name) {
			return &summary, nil
		}
	}

	return nil, nil
}

func pollForStatus(ctx context.Context, serverId string, goal string, gsm *controlplaneclient.GsmApiClient) (achieved bool, last string, err error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.NewTimer(ServerStatusPollTimeoutSeconds * time.Second)
	defer timeout.Stop()

	var lastStatus string

	for {
		select {
		case <-ctx.Done():
			return false, lastStatus, ctx.Err()

		case <-timeout.C:
			return false, lastStatus, nil

		case <-ticker.C:
			summary, err := gsm.GetServer(ctx, serverId)

			if err != nil {
				return false, lastStatus, err
			}

			lastStatus = string(summary.Status)
			if strings.EqualFold(goal, string(summary.Status)) {
				return true, lastStatus, nil
			}
		}
	}
}
