#!/bin/bash

DIRECTORY="/etc/nginx/conf.d"

# Command to run on file changes (e.g., reload Nginx)
COMMAND="nginx -s reload"

# Run fswatch and execute the command on changes
fswatch --event All --poll "$DIRECTORY" | while read event
do
  echo "Change detected: $event"
  $COMMAND
done

