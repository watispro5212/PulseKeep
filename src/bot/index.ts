import { Client, GatewayIntentBits, Events, EmbedBuilder } from 'discord.js';
import type { Interaction } from 'discord.js';
import { Cache } from '../cache/index.js';
import { AutomodEngine } from './automod/index.js';
import { EconomyStore } from './economy/index.js';

export class Bot {
    private client: Client;
    private cache: Cache;
    private db: any;
    private automod: AutomodEngine;
    private economy: EconomyStore;
    private statusWebhookURL: string;

    constructor(token: string, cache: Cache, db: any, statusWebhookURL: string) {
        this.client = new Client({
            intents: [
                GatewayIntentBits.Guilds,
                GatewayIntentBits.GuildMessages,
                GatewayIntentBits.MessageContent,
                GatewayIntentBits.GuildMembers,
            ],
        });
        this.cache = cache;
        this.db = db;
        this.statusWebhookURL = statusWebhookURL;
        this.automod = new AutomodEngine();
        this.economy = new EconomyStore(db);

        this.registerEvents();
    }

    private registerEvents() {
        this.client.once(Events.ClientReady, (c) => {
            console.log(`Ready! Logged in as ${c.user.tag}`);
            this.cache.setGuildsCount(this.client.guilds.cache.size);
        });

        this.client.on(Events.InteractionCreate, async (interaction: Interaction) => {
            if (!interaction.isChatInputCommand()) return;

            const { commandName } = interaction;
            this.cache.incrementCommandsRun();

            try {
                if (commandName === 'ping') {
                    const latency = Date.now() - interaction.createdTimestamp;
                    this.cache.addLatency(latency);
                    await interaction.reply(`Pong! Latency is ${latency}ms. API Latency is ${Math.round(this.client.ws.ping)}ms`);
                } else if (commandName === 'invite') {
                    const embed = new EmbedBuilder()
                        .setTitle('Invite PulseKeep')
                        .setDescription('Add PulseKeep to your server or join the support community.')
                        .addFields(
                            { name: 'Invite Bot', value: '[Click here](https://discord.com/oauth2/authorize?client_id=1507498795569512598&permissions=8&scope=bot%20applications.commands)', inline: false },
                            { name: 'Support Server', value: '[Click here](https://discord.gg/pulsekeep)', inline: false },
                            { name: 'Top.gg', value: '[Vote for us](https://top.gg/bot/1507498795569512598/vote)', inline: false }
                        )
                        .setColor(0x0099FF)
                        .setTimestamp();
                    await interaction.reply({ embeds: [embed], ephemeral: true });
                }
                // Handle other commands...
            } catch (error) {
                console.error(error);
                await interaction.reply({ content: 'There was an error while executing this command!', ephemeral: true });
            }
        });

        this.client.on(Events.MessageCreate, async (message) => {
            if (message.author.bot || !message.guild) return;

            const result = this.automod.checkMessage(message.guild.id, message.author.id, message.content);
            if (result.deleteMsg) {
                try {
                    await message.delete();
                    // Optional: Send a warning or log the action
                } catch (err) {
                    console.error('Failed to delete message:', err);
                }
            }
        });
    }

    public async start(token: string) {
        await this.client.login(token);
    }

    public async stop() {
        this.client.destroy();
    }

    public getAutomod(): AutomodEngine {
        return this.automod;
    }
}
