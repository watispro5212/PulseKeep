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

const ECONOMY_COMMANDS = new Set([
  'balance','daily','weekly','work','gamble','blackjack','slots','rob',
  'fish','mine','shop','buy','inventory','leaderboard','tip',
]);

const TICKET_COMMANDS = new Set(['ticketpanel']);

export class Bot {
  readonly client: Client;
  readonly cache: Cache;
  readonly db: any;
  readonly commands = new Collection<string, SlashCommand>();
  private rest: REST;
  private statusWebhookURL: string;
  private guildToggles = new Map<string, { economy?: boolean; tickets?: boolean }>();

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
      const cmdName = interaction.commandName;

      // Guild toggle enforcement
      if (interaction.guildId) {
        const toggles = await this.getGuildToggles(interaction.guildId);
        if (ECONOMY_COMMANDS.has(cmdName) && toggles.economy === false) {
          await interaction.reply({
            embeds: [new EmbedBuilder()
              .setTitle('Feature Disabled')
              .setDescription('The economy system is disabled in this server. An admin can enable it from the dashboard.')
              .setColor(Colors.Warning)],
            ephemeral: true,
          });
          return;
        }
        if (TICKET_COMMANDS.has(cmdName) && toggles.tickets === false) {
          await interaction.reply({
            embeds: [new EmbedBuilder()
              .setTitle('Feature Disabled')
              .setDescription('The ticket system is disabled in this server. An admin can enable it from the dashboard.')
              .setColor(Colors.Warning)],
            ephemeral: true,
          });
          return;
        }
      }

      const cmd = this.commands.get(cmdName);
      if (!cmd) {
        await interaction.reply({ content: 'Command not found. Use /help to see available commands.', ephemeral: true });
        return;
      }

      this.cache.incrementCommandsRun();

      // Log command to DB
      if (this.db) {
        try {
          const { commandLogs } = await import('../db/schema.js');
          await this.db.insert(commandLogs).values({
            guildId: interaction.guildId,
            userId: interaction.user.id,
            commandName: cmdName,
          });
        } catch {}
      }

      try {
        await cmd.execute({ bot: this, cache: this.cache, db: this.db }, interaction);
      } catch (err) {
        console.error(`Error executing /${cmdName}:`, err);
        const msg = 'An error occurred while executing this command. Please try again.';
        if (interaction.replied || interaction.deferred) {
          await interaction.followUp({ content: msg, ephemeral: true });
        } else {
          await interaction.reply({ content: msg, ephemeral: true });
        }
      }
    });
  }

  private async getGuildToggles(guildId: string): Promise<{ economy?: boolean; tickets?: boolean }> {
    const cached = this.guildToggles.get(guildId);
    if (cached) return cached;

    if (!this.db) return {};

    try {
      const { guildConfigs } = await import('../db/schema.js');
      const { eq } = await import('drizzle-orm');
      const rows: any[] = await this.db
        .select()
        .from(guildConfigs)
        .where(eq(guildConfigs.guildId, guildId))
        .limit(1);

      if (rows.length > 0 && rows[0]) {
        const cfg = rows[0];
        const toggles = {
          economy: cfg.economyEnabled !== false,
          tickets: cfg.ticketsEnabled !== false,
        };
        this.guildToggles.set(guildId, toggles);
        return toggles;
      }
    } catch {}

    return {};
  }

  private async handleTicketOpen(interaction: any) {
    await interaction.deferReply({ ephemeral: true });
    const guild = interaction.guild;
    if (!guild) {
      await interaction.editReply({ content: 'This can only be used in a server.' });
      return;
    }

    // Check tickets toggle
    const toggles = await this.getGuildToggles(guild.id);
    if (toggles.tickets === false) {
      await interaction.editReply({ content: 'The ticket system is disabled in this server.' });
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
        { id: guild.id, deny: [PermissionFlagsBits.ViewChannel] },
        { id: interaction.user.id, allow: [PermissionFlagsBits.ViewChannel, PermissionFlagsBits.SendMessages, PermissionFlagsBits.ReadMessageHistory] },
        { id: interaction.client.user!.id, allow: [PermissionFlagsBits.ViewChannel, PermissionFlagsBits.SendMessages, PermissionFlagsBits.ReadMessageHistory, PermissionFlagsBits.ManageChannels] },
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
    const cmds = this.commands.map((c) => c.data);
    if (cmds.length === 0) return;

    try {
      if (guildId) {
        await this.rest.put(Routes.applicationGuildCommands(this.client.user!.id, guildId), { body: cmds });
        console.log(`Registered ${cmds.length} commands for guild ${guildId}`);
      } else {
        await this.rest.put(Routes.applicationCommands(this.client.user!.id), { body: cmds });
        console.log(`Registered ${cmds.length} global commands`);
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
