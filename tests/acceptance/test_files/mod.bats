load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

function create_work_dir() {
  # create the work folder to run the tests
  tmp_dir="$(mktemp -d)"
  mkdir -p "${tmp_dir}"
}

function cleanup_work_dir() {
  # cleanup the work folder
  rm -rf "${tmp_dir}"
}

@test "generic mod install test" {
  # read the test-cases file
  tests=$(cat $FILE_PATH/test_data/source_files/mod_install.json)

  # fetch the keys(test names)
  test_keys=$(echo $tests | jq '. | keys[]')

  for i in $test_keys; do
    # create the work folder to run the tests
    create_work_dir
    # change to the temp work directory
    cd "$tmp_dir"

    name=$(echo $tests | jq -c ".[${i}]" | jq -r ".name")
    dir=$(echo $tests | jq -c ".[${i}]" | jq -r ".dir")
    cmd=$(echo $tests | jq -c ".[${i}]" | jq -r ".cmd")

    echo ""
    echo "Running test: $name"
    echo "Command: $cmd"
    echo "Mod directory: $FILE_PATH/test_data/mods/test_mods/$dir"
    echo ""

    # Copy the specified mod install dir to the temporary work directory
    cp -r "$FILE_PATH/test_data/mods/test_mods/$dir" "$tmp_dir"

    # switch to the copied mod directory
    cd "$tmp_dir/$dir"

    # run the powerpipe cmd
    echo "Command output:"
    run $cmd
    echo $output

    # check command output matches the expected output
    assert_output "$(cat $TEST_DATA_DIR/top_level_mod_upgraded.txt)"

    # check the folder structure matches the expected structure
    run ls .powerpipe/mods/github.com/pskrbasu/
    assert_output "$(cat $TEST_DATA_DIR/mod_folder_structure.txt)"

    # check the files match the expected
    run ls .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0
    assert_output "$(cat $TEST_DATA_DIR/mod_files.txt)"
  done
}