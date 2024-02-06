# duckdb_mod

### Description

This is a simple mod used for testing Powerpipe's ability to connect to a DuckDB backend.

### Usage

This is a simple mod used for testing Powerpipe's ability to connect to a DuckDB backend. The `employee.duckdb` is the database file. This mod is also used to verify that passing params to a DuckDB backend works as expected.

### Connection ###

#### Connect using duckdb ####

```sh
$ duckdb /path/to/the/database/file/employee.duckdb 
```

#### Connect using powerpipe ####

Run the available query(total_employee):
```sh
$ powerpipe query run query.total_employee --database duckdb:////path/to/the/database/file/employee.duckdb
```

Pass params to the query(total_employee):
```sh
$ powerpipe query run "query.total_employee(p1 => \"command_param_1\")" --database duckdb:////path/to/the/database/file/employee.duckdb
```

Start dashboard server:
```sh
$ powerpipe server --database duckdb:////path/to/the/database/file/employee.duckdb
```