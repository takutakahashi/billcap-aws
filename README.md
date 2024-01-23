# billcap-aws
Billing capture and transfer for AWS Cost Explorer

## Prerequirements

### AWS Credentials
Only Access token is supported
### Data store
Only BigQuery is supported

### Execution Environment
  Only k8s is supported

## Usage

### 1. Set AWS Token

ENV `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`

### 2. Set Google Cloud Credentials

ENV `GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_PATH`

### 3. Deploy as CronJob

### 4. Check BigQuery data
