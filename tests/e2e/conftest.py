import logging
from collections.abc import MutableMapping
from typing import Any
from faker import Faker
import pytest

fake = Faker()

def pytest_tavern_beta_before_every_request(request_args: MutableMapping):
    message = f"Request: {request_args['method']} {request_args['url']}"

    params = request_args.get('params', None)
    if params:
        message += f"\nQuery parameters: {params}"

    message += f"\nRequest body: {request_args.get('json', '<no body>')}"

    logging.info(message)

def pytest_tavern_beta_after_every_response(expected: Any, response: Any) -> None:
    logging.info(f"Response: {response.status_code} {response.text}")

@pytest.fixture
def generate_random_email() -> str:
    return fake.email()
