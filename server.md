# PulseKeep Support Server Blueprint

This document contains a comprehensive plan for creating and configuring a fully functional Discord support server for the **PulseKeep** bot. All channels inside this server are organized into thematic categories and are designed to inherit their permissions directly from the parent category.

---

## 👥 1. Role Configurations & Hierarchy

Create the following roles in order (from top to bottom). The hierarchy is designed for secure server operation, clear styling, and effective permission delegation.

| Role Name | Recommended Hex Color | Core Purpose / Description | Essential Server Permissions |
| :--- | :---: | :--- | :--- |
| **Founder 👑** | `#E74C3C` (Red) | Server owner and project lead. | **Administrator** (Full permissions bypass) |
| **Administrator ⚙️** | `#E67E22` (Orange) | Senior staff, full operational control. | Manage Server, Manage Channels, Manage Roles, View Audit Log, Kick Members, Ban Members, Moderate Members, Manage Messages |
| **Moderator 🛡️** | `#2ECC71` (Green) | Chat and behavior moderation. | Kick Members, Ban Members, Moderate Members, Manage Messages, Mute Members, Deafen Members, Move Members, Timeout Members |
| **Support Team 🎫** | `#3498DB` (Blue) | Ticketing and user support handlers. | Manage Messages, Read Message History, Send Messages, Use Slash Commands, Attach Files, Embed Links |
| **PulseKeep Bot 🤖** | `#9B59B6` (Purple) | The active bot instance. | Manage Messages, Send Messages, Embed Links, Attach Files, Read Message History, Add Reactions, Use External Emojis |
| **Server Booster ✨** | `#FD79A8` (Pink) | Users actively boosting the server. | Change Nickname, Send Files, Use External Emojis, Add Reactions, Send Messages in Threads |
| **Verified Member ✅** | `#95A5A6` (Grey) | Standard members who passed verification. | View Channels, Send Messages, Read Message History, Add Reactions, Use External Emojis, Use Slash Commands |
| **@everyone** | Default | Base guest permissions (pre-verification). | View Channels (Specific lobbies only), Read Message History |

---

## 📁 2. Category-Level Permissions & Channels

To keep permissions simple and perfectly secure, all channels will inherit permissions from their parent category. Adjust category permissions according to the tables below, and do not specify individual channel overrides except where explicitly noted.

### Category 1: 📢 | INFORMATION & NEWS
**Purpose**: Read-only resources, rules, updates, and announcements.

#### Category Permissions Matrix
*   **Founder 👑 / Administrator ⚙️ / Moderator 🛡️**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Manage Messages: `ALLOW`
*   **PulseKeep Bot 🤖**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Embed Links: `ALLOW`
*   **@everyone / Verified Member ✅**:
    *   View Channel: `ALLOW`
    *   Send Messages: `DENY`
    *   Send Messages in Threads: `DENY`
    *   Add Reactions: `ALLOW`
    *   Read Message History: `ALLOW`

#### Channels (Inheriting All Permissions)
1.  `#rules-and-info` — Server guidelines, code of conduct, and bot invitation link.
2.  `#announcements` — Product updates, changelogs, and key releases.
3.  `#status-logs` — Automated bot status notifications and uptime reports.

---

### Category 2: 💬 | PULSEKEEP COMMUNITY
**Purpose**: General chat, discussions, community events, and user interaction.

#### Category Permissions Matrix
*   **Founder 👑 / Administrator ⚙️ / Moderator 🛡️**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Manage Messages: `ALLOW`
*   **Server Booster ✨**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Embed Links: `ALLOW`
    *   Attach Files: `ALLOW`
*   **Verified Member ✅**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Embed Links: `DENY` (Prevents link-spam from non-trusted members)
    *   Attach Files: `DENY` (Prevents image-spam from non-trusted members)
    *   Read Message History: `ALLOW`
    *   Add Reactions: `ALLOW`
*   **@everyone**:
    *   View Channel: `DENY` (Force users to pass gatekeeping verification)
*   **PulseKeep Bot 🤖**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Embed Links: `ALLOW`
    *   Attach Files: `ALLOW`

#### Channels (Inheriting All Permissions)
1.  `#welcome` — Entry point for new users (verification instructions).
2.  `#general-chat` — Core main text channel for general community interaction.
3.  `#bot-discussion` — Ask questions about bot usage and features.
4.  `#showcase` — Share configurations, servers, and panels powered by PulseKeep.

