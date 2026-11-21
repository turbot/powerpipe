"""Helpers for server.bats: read the load state of a running `powerpipe server` and wait on it.

python3 -I workspace_load.py free-port
python3 -I workspace_load.py ctime <file>
python3 -I workspace_load.py wait-up <base-url>
python3 -I workspace_load.py state <base-url>
python3 -I workspace_load.py wait-loaded <base-url> <file>...
    Waits until the workspace is not loading and last_load.newest_file_time is at or after the newest change
    time (ctime) of the files, or of the folder of a file that is gone, then prints last_load as JSON. This is
    the wait a client does after it changes files. Exits 1 after 120 seconds.
python3 -I workspace_load.py probe-reload <base-url> <benchmark>
    Waits until the workspace is loading, then asks to start a run of the benchmark and prints the status
    code and the Retry-After header. Exits 1 if no reload starts within 60 seconds.
python3 -I workspace_load.py generate <file> <count>
    Writes <count> controls to <file>, so that a load takes long enough to be seen.
"""
import json
import os
import re
import socket
import sys
import time
import urllib.error
import urllib.request
from datetime import datetime


def state(base):
    with urllib.request.urlopen(base + "/api/v0/workspace", timeout=10) as r:
        return json.load(r)


def time_ns(value):
    # RFC 3339 with up to nine fraction digits, which datetime cannot parse
    head, _, tail = value.partition(".")
    fraction = re.match(r"\d*", tail).group(0)
    zone = tail[len(fraction):] if fraction else ""
    if not fraction:
        head, zone = re.match(r"(.*?)(Z|[+-]\d\d:\d\d)$", head).groups()
    zone = "+00:00" if zone == "Z" else zone
    return int(datetime.fromisoformat(head + zone).timestamp()) * 10**9 + int((fraction + "000000000")[:9])


def main(args):
    command = args[0]
    if command == "ctime":
        print(os.stat(args[1]).st_ctime_ns)
        return 0
    if command == "free-port":
        with socket.socket() as s:
            s.bind(("127.0.0.1", 0))
            print(s.getsockname()[1])
        return 0
    base = args[1]
    if command == "wait-up":
        end = time.monotonic() + 120
        while time.monotonic() < end:
            try:
                state(base)
                return 0
            except OSError:
                time.sleep(0.2)
        print("the server did not answer within 120 seconds", file=sys.stderr)
        return 1
    if command == "state":
        print(json.dumps(state(base)))
        return 0
    if command == "wait-loaded":
        # a file that is gone was deleted: its folder changed when it was
        goal = max(os.stat(f if os.path.exists(f) else os.path.dirname(os.path.abspath(f))).st_ctime_ns for f in args[2:])
        end = time.monotonic() + 120
        while time.monotonic() < end:
            s = state(base)
            last = s.get("last_load") or {}
            newest = last.get("newest_file_time")
            if not s["loading"] and newest and time_ns(newest) >= goal:
                print(json.dumps(dict(last, load_count=s["load_count"])))
                return 0
            time.sleep(0.1)
        print(f"the change did not load within 120 seconds: {json.dumps(state(base))}", file=sys.stderr)
        return 1
    if command == "probe-reload":
        end = time.monotonic() + 60
        while time.monotonic() < end:
            if state(base)["loading"]:
                req = urllib.request.Request(f"{base}/api/v0/benchmark/{args[2]}/run", data=b"{}", method="POST")
                req.add_header("Content-Type", "application/json")
                try:
                    with urllib.request.urlopen(req, timeout=10) as r:
                        print(r.status, r.headers.get("Retry-After", ""))
                except urllib.error.HTTPError as e:
                    print(e.code, e.headers.get("Retry-After", ""))
                return 0
            time.sleep(0.01)
        print("no reload started within 60 seconds", file=sys.stderr)
        return 1
    if command == "generate":
        with open(args[1], "w") as f:
            for i in range(int(args[2])):
                f.write(f'control "generated_{i}" {{\n  title = "Generated {i}"\n  sql   = "select 1"\n}}\n\n')
        return 0
    print(f"unknown command {command}", file=sys.stderr)
    return 2


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
