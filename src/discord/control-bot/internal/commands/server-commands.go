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
	return func(ctx *messaging.Context) error {
		if len(ctx.Command.Args) < 1 {
			ctx.Client.Reply("Please include the name of the server, with syntax !serverstatus [server name]. (use !listservers to see your options)")
			return nil
		}

		serverName := ctx.Command.Args[0]

		summary, err := summaryFromName(serverName, gsm)
		if err != nil {
			return err
		}

		if summary == nil {
			ctx.Client.Reply(fmt.Sprintf("I could not find a server '%s', are you sure you spelled that right?", serverName))
			return nil
		}

		ctx.Client.Reply(fmt.Sprintf("Server: '%s' is: %s", serverName, summary.Status))

		return nil
	}
}

func StartServerCommand(gsm *controlplaneclient.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.Context) error {
		return handleStateChangeCommand("start", ctx, gsm)
	}
}

func StopServerCommand(gsm *controlplaneclient.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.Context) error {
		return handleStateChangeCommand("stop", ctx, gsm)
	}
}

func ListServersCommand(gsm *controlplaneclient.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.Context) error {

		summaryArr, err := gsm.ListServers(context.TODO())
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

		ctx.Client.Reply(sb.String())
		return nil
	}
}

func handleStateChangeCommand(statecmd string, ctx *messaging.Context, gsm *controlplaneclient.GsmApiClient) error {
	if len(ctx.Command.Args) < 1 {
		ctx.Client.Reply(
			"Please include the server name, with syntax !serverstatus [server name]. (use !listservers to see your options)")
		return nil
	}

	serverName := ctx.Command.Args[0]

	summary, err := summaryFromName(serverName, gsm)
	if err != nil {
		return err
	}

	if summary == nil {
		ctx.Client.Reply(
			fmt.Sprintf(
				"I could not find a server '%s', are you sure you spelled that right?",
				serverName))
		return nil
	}

	if statecmd == "start" {
		_, err = gsm.StartServer(context.TODO(), summary.ServerID)
	} else {
		_, err = gsm.StopServer(context.TODO(), summary.ServerID)
	}

	if err != nil {
		return err
	}

	ctx.Client.Reply(fmt.Sprintf("'%s' request issued, polling the server for appropriate status...", statecmd))

	desiredState := "running"
	if statecmd == "stop" {
		desiredState = "stopped"
	}

	achieved, lastStatus, err := pollForStatus(summary.ServerID, desiredState, gsm)
	if err != nil {
		return err
	}

	if !achieved {
		ctx.Client.Reply(
			fmt.Sprintf(
				"Timed out waiting for the desired status (%s) on server '%s'.  Last seen status: %s",
				desiredState,
				serverName,
				lastStatus))

		return nil
	}

	ctx.Client.Reply(fmt.Sprintf("Succes, server: '%s' is: %s", serverName, lastStatus))

	return nil
}

func summaryFromName(name string, gsm *controlplaneclient.GsmApiClient) (summary *controlplaneclient.ServerSummary, err error) {
	summaryArr, err := gsm.ListServers(context.TODO())
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

func pollForStatus(serverId string, goal string, gsm *controlplaneclient.GsmApiClient) (achieved bool, last string, err error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.After(ServerStatusPollTimeoutSeconds * time.Second)

	var lastStatus string

	for {
		select {
		case <-ticker.C:
			summary, err := gsm.GetServer(context.TODO(), serverId)

			if err != nil {
				return false, lastStatus, err
			}

			lastStatus = string(summary.Status)
			if strings.EqualFold(goal, string(summary.Status)) {
				return true, lastStatus, nil
			}

		case <-timeout:
			return false, lastStatus, nil
		}
	}
}
