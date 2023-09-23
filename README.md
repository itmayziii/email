# email-api

```shell
gcloud functions deploy SendEmail \
  --gen2 \
  --trigger-topic="send-email" \
  --runtime="go121" \
  --entry-point="SendEmail" \
  --region="us-central1" \
  --source="." \
  --ingress-settings="internal-only" \
  --no-allow-unauthenticated \
  --retry \
  --trigger-service-account="eventarc-trigger@itmayziii.iam.gserviceaccount.com" \
  --run-service-account="app-email-api@itmayziii.iam.gserviceaccount.com" \
  --service-account="app-email-api@itmayziii.iam.gserviceaccount.com" \
  --set-secrets="MG_API_KEY_MG_TOMMYMAY_DEV=mailgun-api-key-mg-tommymay-dev:latest"
  --set-env-vars="PROJECT_ID=itmayziii"
```
