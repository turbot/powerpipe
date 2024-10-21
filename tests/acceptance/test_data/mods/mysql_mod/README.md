# duckdb_mod

### Description

This is a simple mod used for testing Powerpipe's ability to connect to a MySQL backend.

### Usage

This is a simple mod used for testing Powerpipe's ability to connect to a MySQL backend. This mod is also used to verify that passing params to a MySQL backend works as expected.

### Connection ###

#### Connect using mysql ####

Start your MySQL server and connect to it.
Create a table(something like this to test the queries):
mysql> SELECT * FROM employees;
+----+---------+------+---------+
| id | name    | age  | salary  |
+----+---------+------+---------+
|  1 | Alice   |   25 | 50000.5 |
|  2 | Bob     |   30 | 60000.8 |
|  3 | Charlie |   35 | 55000.2 |
+----+---------+------+---------+

#### Connect using powerpipe ####

Run the available query(total_employee):
```sh
$ powerpipe query run query.total_employee --database mysql://user:password@/dbname
```

Pass params to the query(total_employee):
```sh
$ powerpipe query run "query.total_employee(p1 => \"command_param_1\")" --database mysql://user:password@/dbname
```

Start dashboard server:
```sh
$ powerpipe server --database mysql://user:password@/dbname
```