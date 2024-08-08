#!/usr/bin/env bash

hv() {
  local schema="${1}"
  local heryfile="${2}"
  local tmpfile="${heryfile%.hery}.yaml"

  cp "${heryfile}" "$tmpfile"
  jv "${schema}" "${tmpfile}"
  result=$?
  rm "${tmpfile}"
  return $result
}
hv "${@}"
