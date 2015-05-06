function run_tests {
####################################################
# Tests go here.
####################################################

  if test -f Dockme.yml; then
    echo "Backing up existing Dockme.yml to Dockme.bak"
    mv Dockme.yml .dockme.bak
  fi

  CMD_WO_IMG="go run bin/dockme.go --dryrun"
  CMD_W_IMG="go run bin/dockme.go --image jmervine/nodebox --dryrun"

  assert_grep "$CMD_W_IMG --help" "USAGE:" \
    "deplay usage with help"

  assert_grep "$CMD_WO_IMG" "USAGE:" \
    "show usage without image"

  assert_grep "$CMD_W_IMG --image jmervine/nodebox" \
    "+ docker run --hostname=nodebox --workdir=/src --rm --tty --interactive --volume=" \
    "show expected defaults"

  assert_grep "$CMD_W_IMG --bad" "^Incorrect Usage." \
    "warn on bad options"

  assert_grep "$CMD_W_IMG --workdir /foo" \
    " --workdir=/foo" \
    "verify arg with space"

  assert_grep "$CMD_W_IMG --workdir=/foo" \
    " --workdir=/foo" \
    "verify arg with equal"

  assert_grep "$CMD_W_IMG -T ruby --sudo" \
    "+ sudo docker run " \
    "verify sudo"

  assert_grep "$CMD_WO_IMG -T ruby" \
    "ruby:latest" \
    "verify template"

  assert_grep "$CMD_W_IMG echo hello" \
    "echo hello$" \
    "verify command"

  assert_grep "$CMD_W_IMG --save" \
    "Wrote Dockme.yml" \
    "verify save"

  assert_file "Dockme.yml" \
    "save file created"
  rm -f Dockme.yml
  if test -f .dockme.bak; then
    echo "Replacing existing Dockme.yml"
    mv .dockme.bak Dockme.yml
  fi


####################################################
}
