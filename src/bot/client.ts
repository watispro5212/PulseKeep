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
  GuildMember,
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
} from 'discord.js';
import type { Cache } from '../cache/index.js';
import type { Config } from '../config.js';
import type { SlashCommand } from './types.js';
import { Colors, footer, timestamp } from '../utils/embed.js';
import { eq } from 'drizzle-orm';
import { startSpamCleanup } from './automod/index.js';
import { guildConfigs, commandLogs } from '../db/schema.js';

const ECONOMY_COMMANDS = new Set([
  'balance','daily','weekly','work','gamble','blackjack','slots','rob','pay',
  'fish','mine','shop','buy','inventory','use','leaderboard','tip','vote','search',
]);

const TICKET_COMMANDS = new Set(['ticketpanel', 'ticket']);

export class Bot {
  readonly client: Client;
  readonly cache: Cache;
  readonly db: any;
  readonly config: Config;
  readonly commands = new Collection<string, SlashCommand>();
  private rest: REST;
  private statusWebhookURL: string;
  private guildToggles = new Map<string, Record<string, any>>();
  private guildTogglesTtl = new Map<string, number>();
  // tracks mod actions per user — stops someone from getting spam-banned
  private modSpam = new Map<string, { count: number; windowStart: number }>();
  // userId:cmdName -> when cooldown expires
  private commandCooldowns = new Map<string, number>();
  // tiny event bus so the API server can react to guild config changes
  private listeners = new Map<string, Set<(...args: any[]) => void>>();

  on(event: 'guildConfigUpdated', listener: (guildId: string) => void): this;
  on(event: string, listener: (...args: any[]) => void): this;
  on(event: string, listener: (...args: any[]) => void): this {
    if (!this.listeners.has(event)) this.listeners.set(event, new Set());
    this.listeners.get(event)!.add(listener);
    return this;
  }

  emit(event: string, ...args: any[]): void {
    const set = this.listeners.get(event);
    if (!set) return;
    for (const fn of set) {
      try { fn(...args); } catch (err) { console.error(`[Event:${event}]`, err); }
    }
  }

  // returns 0 if not on cooldown, seconds remaining otherwise
  // economy = 3s, mod = 2s, everything else = 1s
  checkCommandCooldown(userId: string, commandName: string): number {
    const now = Date.now();
    const key = `${userId}:${commandName}`;
    const expire = this.commandCooldowns.get(key);
    if (expire && now < expire) return (expire - now) / 1000;
    const MOD = new Set(['warn','mute','kick','ban','softban','purge','clean','slowmode','lock','unlock','nick','role','move','vckick','announce']);
    const seconds = ECONOMY_COMMANDS.has(commandName) ? this.config.cooldownEconomy : MOD.has(commandName) ? this.config.cooldownModeration : 1;
    this.commandCooldowns.set(key, now + seconds * 1000);
    return 0;
  }

  /**
   * Returns true if a moderator has triggered too many mod actions in a short
   * window. Used as a soft brake to avoid ban-spam accidents.
   */
  registerModAction(guildId: string, userId: string): { throttled: boolean; resetIn: number } {
    const key = `${guildId}:${userId}`;
    const now = Date.now();
    const WINDOW = 10_000;
    const LIMIT = 6;
    const cur = this.modSpam.get(key);
    if (!cur || now - cur.windowStart > WINDOW) {
      this.modSpam.set(key, { count: 1, windowStart: now });
      return { throttled: false, resetIn: 0 };
    }
    cur.count += 1;
    if (cur.count > LIMIT) {
      return { throttled: true, resetIn: WINDOW - (now - cur.windowStart) };
    }
    return { throttled: false, resetIn: 0 };
  }

  async getGuildConfig(guildId: string): Promise<any> {
    if (!this.db) return {};
    try {
      const rows: any[] = await this.db
        .select()
        .from(guildConfigs)
        .where(eq(guildConfigs.guildId, guildId))
        .limit(1);
      return rows[0] || {};
    } catch {
      return {};
    }
  }

  invalidateGuildToggles(guildId: string) {
    this.guildToggles.delete(guildId);
    this.guildTogglesTtl.delete(guildId);
    this.emit('guildConfigUpdated', guildId);
  }

