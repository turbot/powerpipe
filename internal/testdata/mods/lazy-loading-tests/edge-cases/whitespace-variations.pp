# Resources with different indentation patterns

# Compact style (minimal whitespace)
query "compact" {title="Compact Query"
sql="SELECT 1"}

# Expanded style (maximum whitespace)
query "expanded" {

  title       = "Expanded Query"

  description = "Query with extra whitespace"

  sql         = "SELECT 1"

  tags = {

    style = "expanded"

  }

}

# Mixed indentation
control "mixed_indent" {
	title = "Tab Indented"
    description = "Uses tabs and spaces"
	sql = "SELECT 'pass' as status"
}

# Trailing whitespace (intentional for testing)
query "trailing_whitespace" {
  title       = "Trailing Whitespace"
  description = "Has trailing spaces"
  sql         = "SELECT 1"
}

# Long single line
benchmark "single_line_benchmark" { title = "Single Line" children = [ control.mixed_indent ] }

# Deep indentation in arrays
benchmark "deep_array_indent" {
  title = "Deep Array Indent"
  children = [
                        control.mixed_indent
  ]
}

# No space around equals
query "no_space_equals" {
  title="No Space Around Equals"
  sql="SELECT 1 as value"
}

dashboard "whitespace_dashboard" {
  title = "Whitespace Variations"

  # Tightly packed container
  container{title="Tight"
    card{sql="SELECT 1"}
  }

  # Loosely packed container
  container {

    title = "Loose"

    card {

      sql = "SELECT 2"

    }

  }
}
