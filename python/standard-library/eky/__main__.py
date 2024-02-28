import argparse
import json
from pathlib import Path
from typing import Any


class Eky:
    def __init__(self) -> None:
        self.path = Path.home() / ".eky.json1"
        self.data: dict[str, Any] = {}
        self._load_file()

    def list(self) -> None:
        for key in self.data.keys():
            print(key)

    def __getitem__(self, key: str) -> Any:
        val = self.data.get(key, None)
        if val:
            return json.dumps(val, indent=2)

    def __setitem__(self, key: str, value: Any) -> None:
        try:
            self.data[key] = json.loads(value)
        except json.decoder.JSONDecodeError:
            self.data[key] = value
        self._save_file()

    def _load_file(self) -> None:
        if not self.path.exists():
            self.path.touch()
        try:
            self.data = json.loads(self.path.read_text())
        except json.decoder.JSONDecodeError:
            self.data = {}

    def _save_file(self) -> None:
        self.path.write_text(json.dumps(self.data))

    def _delete_file(self) -> None:
        if self.path.exists() and self.path.is_file():
            self.path.unlink()


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    subparsers = parser.add_subparsers(
        title="subcommands", dest="subcommand", help="Subcommand help"
    )

    get_parser = subparsers.add_parser("get")
    get_parser.add_argument("key", help="Get value by [key]")

    set_parser = subparsers.add_parser("set")
    set_parser.add_argument("key")
    set_parser.add_argument("value", nargs="...")

    subparsers.add_parser("list")
    subparsers.add_parser("clear")

    return parser.parse_args()


def main() -> None:
    args = parse_args()
    if not args.subcommand:
        exit(SystemExit("Use list, get <value> or set <key> <value>..."))

    eky = Eky()
    if args.subcommand == "list":
        for key in eky.data.keys():
            print(key)
    elif args.subcommand == "clear":
        eky._delete_file()
    elif args.subcommand == "get":
        print(eky[args.key])
    elif args.subcommand == "set":
        eky[args.key] = " ".join(args.value)


if __name__ == "__main__":
    main()
