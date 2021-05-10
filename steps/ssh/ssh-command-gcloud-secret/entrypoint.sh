#!/bin/bash

export AUTH_JSON=`echo $AUTH_CODE | base64 --decode`
echo $AUTH_JSON > /auth.json 

gcloud auth activate-service-account --key-file /auth.json --quiet --verbosity=none

SSH_KEY=`gcloud secrets versions access latest --secret $SSH_KEY_SECRET_NAME --project $GCLOUD_PROJECT`

echo -e "$SSH_KEY"  > /key.pem
chmod 400 /key.pem

ssh  -o "StrictHostKeyChecking=no" -o "LogLevel=ERROR" -i "/key.pem" -oConnectTimeout=$CONNECTION_TIMEOUT $SSH_CLIENT $SSH_COMMAND
