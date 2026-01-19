#!/bin/sh
set -e

awslocal sqs create-queue --queue-name proposals >/dev/null
