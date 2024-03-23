#!/bin/bash

aws_region="us-east-1"
aws_account="000000000000"

account_topic="account-service-topic"
account_topic_arn="arn:aws:sns:$aws_region:$aws_account:$account_topic"

account_queue="account-service"
account_queue_arn="arn:aws:sqs:$aws_region:$aws_account:$account_queue"
account_queue_dlq="account-service-dlq"
account_queue_dlq_arn="arn:aws:sqs:$aws_region:$aws_account:$account_queue_dlq"

webhook_schedule_queue="webhook-schedule"
webhook_schedule_queue_arn="arn:aws:sqs:$aws_region:$aws_account:$webhook_schedule_queue"
webhook_schedule_queue_dlq="webhook-schedule-dlq"
webhook_schedule_queue_dlq_arn="arn:aws:sqs:$aws_region:$aws_account:$webhook_schedule_queue_dlq"

webhook_sender_queue="webhook-sender"
webhook_sender_queue_arn="arn:aws:sqs:$aws_region:$aws_account:$webhook_sender_queue"
webhook_sender_queue_dlq="webhook-schedule-dlq"
webhook_sender_queue_dlq_arn="arn:aws:sqs:$aws_region:$aws_account:$webhook_sender_queue_dlq"

# creating topics
awslocal sns create-topic --name $account_topic

# creating queues
awslocal sqs create-queue --queue-name $account_queue_dlq
awslocal sqs create-queue --queue-name $account_queue \
    --attributes '{"RedrivePolicy": "{\"deadLetterTargetArn\":\"'$account_queue_dlq_arn'\",\"maxReceiveCount\":\"3\"}"}'

awslocal sqs create-queue --queue-name $webhook_schedule_queue_dlq
awslocal sqs create-queue \
    --queue-name $webhook_schedule_queue \
    --attributes '{"RedrivePolicy": "{\"deadLetterTargetArn\":\"'$webhook_schedule_queue_dlq_arn'\",\"maxReceiveCount\":\"3\"}"}'

awslocal sqs create-queue --queue-name $webhook_sender_queue_dlq
awslocal sqs create-queue \
    --queue-name $webhook_sender_queue \
    --attributes '{"RedrivePolicy": "{\"deadLetterTargetArn\":\"'$webhook_sender_queue_dlq_arn'\",\"maxReceiveCount\":\"3\"}"}'


# subscribe queue on topic
echo "### Subscribing $account_queue_arn to topics $account_topic_arn\n"
event_subscription_arn=$(awslocal sns subscribe \
    --topic-arn $account_topic_arn \
    --protocol sqs \
    --notification-endpoint $account_queue_arn \
    --output text)

awslocal sns set-subscription-attributes \
    --subscription-arn "$event_subscription_arn" \
    --attribute-name FilterPolicy \
    --attribute-value '{"EventType":[{"anything-but": ["event.created"]}]}'

echo "### Subscribing $webhook_schedule_queue_arn to topics $account_topic_arn\n"
webhook_schedule_subscription_arn=$(awslocal sns subscribe \
    --topic-arn $account_topic_arn \
    --protocol sqs \
    --notification-endpoint $webhook_schedule_queue_arn \
    --output text)

awslocal sns set-subscription-attributes \
    --subscription-arn "$webhook_schedule_subscription_arn" \
    --attribute-name FilterPolicy \
    --attribute-value '{"EventType":["event.created"]}'

echo \n
echo "###################################"
echo "AWS resources created successfully."
echo "###################################"
echo \n