---

### Category 3: 🎫 | CLIENT SUPPORT
**Purpose**: Entrance lobby for ticketing, help commands, and general assistance.

#### Category Permissions Matrix
*   **Founder 👑 / Administrator ⚙️ / Moderator 🛡️ / Support Team 🎫**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Manage Messages: `ALLOW`
*   **Verified Member ✅**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Read Message History: `ALLOW`
    *   Use Slash Commands: `ALLOW`
*   **@everyone**:
    *   View Channel: `DENY`
*   **PulseKeep Bot 🤖**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Embed Links: `ALLOW`
    *   Attach Files: `ALLOW`

#### Channels (Inheriting All Permissions)
1.  `#support-faq` — Curated read-only channel with answers to common errors and guides.
2.  `#open-a-ticket` — Contains a sticky panel command (via PulseKeep Bot) enabling users to react/button-click to open a private ticket.
3.  `#pre-sales` — Inquiries regarding premium features or custom setups before buying.

---

### Category 4: 🔒 | ACTIVE SUPPORT TICKETS
**Purpose**: Dynamically generated private channels where users receive support.

#### Category Permissions Matrix
*   **Founder 👑 / Administrator ⚙️ / Support Team 🎫**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Read Message History: `ALLOW`
    *   Attach Files: `ALLOW`
    *   Embed Links: `ALLOW`
    *   Manage Messages: `ALLOW`
*   **PulseKeep Bot 🤖**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Embed Links: `ALLOW`
    *   Attach Files: `ALLOW`
*   **@everyone / Verified Member ✅**:
    *   View Channel: `DENY` (Strict isolation — ensure user private data is shielded)
    *   Send Messages: `DENY`

> [!NOTE]
> **Dynamic Override**: When a user clicks "Open Ticket" in `#open-a-ticket`, the bot generates a new channel in this category (e.g., `#ticket-1042`) and appends **only one override**: grant the *ticket creator* `View Channel: ALLOW` and `Send Messages: ALLOW` permissions. All other users remain locked out, inheriting the base category's hidden settings.

#### Channels (Generated Dynamically)
*   `#ticket-XXXX` — Private ticket chat channels.

---

### Category 5: 🤖 | BOT TESTING & SANDBOX
**Purpose**: Sandbox region for members to interact with the bot and try features safely.

#### Category Permissions Matrix
*   **Founder 👑 / Administrator ⚙️ / Moderator 🛡️**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Manage Messages: `ALLOW`
*   **Verified Member ✅**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Use Slash Commands: `ALLOW`
    *   Read Message History: `ALLOW`
    *   *Channel setting hint*: Turn on a 5-second slowmode on sandbox channels to protect performance.
*   **@everyone**:
    *   View Channel: `DENY`
*   **PulseKeep Bot 🤖**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Embed Links: `ALLOW`
    *   Attach Files: `ALLOW`
    *   Add Reactions: `ALLOW`

#### Channels (Inheriting All Permissions)
1.  `#sandbox-1` — Public spam channel for trying commands.
2.  `#sandbox-2` — Secondary test channel for advanced setup commands.

---

### Category 6: 🛡️ | STAFF QUARTERS (PRIVATE)
**Purpose**: Hidden internal coordination, staff chats, and bot administration.

#### Category Permissions Matrix
*   **Founder 👑 / Administrator ⚙️ / Moderator 🛡️ / Support Team 🎫**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Read Message History: `ALLOW`
    *   Attach Files: `ALLOW`
    *   Embed Links: `ALLOW`
    *   Manage Messages: `ALLOW`
*   **PulseKeep Bot 🤖**:
    *   View Channel: `ALLOW`
    *   Send Messages: `ALLOW`
    *   Embed Links: `ALLOW`
*   **@everyone / Verified Member ✅**:
    *   View Channel: `DENY` (Staff-only area is completely hidden)

#### Channels (Inheriting All Permissions)
1.  `#staff-chat` — Coordination channel for the support and moderation team.
2.  `#mod-logs` — Automated moderation logging (timeouts, kicks, bans).
3.  `#ticket-archives` — Dynamic transcripts of closed tickets posted by the bot.
