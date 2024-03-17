#!/bin/bash

aws_region="us-east-1"
aws_account="000000000000"

account_topic="account-service-topic"
account_topic_arn="arn:aws:sns:$aws_region:$aws_account:$account_topic"

account_queue="account-service-queue"
account_queue_arn="arn:aws:sqs:$aws_region:$aws_account:$account_queue"

# creating topics
awslocal sns create-topic --name $account_topic

# creating queue
awslocal sqs create-queue --queue-name $account_queue

# subscribe queue on topic
awslocal sns subscribe --topic-arn $account_topic_arn --protocol sqs --notification-endpoint $account_queue_arn

echo \n
echo "###################################"
echo "AWS resources created successfully."
echo "###################################"
echo \n