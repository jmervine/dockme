function run_tests {
####################################################
# Tests go here.
####################################################

  COMMAND="go run dockme.go --dryrun"

  assert_grep "$COMMAND --help" "USAGE:" \
    "deplay usage with help"

  assert_grep "$COMMAND" "+ docker run --name=dockme --hostname=zshrc --rm --tty --interactive --volume=" \
    "show expected defaults"

  assert_grep "$COMMAND --bad" "^Incorrect Usage." \
    "warn on bad options"

  assert_grep "$COMMAND --workdir /foo" \
    " --workdir=/foo" \
    "verify arg with space"

  assert_grep "$COMMAND --workdir=/foo" \
    " --workdir=/foo" \
    "verify arg with equal"

  assert_grep "$COMMAND -T ruby" \
    "ruby:latest" \
    "verify template"

  assert_grep "$COMMAND echo hello" \
    "echo hello$" \
    "verify commands"

  assert_grep "$COMMAND --save" \
    "# file: .dockmerc" \
    "verify save"

  assert_file ".dockmerc" \
    "save file created"
  rm -f .dockmerc

####################################################
}
