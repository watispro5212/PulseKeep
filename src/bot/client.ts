import {
  Client,
  GatewayIntentBits,
  Events,
  REST,
  Routes,
  Collection,
  EmbedBuilder,
  ChannelType,
  PermissionFlagsBits,
} from 'discord.js';
import type { Cache } from '../cache/index.js';
import type { SlashCommand } from './types.js';
import { Colors, footer, timestamp } from '../utils/embed.js';

export class Bot {
  readonly client: Client;
  readonly cache: Cache;
  readonly db: any;
  readonly commands = new Collection<string, SlashCommand>();
  private rest: REST;
  private statusWebhookURL: string;

  constructor(token: string, cache: Cache, db: any, statusWebhookURL: string) {
    this.client = new Client({
      intents: [
        GatewayIntentBits.Guilds,
        GatewayIntentBits.GuildMessages,
        GatewayIntentBits.MessageContent,
        GatewayIntentBits.GuildMembers,
        GatewayIntentBits.GuildVoiceStates,
      ],
    });
    this.cache = cache;
    this.db = db;
    this.statusWebhookURL = statusWebhookURL;
    this.rest = new REST({ version: '10' }).setToken(token);

    this.registerListeners();
  }

  private registerListeners() {
    this.client.once(Events.ClientReady, (c) => {
      console.log(`Logged in as ${c.user.tag}`);
      this.cache.setGuildsCount(this.client.guilds.cache.size);
    });

    this.client.on(Events.InteractionCreate, async (interaction) => {
      if (interaction.isButton() && interaction.customId === 'ticket_open') {
        await this.handleTicketOpen(interaction);
        return;
      }

      if (!interaction.isChatInputCommand()) return;

      const cmd = this.commands.get(interaction.commandName);
      if (!cmd) {
        await interaction.reply({ content: 'Command not found.', ephemeral: true });
        return;
      }

      this.cache.incrementCommandsRun();

      try {
        await cmd.execute({ bot: this, cache: this.cache, db: this.db }, interaction);
      } catch (err) {
        console.error(`Error executing ${interaction.commandName}:`, err);
        const msg = 'An error occurred while executing this command.';
        if (interaction.replied || interaction.deferred) {
          await interaction.followUp({ content: msg, ephemeral: true });
        } else {
          await interaction.reply({ content: msg, ephemeral: true });
        }
      }
    });
  }

  private async handleTicketOpen(interaction: any) {
    await interaction.deferReply({ ephemeral: true });
    const guild = interaction.guild;
    if (!guild) {
      await interaction.editReply({ content: 'This can only be used in a server.' });
      return;
    }

    let categoryId: string | undefined;
    if (this.db) {
      try {
        const { guildConfigs } = await import('../db/schema.js');
        const { eq } = await import('drizzle-orm');
        const rows: any[] = await this.db
          .select()
          .from(guildConfigs)
          .where(eq(guildConfigs.guildId, guild.id))
          .limit(1);
        if (rows[0]?.ticketCategoryId) categoryId = rows[0].ticketCategoryId;
      } catch {}
    }

    const existing = guild.channels.cache.find(
      (c: any) => c.name === `ticket-${interaction.user.id}`,
    );
    if (existing) {
      await interaction.editReply({ content: `You already have an open ticket: ${existing}.` });
      return;
    }

    const channel = await guild.channels.create({
      name: `ticket-${interaction.user.id}`,
      type: ChannelType.GuildText,
      parent: categoryId || undefined,
      permissionOverwrites: [
        {
          id: guild.id,
          deny: [PermissionFlagsBits.ViewChannel],
        },
        {
          id: interaction.user.id,
          allow: [PermissionFlagsBits.ViewChannel, PermissionFlagsBits.SendMessages, PermissionFlagsBits.ReadMessageHistory],
        },
        {
          id: interaction.client.user!.id,
          allow: [PermissionFlagsBits.ViewChannel, PermissionFlagsBits.SendMessages, PermissionFlagsBits.ReadMessageHistory, PermissionFlagsBits.ManageChannels],
        },
      ],
    });

    const emb = new EmbedBuilder()
      .setTitle('Ticket Created')
      .setDescription(`Welcome ${interaction.user}! A staff member will be with you shortly.`)
      .setColor(Colors.Tickets);

    await channel.send({ embeds: [footer(timestamp(emb))] });
    await interaction.editReply({ content: `Your ticket has been created: ${channel}` });
  }

  registerCommand(command: SlashCommand) {
    this.commands.set(command.data.name, command);
  }

  async registerSlashCommands(guildId?: string) {
    const commands = this.commands.map((c) => c.data);
    if (commands.length === 0) return;

    try {
      if (guildId) {
        await this.rest.put(Routes.applicationGuildCommands(this.client.user!.id, guildId), {
          body: commands,
        });
        console.log(`Registered ${commands.length} commands for guild ${guildId}`);
      } else {
        await this.rest.put(Routes.applicationCommands(this.client.user!.id), {
          body: commands,
        });
        console.log(`Registered ${commands.length} global commands`);
      }
    } catch (err) {
      console.error('Failed to register commands:', err);
    }
  }

  async start(token: string) {
    await this.client.login(token);
  }

  async stop() {
    this.client.destroy();
  }
}
