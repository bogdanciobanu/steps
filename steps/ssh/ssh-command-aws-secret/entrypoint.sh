#!/bin/bash


SSH_KEY=`aws secretsmanager get-secret-value --secret-id $SSH_KEY_SECRET_NAME | jq .SecretString | tr -d '"'`
echo -e "$SSH_KEY"  > /key.pem
chmod 400 /key.pem

ssh  -o "StrictHostKeyChecking=no" -o "LogLevel=ERROR" -i "/key.pem" -oConnectTimeout=$CONNECTION_TIMEOUT $SSH_CLIENT $SSH_COMMAND
