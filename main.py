import os
import typer
import requests
import time
from dotenv import load_dotenv
from rich.console import Console
from requests_auth_aws_sigv4 import AWSSigV4

load_dotenv()
ENDPOINT_URL = "https://vd5vhqweprchm3fixxqb4emphe0zpolz.lambda-url.us-east-1.on.aws/"
SERVER_WAIT_TIMEOUT = 30
GSM_ACCESS_KEY = os.getenv("GSM_ACCESS_KEY")
GSM_ACCESS_SECRET = os.getenv("GSM_ACCESS_SECRET")
GSM_REGION = os.getenv("GSM_REGION")
GSM_SERVICE = os.getenv("GSM_SERVICE")

console = Console()
app = typer.Typer()

def __sign():
    return AWSSigV4(
        aws_access_key_id = GSM_ACCESS_KEY,
        aws_secret_access_key = GSM_ACCESS_SECRET,
        region = GSM_REGION,
        service = GSM_SERVICE)

def __send(body: object):
    sig = __sign()
    r = requests.request(
        "POST",
        ENDPOINT_URL,
        json = body,
        auth = sig
    )
    return r

def __status():
    body = { "action" : "status" }
    r = __send(body)
    return r.json()

def __waitForStatus(desired: str):
    with console.status("[bold blue]Waiting for server status...") as status:
        reached = False
        started_at = time.time()
        while not reached:
            s = __status()["status"]

            status.update("[bold blue]Waiting for server status... Status: " + s)
            reached = (s == desired)
            elapsed = (time.time() - started_at)

            if not reached and elapsed > SERVER_WAIT_TIMEOUT:
                console.log("[bold red]Timed out waiting for server. Last server status: " + s)
                return False

            if not reached:
                time.sleep(1)
    return True

@app.command()
def start():
    s = __status()["status"]
    if s != "stopped":
        console.log("[bold red]Server is not stopped, can not start")
        return

    console.log("Starting...")
    body = { "action" : "start" }
    r = __send(body)

    if __waitForStatus("running"):
        console.log("[bold green]done")
    else:
        console.log("[bold red]Server did not transition states in timely manner, validate manually")

@app.command()
def stop():
    s = __status()
    console.log("Stopping")
    body = { "action" : "stop" }
    r = __send(body)

    if __waitForStatus("stopped"):
        console.log("[bold green]done")
    else:
        console.log("[bold red]Server did not transition states in timely manner, validate manually")

@app.command()
def status():
    console.log("Status")
    s = __status()["status"]
    console.log(s)

if __name__ == "__main__":
    app()