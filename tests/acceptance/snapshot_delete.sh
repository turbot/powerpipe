#!/bin/bash
#
# Functions the snapshot acceptance tests use to delete the snapshot they publish.
# Source this file; it runs nothing.
#
# Powerpipe prints a snapshot's URL as
#   https://<host>/org/<org>/workspace/<workspace>/powerpipe/snapshot/<id>
# but the Pipes API serves the snapshot at
#   https://<host>/api/v0/org/<org>/workspace/<workspace>/snapshot/<id>
# so the delete URL is built from the host, the org, the workspace and the id (the last
# path segment) of the printed URL, never from the printed path.

# Prints the API URL that deletes the snapshot a printed snapshot URL names.
# Exits 1, with a message on stderr, when the URL does not have that form.
snapshot_delete_url() {
  local printed="${1:-}"
  # a query string or fragment is not part of the path
  printed="${printed%%[?#]*}"
  local re='^(https?://[^/]+)/org/([^/]+)/workspace/([^/]+)(/[^/]+)*/snapshot/([^/]+)/?$'
  if [[ "$printed" =~ $re ]]; then
    echo "${BASH_REMATCH[1]}/api/v0/org/${BASH_REMATCH[2]}/workspace/${BASH_REMATCH[3]}/snapshot/${BASH_REMATCH[5]}"
  else
    echo "cannot build a snapshot delete URL from '${1:-}'" >&2
    return 1
  fi
}

# Deletes the snapshot a printed snapshot URL names, with the token in SPIPETOOLS_TOKEN.
# Returns 1, and says the delete URL and the HTTP status on stderr, unless the API answers 2xx.
snapshot_delete() {
  local delete_url status
  delete_url=$(snapshot_delete_url "${1:-}") || return 1
  # the token goes in a curl config on stdin, so it never shows in a process list
  status=$(printf 'header = "Authorization: Bearer %s"\n' "${SPIPETOOLS_TOKEN:-}" |
    curl -s -o /dev/null -w '%{http_code}' -X DELETE -K - "$delete_url") || true
  if [[ ! "$status" =~ ^2[0-9][0-9]$ ]]; then
    echo "DELETE $delete_url answered HTTP ${status:-000}, expected 2xx" >&2
    return 1
  fi
}
