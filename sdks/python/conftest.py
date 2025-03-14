import logging
import os
import subprocess
import time
from io import BytesIO
from threading import Thread
from typing import AsyncGenerator, Callable, Generator, cast

import psutil
import pytest
import pytest_asyncio
import requests

from hatchet_sdk import ClientConfig, Hatchet
from hatchet_sdk.config import ClientTLSConfig


@pytest.fixture(scope="session", autouse=True)
def token() -> str:
    result = subprocess.run(
        [
            "docker",
            "compose",
            "run",
            "--no-deps",
            "setup-config",
            "/hatchet/hatchet-admin",
            "token",
            "create",
            "--config",
            "/hatchet/config",
            "--tenant-id",
            "707d0855-80ab-4e1f-a156-f1c4546cbf52",
        ],
        capture_output=True,
        text=True,
    )

    token = result.stdout.strip()

    os.environ["HATCHET_CLIENT_TOKEN"] = token

    return token


@pytest_asyncio.fixture(scope="session")
async def aiohatchet(token: str) -> AsyncGenerator[Hatchet, None]:
    yield Hatchet(
        debug=True,
        config=ClientConfig(
            token=token,
            tls_config=ClientTLSConfig(strategy="none"),
        ),
    )


@pytest.fixture(scope="session")
def hatchet(token: str) -> Hatchet:
    return Hatchet(
        debug=True,
        config=ClientConfig(
            token=token,
            tls_config=ClientTLSConfig(strategy="none"),
        ),
    )


def wait_for_worker_health() -> bool:
    worker_healthcheck_attempts = 0
    max_healthcheck_attempts = 10

    while True:
        if worker_healthcheck_attempts > max_healthcheck_attempts:
            raise Exception(
                f"Worker failed to start within {max_healthcheck_attempts} seconds"
            )

        try:
            requests.get("http://localhost:8001/health", timeout=5)

            return True
        except Exception:
            time.sleep(1)

        worker_healthcheck_attempts += 1


@pytest.fixture()
def worker(
    request: pytest.FixtureRequest,
) -> Generator[subprocess.Popen[bytes], None, None]:
    example = cast(str, request.param)

    command = ["poetry", "run", example]

    logging.info(f"Starting background worker: {' '.join(command)}")

    proc = subprocess.Popen(
        command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, env=os.environ.copy()
    )

    # Check if the process is still running
    if proc.poll() is not None:
        raise Exception(f"Worker failed to start with return code {proc.returncode}")

    wait_for_worker_health()

    def log_output(pipe: BytesIO, log_func: Callable[[str], None]) -> None:
        for line in iter(pipe.readline, b""):
            print(line.decode().strip())

    Thread(target=log_output, args=(proc.stdout, logging.info), daemon=True).start()
    Thread(target=log_output, args=(proc.stderr, logging.error), daemon=True).start()

    yield proc

    logging.info("Cleaning up background worker")
    parent = psutil.Process(proc.pid)
    children = parent.children(recursive=True)
    for child in children:
        child.terminate()
    parent.terminate()

    _, alive = psutil.wait_procs([parent] + children, timeout=5)
    for p in alive:
        logging.warning(f"Force killing process {p.pid}")
        p.kill()
