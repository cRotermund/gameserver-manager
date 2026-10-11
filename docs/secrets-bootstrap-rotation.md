# Production Secret Bootstrap And Rotation SOP

This procedure provisions the runtime secrets used by the production API.
Secret values remain in a local, gitignored environment file and are written
to AWS Systems Manager Parameter Store as `SecureString` parameters.
Kubernetes reads those parameters through External Secrets Operator; values
are never committed to Git or placed in a Kubernetes manifest.

## Why This Exists

The production API needs AWS credentials to assume its game-server control
role. SSM is the source of truth for these values, while the application
repository contains only ExternalSecret references.

The harness writes the following parameters:

| Parameter | Purpose |
| --- | --- |
| `/kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-access-key-id` | API source credential access key |
| `/kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-secret-access-key` | API source credential secret key |
| `/kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-role-arn` | Role assumed by the API for EC2 operations |

All three values are stored as `SecureString` parameters. The AWS CLI identity
used to run the harness must be allowed to write these parameters and use the
configured KMS key, if a customer-managed key is used.

## Prerequisites

- AWS CLI is installed and authenticated with an operator identity allowed to write the approved SSM path.
- The AWS region is `us-east-1`, unless overridden with `--region`.
- The API access key can call `sts:AssumeRole` for the configured API role.
- The API role grants the required EC2 actions for managed game-server instances.
- The GitOps `chad-applications` ExternalSecret configuration is installed and healthy.

The AWS CLI identity used by the operator is separate from the `AWS_ACCESS_KEY_ID`
and `AWS_SECRET_ACCESS_KEY` values stored for the API. Do not use the API
runtime credentials as the operator identity unless that is an intentional,
temporary bootstrap arrangement.

## Prepare The Local Secret File

Copy the example file and restrict access to it:

```bash
cp infra/k8s/prod/secrets.env.example infra/k8s/prod/secrets.env
chmod 600 infra/k8s/prod/secrets.env
```

Set the values in `infra/k8s/prod/secrets.env`:

```dotenv
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
AWS_ROLE_ARN=arn:aws:iam::123456789012:role/game-server-manager
```

On Windows, keep the file inside the repository but confirm it remains ignored
before running the command:

```powershell
git check-ignore infra/k8s/prod/secrets.env
```

Never paste these values into command-line arguments, commit them, or include
them in screenshots, logs, pull requests, or issue comments.

## Bootstrap

Bootstrap creates the parameters and refuses to overwrite an existing
parameter:

```bash
./harness.sh secrets bootstrap
```

Use an explicit AWS profile or region when needed:

```bash
./harness.sh secrets bootstrap \
  --profile personal \
  --region us-east-1
```

The harness prints parameter names only. Secret values are written through a
short-lived, permission-restricted AWS CLI JSON input file and are removed
after each write.

## Rotate

Update the local `secrets.env` file with the new API credential or role ARN,
then overwrite the existing SSM values:

```bash
./harness.sh secrets rotate
```

Rotation updates all three parameters from the file. The operation is
deterministic, but AWS SSM writes are not transactional; if a write fails,
inspect the reported parameter and rerun rotation after correcting the cause.
After the External Secrets refresh interval, verify the generated Kubernetes Secrets:

```bash
kubectl get externalsecret control-plane-api -n chad-applications
```

For an API credential rotation, restart the API Deployment after the
`control-plane-api` ExternalSecret reports `Ready=True` so the new environment
values are loaded:

```bash
kubectl rollout restart deployment/control-plane-api -n chad-applications
kubectl rollout status deployment/control-plane-api -n chad-applications
```

## Verification And Cleanup

Verify only metadata and Secret types. Do not print secret data:

```bash
aws ssm get-parameter \
  --name /kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-role-arn \
  --with-decryption \
  --query 'Parameter.{Name:Name,Type:Type,Version:Version}'
```

The SSM parameter metadata should report:

```text
SecureString
```

After a successful bootstrap or rotation, remove the local `secrets.env` file
when it is no longer needed. If it must be retained for future rotation, keep
it encrypted at rest and restrict filesystem access. Revoke superseded API
credentials immediately after confirming the replacement.

## Failure Handling

- If bootstrap reports that a parameter already exists, use `rotate` only after confirming the local file contains the intended replacement values.
- If AWS reports `AccessDenied`, check the operator IAM policy, region, parameter path, and KMS permissions.
- If an ExternalSecret is not ready, inspect its status and the External Secrets controller logs without printing generated Secret data.
- If the API cannot assume its role after rotation, verify the trust policy and the `AWS_ROLE_ARN` value.
