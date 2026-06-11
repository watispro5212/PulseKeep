import { warnCommand } from './warn.js';
import { warningsCommand } from './warnings.js';
import { clearwarnsCommand } from './clearwarns.js';
import { moveCommand } from './move.js';
import { vckickCommand } from './vckick.js';
import { purgeCommand } from './purge.js';
import { kickCommand } from './kick.js';
import { banCommand } from './ban.js';
import { softbanCommand } from './softban.js';
import { muteCommand, unmuteCommand } from './mute.js';
import { nickCommand } from './nick.js';
import { roleCommand } from './role.js';
import { slowmodeCommand } from './slowmode.js';
import { lockCommand, unlockCommand } from './lockdown.js';
import { announceCommand } from './announce.js';
import { cleanCommand } from './clean.js';
import { historyCommand } from './history.js';

export const moderationCommands = [
  warnCommand,
  warningsCommand,
  clearwarnsCommand,
  moveCommand,
  vckickCommand,
  purgeCommand,
  kickCommand,
  banCommand,
  softbanCommand,
  muteCommand,
  unmuteCommand,
  nickCommand,
  roleCommand,
  slowmodeCommand,
  lockCommand,
  unlockCommand,
  announceCommand,
  cleanCommand,
  historyCommand,
];
