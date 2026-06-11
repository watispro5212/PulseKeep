import { economyCommands } from './economy/index.js';
import { moderationCommands } from './moderation/index.js';
import { utilityCommands } from './utility/index.js';
import { ticketCommands } from './tickets/index.js';
import { pingCommand } from './ping.js';
import { inviteCommand } from './invite.js';
import { helpCommand } from './help.js';
import { configureCommand } from './configure.js';

export const commands = [
  helpCommand,
  pingCommand,
  inviteCommand,
  configureCommand,
  ...economyCommands,
  ...moderationCommands,
  ...utilityCommands,
  ...ticketCommands,
];
