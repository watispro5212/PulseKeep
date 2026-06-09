"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Bot = void 0;
const discord_js_1 = require("discord.js");
const cache_1 = require("../cache");
const automod_1 = require("./automod");
const economy_1 = require("./economy");
class Bot {
    client;
    cache;
    db;
    automod;
    economy;
    statusWebhookURL;
    constructor(token, cache, db, statusWebhookURL) {
        this.client = new discord_js_1.Client({
            intents: [
                discord_js_1.GatewayIntentBits.Guilds,
                discord_js_1.GatewayIntentBits.GuildMessages,
                discord_js_1.GatewayIntentBits.MessageContent,
                discord_js_1.GatewayIntentBits.GuildMembers,
            ],
        });
        this.cache = cache;
        this.db = db;
        this.statusWebhookURL = statusWebhookURL;
        this.automod = new automod_1.AutomodEngine();
        this.economy = new economy_1.EconomyStore(db);
        this.registerEvents();
    }
    registerEvents() {
        this.client.once(discord_js_1.Events.ClientReady, (c) => {
            console.log(`Ready! Logged in as ${c.user.tag}`);
            this.cache.setGuildsCount(this.client.guilds.cache.size);
        });
        this.client.on(discord_js_1.Events.InteractionCreate, async (interaction) => {
            if (!interaction.isChatInputCommand())
                return;
            const { commandName } = interaction;
            this.cache.incrementCommandsRun();
            try {
                if (commandName === 'ping') {
                    const latency = Date.now() - interaction.createdTimestamp;
                    this.cache.addLatency(latency);
                    await interaction.reply(`Pong! Latency is ${latency}ms. API Latency is ${Math.round(this.client.ws.ping)}ms`);
                }
                else if (commandName === 'invite') {
                    const embed = new discord_js_1.EmbedBuilder()
                        .setTitle('Invite PulseKeep')
                        .setDescription('Add PulseKeep to your server or join the support community.')
                        .addFields({ name: 'Invite Bot', value: '[Click here](https://discord.com/oauth2/authorize?client_id=1507498795569512598&permissions=8&scope=bot%20applications.commands)', inline: false }, { name: 'Support Server', value: '[Click here](https://discord.gg/pulsekeep)', inline: false }, { name: 'Top.gg', value: '[Vote for us](https://top.gg/bot/1507498795569512598/vote)', inline: false })
                        .setColor(0x0099FF)
                        .setTimestamp();
                    await interaction.reply({ embeds: [embed], ephemeral: true });
                }
                // Handle other commands...
            }
            catch (error) {
                console.error(error);
                await interaction.reply({ content: 'There was an error while executing this command!', ephemeral: true });
            }
        });
        this.client.on(discord_js_1.Events.MessageCreate, async (message) => {
            if (message.author.bot || !message.guild)
                return;
            const result = this.automod.checkMessage(message.guild.id, message.author.id, message.content);
            if (result.deleteMsg) {
                try {
                    await message.delete();
                    // Optional: Send a warning or log the action
                }
                catch (err) {
                    console.error('Failed to delete message:', err);
                }
            }
        });
    }
    async start(token) {
        await this.client.login(token);
    }
    async stop() {
        this.client.destroy();
    }
    getAutomod() {
        return this.automod;
    }
}
exports.Bot = Bot;
//# sourceMappingURL=index.js.map