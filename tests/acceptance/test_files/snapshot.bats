load "$LIB_BATS_ASSERT/load.bash"
load "$LIB_BATS_SUPPORT/load.bash"

# These set of tests are skipped locally
# To run these tests locally set the SPIPETOOLS_TOKEN env var.
# These tests will be skipped locally unless the below env var is set.

function setup() {
  if [[ -z "${SPIPETOOLS_TOKEN}" ]]; then
    skip
  fi
  req_url=""
}

# Runs after every test, including failing ones. Deleting the snapshot inline at
# the end of a test only happened when every assertion passed, so any failure
# after the upload left the snapshot behind in the workspace for good.
function teardown() {
  rm -f output.*

  if [[ -n "${req_url}" ]]; then
    curl -s -X DELETE "$req_url" -H "Authorization: Bearer $SPIPETOOLS_TOKEN" || true
  fi
}

@test "snapshot mode - query output csv" {
  cd $FILE_PATH/test_data/mods/functionality_test_mod

  powerpipe query run static_query2 --snapshot --output csv --pipes-token $SPIPETOOLS_TOKEN --snapshot-location turbot-ops/clitesting > output.csv

  # extract the snapshot url from the output
  url=$(grep -o 'http[^"]*' output.csv)
  echo $url

  # create the snapshot DELETE Request URL, so teardown can remove it whatever
  # the assertions below do
  req_url=$($FILE_PATH/url_parse.sh $url)
  echo $req_url

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 15th line, since it contains snapshot upload link, which will be different in each run
  if [[ "$OSTYPE" == "darwin"* ]]; then
    run sed -i ".csv" "15d" output.csv
  else
    run sed -i "15d" output.csv
  fi
  cat output.csv

  assert_equal "$(cat output.csv)" "$(cat $TEST_DATA_DIR/expected_static_query_csv_snapshot_mode.csv)"
}

@test "snapshot mode - query output json" {
  skip "cannot compare json output as it contains duration."
  cd $FILE_PATH/test_data/mods/functionality_test_mod

  powerpipe query run static_query2 --snapshot --output json --pipes-token $SPIPETOOLS_TOKEN --snapshot-location turbot-ops/clitesting > output.json

  # extract the snapshot url from the output
  url=$(grep -o 'http[^"]*' output.json)
  echo $url

  # create the snapshot DELETE Request URL, so teardown can remove it whatever
  # the assertions below do
  req_url=$($FILE_PATH/url_parse.sh $url)
  echo $req_url

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 64th line, since it contains snapshot upload link, which will be different in each run
  if [[ "$OSTYPE" == "darwin"* ]]; then
    run sed -i ".csv" "64d" output.json
  else
    run sed -i "64d" output.json
  fi
  cat output.json

  assert_equal "$(cat output.json)" "$(cat $TEST_DATA_DIR/expected_static_query_json_snapshot_mode.json)"
}

@test "snapshot mode - query output table" {
  cd $FILE_PATH/test_data/mods/functionality_test_mod

  powerpipe query run static_query2 --snapshot --output table --pipes-token $SPIPETOOLS_TOKEN --snapshot-location turbot-ops/clitesting > output.txt

  # extract the snapshot url from the output
  url=$(grep -o 'http[^"]*' output.txt)
  echo $url

  # create the snapshot DELETE Request URL, so teardown can remove it whatever
  # the assertions below do
  req_url=$($FILE_PATH/url_parse.sh $url)
  echo $req_url

  # checking for OS type, since sed command is different for linux and OSX
  # removing the 18th line, since it contains snapshot upload link, which will be different in each run
  if [[ "$OSTYPE" == "darwin"* ]]; then
    run sed -i ".csv" "18d" output.txt
  else
    run sed -i "18d" output.txt
  fi
  cat output.txt

  assert_equal "$(cat output.txt)" "$(cat $TEST_DATA_DIR/expected_static_query_table_snapshot_mode.txt)"
}
