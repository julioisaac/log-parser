
#!/usr/bin/env bash

set -e

echo 'Starting Quake parser
                   _
                  | |
  __ _ _   _  __ _| | _____
 / _` | | | |/ _` | |/ / _ \
| (_| | |_| | (_| |   <  __/
 \__, |\__,_|\__,_|_|\_\___|
    | |
    |_|

'

COMMAND=${1:-"start"}
echo $COMMAND

case "$COMMAND" in
  baseCheck)
    exit 0
    ;;
  start)
    /usr/bin/log-parser
    ;;
  *)
    exec sh -c "$*"
    ;;
esac