  async logToChannel(guildId: string, embed: EmbedBuilder) {
    const toggles = await this.getGuildToggles(guildId);
    if (toggles.modlogsEnabled === false) return;
    if (!toggles.logChannelId) return;
    try {
      const guild = this.client.guilds.cache.get(guildId)
        || await this.client.guilds.fetch(guildId).catch(() => null);
      if (!guild) return;
      const channel = guild.channels.cache.get(toggles.logChannelId)
        || await guild.channels.fetch(toggles.logChannelId).catch(() => null);
      if (channel?.isTextBased()) {
        await channel.send({ embeds: [embed] });
      }
    } catch (err) {
      console.error(`[logToChannel] Failed in guild ${guildId}:`, err);
    }
  }

  constructor(token: string, cache: Cache, db: any, statusWebhookURL: string, config: Config) {
    this.client = new Client({
      shards: 'auto',
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
    this.config = config;
    this.statusWebhookURL = statusWebhookURL;
    this.rest = new REST({ version: '10' }).setToken(token);

    this.registerListeners();
    startSpamCleanup();
  }

  private registerListeners() {
    this.client.once(Events.ClientReady, (c) => {
      console.log(`Logged in as ${c.user.username}`);
      this.cache.setGuildsCount(this.client.guilds.cache.size);
      this.cache.setBotGuilds(
        this.client.guilds.cache.map((g: any) => ({ id: g.id, name: g.name }))
      );
      this.cache.setTotalUserCount(
        this.client.guilds.cache.reduce((sum: number, g: any) => sum + (g.memberCount || 0), 0)
      );
      this.sendStatusWebhook();
    });

    // log disconnects so we can tell when the bot flaps
    this.client.on(Events.ShardDisconnect, (close, shardId) => {
      console.warn(`[Bot] Shard ${shardId} disconnected: ${close?.code} ${close?.reason ?? ''}`);
    });
    this.client.on(Events.ShardReconnecting, (shardId) => {
      console.warn(`[Bot] Shard ${shardId} reconnecting…`);
    });
    this.client.on(Events.ShardReady, (shardId) => {
      console.log(`[Bot] Shard ${shardId} ready.`);
    });
    this.client.on(Events.ShardError, (err, shardId) => {
      console.error(`[Bot] Shard ${shardId} error:`, err);
    });
    this.client.on(Events.Warn, (msg) => console.warn('[Discord]', msg));
    this.client.on(Events.Error, (err) => console.error('[Discord]', err));

    this.client.on(Events.GuildCreate, (guild) => {
      this.cache.setGuildsCount(this.client.guilds.cache.size);
      this.cache.setBotGuilds(
        this.client.guilds.cache.map((g: any) => ({ id: g.id, name: g.name }))
      );
      this.cache.setTotalUserCount(
        this.client.guilds.cache.reduce((sum: number, g: any) => sum + (g.memberCount || 0), 0)
      );
      console.log(`[Guild] Joined: ${guild.name} (${guild.id})`);
      this.sendStatusWebhook();
    });

    this.client.on(Events.GuildDelete, (guild) => {
      this.cache.setGuildsCount(this.client.guilds.cache.size);
      this.cache.setBotGuilds(
        this.client.guilds.cache.map((g: any) => ({ id: g.id, name: g.name }))
      );
      this.cache.setTotalUserCount(
        this.client.guilds.cache.reduce((sum: number, g: any) => sum + (g.memberCount || 0), 0)
      );
      console.log(`[Guild] Left/removed: ${guild.id}`);
      this.sendStatusWebhook();
    });

    // welcome messages
    this.client.on(Events.GuildMemberAdd, async (member: GuildMember) => {
      if (member.user.bot) return;
      this.cache.setTotalUserCount(this.cache.getTotalUserCount() + 1);
      const toggles = await this.getGuildToggles(member.guild.id);
      if (!toggles.welcomeEnabled || !toggles.welcomeChannelId) return;
      try {
        const channel = member.guild.channels.cache.get(toggles.welcomeChannelId)
          || await member.guild.channels.fetch(toggles.welcomeChannelId).catch(() => null);
        if (!channel?.isTextBased()) {
          console.warn(`[Welcome] Channel ${toggles.welcomeChannelId} not found or not text-based in guild ${member.guild.id}`);
          return;
        }
        const memberCount = member.guild.memberCount;
        const emb = new EmbedBuilder()
          .setTitle(`👋 Welcome to ${member.guild.name}!`)
          .setDescription(
            `Hey ${member.user} — welcome aboard!\n\n` +
            `You're member **#${memberCount.toLocaleString()}** of **${member.guild.name}**. ` +
            `Take a look around, say hi, and have fun.`
          )
          .setThumbnail(member.user.displayAvatarURL({ size: 256 }))
          .setColor(Colors.Success)
          .setFooter({ text: `Account created` })
          .setTimestamp();
        await channel.send({ embeds: [emb] });
        // also log it to the mod channel
        const logEmb = new EmbedBuilder()
          .setTitle('Member Joined')
          .setDescription(`${member.user.tag} (${member.user.id}) — member #${memberCount.toLocaleString()}`)
          .setThumbnail(member.user.displayAvatarURL({ size: 128 }))
          .setColor(Colors.Utility)
          .setTimestamp();
        this.logToChannel(member.guild.id, logEmb);
      } catch (err) {
        console.error(`[Welcome] Failed to send welcome message in guild ${member.guild.id}:`, err);
      }
    });

    this.client.on(Events.GuildMemberRemove, () => {
      this.cache.setTotalUserCount(Math.max(0, this.cache.getTotalUserCount() - 1));
    });

    // voice state logging
    this.client.on(Events.VoiceStateUpdate, (oldState, newState) => {
      const guildId = oldState.guild?.id || newState.guild?.id;
      if (!guildId) return;
      const userId = oldState.id || newState.id;

      const oldChannel = oldState.channelId;
      const newChannel = newState.channelId;

      if (!oldChannel && newChannel) {
        const logEmb = new EmbedBuilder()
          .setDescription(`<@${userId}> joined voice channel <#${newChannel}>`)
          .setColor(Colors.Utility)
          .setTimestamp();
        this.logToChannel(guildId, logEmb);
      } else if (oldChannel && !newChannel) {
        const logEmb = new EmbedBuilder()
          .setDescription(`<@${userId}> left voice channel <#${oldChannel}>`)
          .setColor(Colors.Utility)
          .setTimestamp();
        this.logToChannel(guildId, logEmb);
      } else if (oldChannel && newChannel && oldChannel !== newChannel) {
        const logEmb = new EmbedBuilder()
          .setDescription(`<@${userId}> moved from <#${oldChannel}> to <#${newChannel}>`)
          .setColor(Colors.Utility)
          .setTimestamp();
        this.logToChannel(guildId, logEmb);
      }
    });

    // auto-mod
    this.client.on(Events.MessageCreate, async (message) => {
      if (message.author.bot) return;
      if (!message.guildId || !message.guild) return;
      const member = message.member;
      if (!member) return;
      const isMod = member.permissions.has('ManageMessages') || member.permissions.has('Administrator');
      if (isMod) return;
      const { runAutomod } = await import('./automod/index.js');
      const { action, reason } = await runAutomod(
        this, message.guildId, message.channelId,
        message.author.id, message.content, member.roles.cache.map(r => r.id),
      );
      if (!action) return;
      switch (action) {
        case 'delete':
          await message.delete().catch(() => {});
          break;
        case 'warn': {
          const warnEmb = new EmbedBuilder()
            .setTitle('⚠️ Auto-Mod Warning')
            .setDescription(`**Reason:** ${reason}`)
            .setColor(Colors.Warning)
            .setTimestamp();
          await message.channel.send({ content: `<@${message.author.id}>`, embeds: [warnEmb] }).then(m => {
            setTimeout(() => m.delete().catch(() => {}), 5000);
          });
          break;
        }
        case 'timeout': {
          await message.delete().catch(() => {});
          await member.timeout(10 * 60 * 1000, reason ?? undefined).catch(() => {});
          const timeoutEmb = new EmbedBuilder()
            .setTitle('🔇 Auto-Mod Timeout')
            .setDescription(`<@${message.author.id}> has been timed out for **10 minutes**.\n**Reason:** ${reason}`)
            .setColor(Colors.Moderation)
            .setTimestamp();
          await message.channel.send({ embeds: [timeoutEmb] }).then(m => {
            setTimeout(() => m.delete().catch(() => {}), 8000);
          });
          break;
        }
      }
    });

    this.client.on(Events.InteractionCreate, async (interaction) => {
      if (interaction.isButton() && interaction.customId === 'ticket_open') {
        await this.handleTicketOpen(interaction);
        return;
      }

      if (interaction.isButton() && interaction.customId === 'ticket_close_button') {
        await this.handleTicketCloseButton(interaction);
        return;
      }

      if (!interaction.isChatInputCommand()) return;
      const cmdName = interaction.commandName;

      // check guild toggles
      if (interaction.guildId) {
        const toggles = await this.getGuildToggles(interaction.guildId);
        if (ECONOMY_COMMANDS.has(cmdName) && toggles.economy === false) {
          await interaction.reply({
            embeds: [new EmbedBuilder()
              .setTitle('Feature Disabled')
              .setDescription('The economy system is disabled in this server. An admin can enable it from the dashboard.')
              .setColor(Colors.Warning)],
            flags: 64,
          });
          return;
        }
        if (TICKET_COMMANDS.has(cmdName) && toggles.tickets === false) {
          await interaction.reply({
            embeds: [new EmbedBuilder()
              .setTitle('Feature Disabled')
              .setDescription('The ticket system is disabled in this server. An admin can enable it from the dashboard.')
              .setColor(Colors.Warning)],
            flags: 64,
          });
          return;
        }
      }

      // check cooldown
      if (interaction.user.id !== this.config.botOwnerID && interaction.user.id !== this.config.botCoOwnerID) {
        const remaining = this.checkCommandCooldown(interaction.user.id, cmdName);
        if (remaining > 0) {
          await interaction.reply({
            embeds: [new EmbedBuilder()
              .setTitle('⏳ Slow Down')
              .setDescription(`Please wait **${remaining.toFixed(1)}s** before using this command again.`)
              .setColor(Colors.Warning)],
            flags: 64,
          });
          return;
        }
      }

      const cmd = this.commands.get(cmdName);
      if (!cmd) {
        await interaction.reply({ content: 'Command not found. Use /help to see available commands.', flags: 64 });
        return;
      }

      this.cache.incrementCommandsRun();

      // log to db
      if (this.db) {
        try {
          await this.db.insert(commandLogs).values({
            guildId: interaction.guildId,
            userId: interaction.user.id,
            commandName: cmdName,
          });
        } catch (logErr) {
          console.error(`[CommandLog] Failed to log command ${cmdName}:`, logErr);
        }
      }

      try {
        await cmd.execute({ bot: this, cache: this.cache, db: this.db, config: this.config }, interaction);
      } catch (err) {
        console.error(`Error executing /${cmdName}:`, err);
        const msg = 'An error occurred while executing this command. Please try again.';
        if (interaction.replied || interaction.deferred) {
          await interaction.followUp({ content: msg, flags: 64 });
        } else {
          await interaction.reply({ content: msg, flags: 64 });
        }
      }
    });
  }

  async getGuildToggles(guildId: string): Promise<Record<string, any>> {
    const cached = this.guildToggles.get(guildId);
    if (cached) {
      const expiresAt = this.guildTogglesTtl.get(guildId) || 0;
      // cache toggles for 60s — stops hammering the db on every message
      if (Date.now() < expiresAt) return cached;
      this.guildToggles.delete(guildId);
      this.guildTogglesTtl.delete(guildId);
    }

    if (!this.db) return {};

    try {
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
          modlogsEnabled: cfg.modlogsEnabled !== false,
          welcomeEnabled: cfg.welcomeEnabled === true,
          automodEnabled: cfg.automodEnabled !== false,
          logChannelId: cfg.logChannelId || undefined,
          welcomeChannelId: cfg.welcomeChannelId || undefined,
          voteChannelId: cfg.voteChannelId || undefined,
          ticketCategoryId: cfg.ticketCategoryId || undefined,
        };
        this.guildToggles.set(guildId, toggles);
        this.guildTogglesTtl.set(guildId, Date.now() + 60_000);
        return toggles;
      }
    } catch (err) {
      console.error(`[getGuildToggles] DB error for ${guildId}:`, err);
    }

    return {};
  }

