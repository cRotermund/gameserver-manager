# Production Deployment

This overlay deploys `control-plane-api` and the Discord `control-bot` into the
`chad-applications` namespace. The web application remains outside this first
production deployment.

## Required SSM parameters

Create these AWS Systems Manager Parameter Store parameters as `SecureString`
values in `us-east-1`:

```text
/kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-access-key-id
/kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-secret-access-key
/kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-role-arn
/kubernetes/appsecrets/gameserver-manager/control-bot/discord-token
```

The names are within the `/kubernetes/appsecrets/` path allowed to the
External Secrets Operator read role. The three API parameters must be
`SecureString` values. The Discord token parameter must contain the bot token
as a `SecureString` value.

The API access key must be allowed to call `sts:AssumeRole` for the role stored
in `aws-role-arn`, and that role must grant `ec2:DescribeInstances`,
`ec2:StartInstances`, `ec2:StopInstances`, and `ec2:RebootInstances`. Scope
those actions to the managed game-server instances where practical. This is
the temporary deployment strategy; a dedicated workload identity role should
replace these static credentials in a later refactor.

## Image

The overlay currently references the `latest` image:

```text
ghcr.io/crotermund/gameserver-manager/control-plane-api:latest
ghcr.io/crotermund/gameserver-manager/control-bot:latest
```

Publish both images before syncing the Argo CD application. The GHCR pull
secret referenced by the Deployments must also exist in the target namespace.
