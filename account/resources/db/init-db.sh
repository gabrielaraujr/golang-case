#!/bin/bash

set -e

CREATE="CREATE DATABASE account_proposals;"
SELECT="SELECT 1 FROM pg_database WHERE datname = 'account_proposals';"

printf "\n\nCreating database account_proposals...\n"
psql -U postgres -tc "${SELECT}" | grep -q 1 || psql -U postgres -c "${CREATE}"

printf "\n\nRunning migrations...\n"
psql -U postgres -d account_proposals -f /migrations/001_create_proposals_table.sql

printf "\n\nDatabase setup completed.\n"
