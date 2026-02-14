import asyncio
import json
import os
from pathlib import Path
from typing import Any, Callable, List, Optional

import yaml
from dotenv import load_dotenv
from google.api_core.exceptions import Conflict
from google.cloud import pubsub_v1
from pydantic import BaseModel

load_dotenv()


project_id = os.getenv("PUBSUB_PROJECT_ID", "test-project")
publisher = pubsub_v1.PublisherClient()
subscriber = pubsub_v1.SubscriberClient()


class Subscriber(BaseModel):
    id: str
    enabled_message_ordering: bool = False


class Topic(BaseModel):
    id: str
    subscribers: List[Subscriber]


class Record(BaseModel):
    data: Any
    attributes: dict
    message_id: str
    ordering_key: Optional[str] = None


class Settings(BaseModel):
    topics: List[Topic]


def create_callback(id: str) -> Callable[[pubsub_v1.subscriber.message.Message], None]:
    def callback(message: pubsub_v1.subscriber.message.Message) -> None:
        print(f"Received message on topic {id}: {message.message_id}")
        save_message_to_file(id, message)
        message.ack()

    return callback


def save_message_to_file(
    topic_id: str, message: pubsub_v1.subscriber.message.Message
) -> None:
    # Create messages directory if not exists
    messages_dir = Path("./messages") / topic_id
    messages_dir.mkdir(parents=True, exist_ok=True)

    # Parse message data
    try:
        data = json.loads(message.data.decode("utf-8"))
    except json.JSONDecodeError, UnicodeDecodeError:
        data = message.data.decode("utf-8")

    attributes = {}
    for k, v in message.attributes.items():
        attributes[k] = v

    message_obj = {
        "publish_time": message.publish_time.isoformat(),
        "message_id": message.message_id,
        "attributes": attributes,
        "ordering_key": message.ordering_key,
        "data": data,
    }

    filename = messages_dir / f"{message.message_id}.json"

    # Write to file
    with open(filename, "w") as f:
        json.dump(message_obj, f, indent=2)


def create_topic_if_not_exists(id: str) -> str:
    name = publisher.topic_path(project_id, id)

    try:
        topic = publisher.create_topic(request={"name": name})
        print(f"Topic created: {topic}")
    except Conflict:
        pass

    return name


def create_subscription_if_not_exists(id: str, topic: str) -> str:
    name = subscriber.subscription_path(project_id, id)

    try:
        subscription = subscriber.create_subscription(
            request={"name": name, "topic": topic}
        )

        print(f"Subscription created: {subscription}")
    except Conflict:
        pass

    return name


def setup() -> list:
    settings: Settings
    # Load configuration from YAML file
    with open("settings.yaml", "r") as file:
        settings = Settings(**yaml.safe_load(file))

    futures = []

    for topic in settings.topics:
        topic_name = create_topic_if_not_exists(topic.id)
        subscriber_name = create_subscription_if_not_exists(
            f"{topic.id}-echo", topic_name
        )

        future = subscriber.subscribe(
            subscriber_name, callback=create_callback(topic.id)
        )

        futures.append(asyncio.to_thread(future.result))

        for sub in topic.subscribers:
            create_subscription_if_not_exists(sub.id, topic_name)

    return futures


async def main():
    futures = setup()
    results = await asyncio.gather(*futures)
    print(f"Results: {results}")


if __name__ == "__main__":
    asyncio.run(main())
