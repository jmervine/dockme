function run_tests {
####################################################
# Tests go here.
####################################################

  COMMAND="./dockme -D"

  assert_grep "$COMMAND --help" "Usage" \
    "deplay usage with help"

  assert_grep "$COMMAND" "+ docker run -it --workdir=/src --rm --volume=" \
    "show expected defaults"

  assert_grep "$COMMAND --bad" "WARN: Unknown option (ignored): --bad" \
    "warn on bad options"

  assert_grep "$COMMAND --workdir /foo" \
    "+ docker run -it --workdir=/foo --rm --volume=" \
    "verify arg with space"

  assert_grep "$COMMAND --workdir=/foo" \
    "+ docker run -it --workdir=/foo --rm --volume=" \
    "verify arg with equal"

  assert_grep "$COMMAND --no-rm" \
    "+ docker run -it --workdir=/src --volume=" \
    "verify no rm"

  assert_grep "$COMMAND --save" \
    "# file: .dockmerc" \
    "verify save"

  assert_file ".dockmerc" \
    "save file created"
  rm -f .dockmerc

####################################################
}
