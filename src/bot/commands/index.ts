import { economyCommands } from './economy/index.js';
import { moderationCommands } from './moderation/index.js';
import { utilityCommands } from './utility/index.js';
import { ticketCommands } from './tickets/index.js';
import { pingCommand } from './ping.js';
import { inviteCommand } from './invite.js';
import { helpCommand } from './help.js';
import { configureCommand } from './configure.js';
import { aboutCommand } from './about.js';
import { userinfoCommand } from './userinfo.js';
import { serverinfoCommand } from './serverinfo.js';
import { createServerCommand } from './createserver.js';

export const commands = [
  helpCommand,
  pingCommand,
  inviteCommand,
  aboutCommand,
  configureCommand,
  userinfoCommand,
  serverinfoCommand,
  createServerCommand,
  ...economyCommands,
  ...moderationCommands,
  ...utilityCommands,
  ...ticketCommands,
];
