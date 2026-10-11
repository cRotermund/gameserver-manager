# Production API Deployment

This overlay deploys only `control-plane-api` into the `chad-applications`
namespace. The Discord bot and web application remain outside this first
production deployment.

## Required SSM parameters

Create these AWS Systems Manager Parameter Store parameters as `SecureString`
values in `us-east-1`:

```text
/kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-access-key-id
/kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-secret-access-key
/kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-role-arn
```

The names are within the `/kubernetes/appsecrets/` path allowed to the
External Secrets Operator read role. The access key must be allowed to call
`sts:AssumeRole` for the role stored in `aws-role-arn`, and that role must
grant `ec2:DescribeInstances`, `ec2:StartInstances`, `ec2:StopInstances`, and
`ec2:RebootInstances`. Scope those actions to the managed game-server
instances where practical. This is the temporary deployment strategy; a
dedicated workload identity role should replace these static credentials in a
later refactor.

## Image

The overlay currently references the `latest` image:

```text
ghcr.io/crotermund/gameserver-manager/control-plane-api:latest
```

Publish that image before syncing the Argo CD application.
