#!/bin/bash

aws_region="us-east-1"
aws_account="000000000000"

account_topic="account-service-topic"
account_topic_arn="arn:aws:sns:$aws_region:$aws_account:$account_topic"

account_queue="account-service"
account_queue_arn="arn:aws:sqs:$aws_region:$aws_account:$account_queue"
account_queue_dlq="account-service-dlq"
account_queue_dlq_arn="arn:aws:sqs:$aws_region:$aws_account:$account_queue_dlq"

webhook_sender_queue="webhook-sender"
webhook_sender_queue_arn="arn:aws:sqs:$aws_region:$aws_account:$webhook_sender_queue"
webhook_sender_queue_dlq="webhook-sender-dlq"
webhook_sender_queue_dlq_arn="arn:aws:sqs:$aws_region:$aws_account:$webhook_sender_queue_dlq"

# creating topics
awslocal sns create-topic --name $account_topic

# creating queues
awslocal sqs create-queue --queue-name $account_queue_dlq
awslocal sqs create-queue --queue-name $account_queue \
    --attributes '{"RedrivePolicy": "{\"deadLetterTargetArn\":\"'$account_queue_dlq_arn'\",\"maxReceiveCount\":\"5\"}"}'

awslocal sqs create-queue --queue-name $webhook_sender_queue_dlq
awslocal sqs create-queue \
    --queue-name $webhook_sender_queue \
    --attributes '{"RedrivePolicy": "{\"deadLetterTargetArn\":\"'$webhook_sender_queue_dlq_arn'\",\"maxReceiveCount\":\"5\"}"}'

# subscribe queue on topic
awslocal sns subscribe \
    --topic-arn $account_topic_arn \
    --protocol sqs \
    --notification-endpoint $account_queue_arn

awslocal sns subscribe \
    --topic-arn $account_topic_arn \
    --protocol sqs \
    --notification-endpoint $webhook_sender_queue_arn


echo \n
echo "###################################"
echo "AWS resources created successfully."
echo "###################################"
echo \n
