#!/bin/bash

set -e

printf "\n\nCreating SQS queues...\n"
awslocal sqs create-queue --queue-name proposals
awslocal sqs create-queue --queue-name risk-results
printf "\n\nSQS queues created successfully!\n"
