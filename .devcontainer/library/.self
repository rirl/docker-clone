#!/usr/bin/env bash
###
#### Minimal script template.
### 
function self() {
  export _self="$(basename -- $0 .sh)"
  export _self_path="$( cd "$( dirname "${BASH_SOURCE[0]}")/" >/dev/null 2>&1 && pwd )"
}
status=1; self
echo "[$_self] --  Command: BEGIN:[$_self_path]"
echo "[$_self] --  Command: END  :[$status]"
exit $status