  private async handleTicketOpen(interaction: any) {
    await interaction.deferReply({ flags: 64 });
    const guild = interaction.guild;
    if (!guild) {
      await interaction.editReply({ content: 'This can only be used in a server.' });
      return;
    }

    // check ticket toggle
    const toggles = await this.getGuildToggles(guild.id);
    if (toggles.tickets === false) {
      await interaction.editReply({ content: 'The ticket system is disabled in this server.' });
      return;
    }

    let categoryId: string | undefined;
    if (this.db) {
      try {
        const rows: any[] = await this.db
          .select()
          .from(guildConfigs)
          .where(eq(guildConfigs.guildId, guild.id))
          .limit(1);
        if (rows[0]?.ticketCategoryId) categoryId = rows[0].ticketCategoryId;
      } catch (ticketErr) {
        console.error('[Ticket] Failed to fetch ticket category:', ticketErr);
      }
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
      .setTitle('🎫 Ticket Created')
      .setDescription(
        `Welcome ${interaction.user}! A staff member will be with you shortly.\n\n` +
        `**What to do next**\n` +
        `• Describe your issue as clearly as you can\n` +
        `• Add screenshots or links if they help\n` +
        `• Use \`/ticket add <user>\` to bring someone else in\n` +
        `• Use \`/ticket close\` or the button below when you're done`
      )
      .setColor(Colors.Tickets);

    const row = new ActionRowBuilder<ButtonBuilder>().addComponents(
      new ButtonBuilder()
        .setCustomId('ticket_close_button')
        .setLabel('Close Ticket')
        .setEmoji('🔒')
        .setStyle(ButtonStyle.Danger),
    );

    await channel.send({ embeds: [footer(timestamp(emb))], components: [row] });
    await interaction.editReply({ content: `✅ Your ticket has been created: ${channel}` });
  }

  private async handleTicketCloseButton(interaction: any) {
    const channel = interaction.channel;
    if (!channel || !channel.name?.startsWith?.('ticket-')) {
      await interaction.reply({ content: 'This button only works inside a ticket channel.', flags: 64 });
      return;
    }

    // Only the ticket creator or staff (ManageChannels) can close.
    const isStaff = interaction.member?.permissions?.has?.('ManageChannels');
    const openerMatch = channel.name === `ticket-${interaction.user.id}`;
    if (!isStaff && !openerMatch) {
      await interaction.reply({
        content: 'Only the ticket creator or staff can close this ticket. Use `/ticket close` as staff if needed.',
        flags: 64,
      });
      return;
    }

    await interaction.deferReply({ flags: 64 });

    // best-effort transcript to mod channel
    try {
      if ('messages' in channel && interaction.guildId) {
        const messages = await (channel as any).messages.fetch({ limit: 100 }).catch(() => null);
        if (messages && messages.size > 0) {
          const transcript = messages
            .sort((a: any, b: any) => a.createdTimestamp - b.createdTimestamp)
            .map((m: any) => `[${new Date(m.createdTimestamp).toISOString()}] ${m.author?.tag ?? 'unknown'}: ${m.content ?? ''}`)
            .join('\n')
            .slice(0, 6000);
          const transcriptEmb = new EmbedBuilder()
            .setTitle(`Ticket transcript — #${channel.name}`)
            .setDescription('```\n' + transcript + '\n```')
            .setColor(Colors.Tickets)
            .setFooter({ text: `${messages.size} message(s) • closed by ${interaction.user.tag}` });
          await this.logToChannel(interaction.guildId, transcriptEmb);
        }
      }
    } catch (err) {
      console.error('[ticket close button] transcript failed:', err);
    }

    const emb = new EmbedBuilder()
      .setTitle('🔒 Ticket Closed')
      .setDescription(`Ticket closed by ${interaction.user}. Deleting channel in **10 seconds**…`)
      .setColor(Colors.Tickets);
    await interaction.editReply({ embeds: [timestamp(emb)] });
    setTimeout(() => channel.delete().catch(() => {}), 10_000);
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

  private async sendStatusWebhook() {
    if (!this.statusWebhookURL) return;
    try {
      const uptime = Math.floor((Date.now() - this.cache.getStartedAt().getTime()) / 1000);
      const d = Math.floor(uptime / 86400);
      const h = Math.floor((uptime % 86400) / 3600);
      const m = Math.floor((uptime % 3600) / 60);
      const uptimeStr = `${d > 0 ? d + 'd ' : ''}${h > 0 || d > 0 ? h + 'h ' : ''}${m}m`;

      const emb = {
        embeds: [{
          title: 'PulseKeep Status Update',
          color: 0x7c5cfc,
          fields: [
            { name: 'Servers', value: String(this.cache.getGuildsCount()), inline: true },
            { name: 'Users', value: String(this.cache.getTotalUserCount()), inline: true },
            { name: 'Commands Run', value: String(this.cache.getCommandsRun()), inline: true },
            { name: 'Avg Latency', value: `${Math.round(this.cache.getAvgLatency())}ms`, inline: true },
            { name: 'Uptime', value: uptimeStr, inline: true },
          ],
          timestamp: new Date().toISOString(),
        }],
      };

      await fetch(this.statusWebhookURL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(emb),
      });
    } catch (err) {
      console.error('[StatusWebhook] Failed to send:', err);
    }
  }

  async start(token: string) {
    await this.client.login(token);
    setInterval(() => this.sendStatusWebhook(), 300000);
  }

  async stop() {
    this.client.destroy();
  }
}
