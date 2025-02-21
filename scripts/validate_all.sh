#!/usr/bin/env bash

schema="./.schema/entity.schema.json"
test_dir="./.schema/test/"

find "$test_dir" -name "*.hery" | while read -r file; do
  echo "Validating $file"
  if echo "$file" | grep -q "invalid.hery"; then
    if ./.script/hv.sh "$schema" "$file"; then
      echo "Error: $file is invalid but passed validation"
      exit 1
    else
      echo "$file correctly failed validation"
    fi
  elif echo "$file" | grep -q "valid.hery"; then
    if ./.script/hv.sh "$schema" "$file"; then
      echo "$file correctly passed validation"
    else
      echo "Error: $file is valid but failed validation"
      exit 1
    fi
  fi
done
