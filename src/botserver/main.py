# bot.py
import os
import discord

from gsmclient import GSMClient
from discord.ext import commands
from dotenv import load_dotenv

load_dotenv()

GSM_ACCESS_KEY = os.getenv("GSM_ACCESS_KEY")
GSM_ACCESS_SECRET = os.getenv("GSM_ACCESS_SECRET")
GSM_REGION = os.getenv("GSM_REGION")
BOT_TOKEN = os.getenv('DISCORD_BOT_TOKEN')
CLIENT_ID = os.getenv('DISCORD_CLIENT_ID')
CLIENT_SECRET = os.getenv('DISCORD_CLIENT_SECRET')

intents = discord.Intents.default()
intents.messages = True
intents.message_content = True
bot = commands.Bot(command_prefix='!', intents=intents)
client = GSMClient(GSM_ACCESS_KEY, GSM_ACCESS_SECRET, GSM_REGION)

async def __wait_for_status(desired: str, friendly_verb: str, ctx: commands.Context):
    try:
        client.wait_for_status(desired, None)
        await ctx.send(f"The server has been {friendly_verb}.")
    except:
        await ctx.send(f"Could not verify the server status, verify manually.")

@bot.command()
async def serverstart(ctx):
    print(f'Executing server start command, requested by {ctx.author}')

    s = client.status()["status"]
    if s != "stopped":
        await ctx.send("Server is not stopped, can not start")
        return
    else:
        await ctx.send(f"I'm working on that...")
        client.start()
        await __wait_for_status("running", "started", ctx)

@bot.command()
async def serverstop(ctx):
    print(f'Executing server stop command, requested by {ctx.author}')

    s = client.status()["status"]
    if s == "stopped":
        await ctx.send("Server is already stopped")
        return
    else:
        await ctx.send(f"I'm working on that...")
        client.stop()
        await __wait_for_status("stopped", "stopped", ctx)

@bot.command()
async def serverstatus(ctx):
    print(f'Executing server status command, requested by {ctx.author}')
    s = client.status()["status"]
    await ctx.send(f"The server's status is: {s}")

bot.run(BOT_TOKEN)