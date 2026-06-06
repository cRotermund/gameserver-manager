import { Client, Events, GatewayIntentBits } from "discord.js";

const token = process.env.DISCORD_TOKEN;
if (!token) {
	console.error("DISCORD_TOKEN environment variable is required");
	process.exit(1);
}

const client = new Client({
	intents: [
		GatewayIntentBits.Guilds, 
		GatewayIntentBits.GuildMessages,
		GatewayIntentBits.MessageContent
	],
});

client.once(Events.ClientReady, (ready) => {
	console.log(`control-bot logged in as ${ready.user.tag}`);
});

client.on(Events.MessageCreate, async (message) => {
	if (message.author.bot) return;

	if (message.content === "!ping") {
		await message.reply("pong");
	}
});

client.login(token);
