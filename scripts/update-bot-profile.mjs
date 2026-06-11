// Script to update the bot's Discord application profile
// Usage: node scripts/update-bot-profile.mjs <token> [name] [description] [avatar-path]
// Example: node scripts/update-bot-profile.mjs "MT..." "PulseKeep" "Your new description" ./avatar.png

const token = process.argv[2];
const name = process.argv[3];
const description = process.argv[4];
const avatarPath = process.argv[5];

if (!token) {
  console.log('Usage: node update-bot-profile.mjs <token> [name] [description] [avatar-path]');
  process.exit(1);
}

async function main() {
  // Get current application info
  const appRes = await fetch('https://discord.com/api/v10/oauth2/applications/@me', {
    headers: { Authorization: `Bot ${token}` },
  });
  if (!appRes.ok) {
    console.error('Failed to fetch application:', await appRes.text());
    process.exit(1);
  }
  const app = await appRes.json();
  console.log('Current bot name:', app.name);
  console.log('Current bot description:', app.description);

  const body = {};
  if (name && name !== app.name) body.name = name;
  if (description) body.description = description;

  if (avatarPath) {
    const fs = await import('fs');
    const mime = avatarPath.endsWith('.png') ? 'image/png' : 'image/jpeg';
    const base64 = fs.readFileSync(avatarPath).toString('base64');
    body.icon = `data:${mime};base64,${base64}`;
  }

  if (Object.keys(body).length === 0) {
    console.log('No changes to make.');
    return;
  }

  const updateRes = await fetch('https://discord.com/api/v10/applications/@me', {
    method: 'PATCH',
    headers: {
      Authorization: `Bot ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(body),
  });

  if (!updateRes.ok) {
    console.error('Failed to update bot profile:', await updateRes.text());
    process.exit(1);
  }

  const updated = await updateRes.json();
  console.log('\nBot profile updated!');
  console.log('Name:', updated.name);
  console.log('Description:', updated.description);
}

main().catch(console.error);
