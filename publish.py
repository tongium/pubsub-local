import argparse
import json
import os
import sys
from typing import Any, Dict, Iterable, List, Optional

from dotenv import load_dotenv
from google.api_core.exceptions import Conflict
from google.cloud import pubsub_v1

load_dotenv()

project_id = os.getenv("PUBSUB_PROJECT_ID", "test-project")
publisher = pubsub_v1.PublisherClient()


def create_topic_if_not_exists(topic_id: str) -> str:
    name = publisher.topic_path(project_id, topic_id)

    try:
        publisher.create_topic(request={"name": name})
    except Conflict:
        pass

    return name


def read_jsonl_lines(path: Optional[str]) -> Iterable[str]:
    if path:
        with open(path, "r", encoding="utf-8") as file:
            for line in file:
                yield line
    else:
        for line in sys.stdin:
            yield line


def normalize_attributes(raw: Any) -> Dict[str, str]:
    if raw is None:
        return {}
    if not isinstance(raw, dict):
        raise ValueError("headers/attributes must be an object")
    return {str(k): str(v) for k, v in raw.items()}


def extract_payload(entry: Dict[str, Any]) -> Any:
    if "payload" in entry:
        return entry["payload"]
    if "data" in entry:
        return entry["data"]
    raise ValueError("missing payload/data")


def extract_topic(entry: Dict[str, Any]) -> str:
    for key in ("topic", "topic_id", "target_topic"):
        value = entry.get(key)
        if value:
            return str(value)
    raise ValueError("missing topic")


def extract_ordering_key(entry: Dict[str, Any]) -> Optional[str]:
    value = entry.get("ordering_key")
    return str(value) if value else None


def publish_entries(lines: Iterable[str]) -> int:
    futures: List[pubsub_v1.publisher.futures.Future] = []
    failures = 0

    for index, line in enumerate(lines, start=1):
        stripped = line.strip()
        if not stripped or stripped.startswith("#"):
            continue

        try:
            entry = json.loads(stripped)
            if not isinstance(entry, dict):
                raise ValueError("entry must be an object")

            topic_id = extract_topic(entry)
            topic_path = create_topic_if_not_exists(topic_id)

            payload = extract_payload(entry)
            data = json.dumps(payload, ensure_ascii=True).encode("utf-8")

            attributes = normalize_attributes(
                entry.get("headers") or entry.get("attributes")
            )
            ordering_key = extract_ordering_key(entry)

            future = publisher.publish(
                topic_path, data, ordering_key=ordering_key or "", **attributes
            )
            futures.append(future)
        except (ValueError, json.JSONDecodeError) as exc:
            failures += 1
            print(f"Line {index}: {exc}", file=sys.stderr)

    for future in futures:
        try:
            future.result()
        except Exception as exc:
            failures += 1
            print(f"Publish failed: {exc}", file=sys.stderr)

    return failures


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description="Publish messages from JSONL for batch testing"
    )
    parser.add_argument(
        "--jsonl",
        dest="jsonl_path",
        default=None,
        help="Path to JSONL file (defaults to stdin)",
    )
    return parser


def main() -> int:
    args = build_parser().parse_args()
    failures = publish_entries(read_jsonl_lines(args.jsonl_path))
    if failures:
        print(f"Completed with {failures} failures", file=sys.stderr)
        return 1
    print("Publish complete")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
