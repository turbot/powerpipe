# sqlite_mod

### Description

This is a simple mod used for testing Powerpipe's ability to connect to a SQLite backend.

### Usage

This is a simple mod used for testing Powerpipe's ability to connect to a SQLite backend. The `chinook.db` is the database file. This mod is also used to verify that passing params to a SQLite backend works as expected.

### Connection ###

#### Connect using sqlite ####

```sh
$ sqlite /path/to/the/database/file/chinook.db 
```

#### Connect using powerpipe ####

Run the available query(albums):
```sh
$ powerpipe query run query.total_albums --database sqlite:////path/to/the/database/file/chinook.db
```

Pass params to the query(albums):
```sh
$ powerpipe query run "query.total_albums(p1 => \"command_param_1\")" --database sqlite:////path/to/the/database/file/chinook.db
```

Start dashboard server:
```sh
$ powerpipe server --database sqlite:////path/to/the/database/file/chinook.db
```