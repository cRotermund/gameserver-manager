import os
import typer

from libs.gsmclient.gsmclient import GSMClient
from dotenv import load_dotenv
from rich.console import Console

load_dotenv()
GSM_ACCESS_KEY = os.getenv("GSM_ACCESS_KEY")
GSM_ACCESS_SECRET = os.getenv("GSM_ACCESS_SECRET")
GSM_REGION = os.getenv("GSM_REGION")

client = GSMClient(GSM_ACCESS_KEY, GSM_ACCESS_SECRET, GSM_REGION)
console = Console()
app = typer.Typer()

def __waitForStatus(desired: str):
    with console.status("[bold blue]Waiting for server status...") as status:
        onpoll = lambda s: status.update(f"[bold blue]Waiting for server status... Status: {s}")
        try:
            client.wait_for_status(desired, onpoll)
        except:
            status = client.status()["status"]
            console.log(f"[bold red]Last server status observed: {status}")
            return False
    return True

@app.command()
def start():
    s = client.status()["status"]
    if s != "stopped":
        console.log("[bold red]Server is not stopped, can not start")
        return

    console.log("Sending the server start request...")
    client.start()

    if __waitForStatus("running"):
        console.log("[bold green]done")
    else:
        console.log("[bold red]Server did not transition states in timely manner, validate manually")

@app.command()
def stop():
    s = client.status()["status"]
    console.log("Sending the stop request...")
    client.stop()

    if __waitForStatus("stopped"):
        console.log("[bold green]done")
    else:
        console.log("[bold red]Server did not transition states in timely manner, validate manually")

@app.command()
def status():
    console.log("Getting server status...")
    s = client.status()["status"]
    console.log("Status: " + s)

if __name__ == "__main__":
    app()