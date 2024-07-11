#!/bin/bash

if [ -z "$1" ]; then
  echo "Please provide a password."
  exit 1
fi

if ! command -v md5sum &> /dev/null; then
  echo "md5sum is not installed. Please install it."
  exit 1
fi

if ! command -v htpasswd &> /dev/null; then
  echo "htpasswd is not installed. Please install it."
  exit 1
fi

hash=$(echo -n $1 | md5sum | cut -d ' ' -f 1)

echo $(htpasswd -bnBC 10 "" $hash | tr -d ':\n' | sed 's/$2y/$2a/')