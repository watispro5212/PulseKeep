# PulseKeep Fly.io Deployment Script
# Run from the project root
# Prerequisites: flyctl installed, logged in

Write-Host "=== PulseKeep Fly.io Deploy ===" -ForegroundColor Cyan
Write-Host ""

# 1. Set Fly.io secrets (required)
Write-Host "Step 1: Set Fly.io secrets..." -ForegroundColor Yellow
Write-Host "Make sure you have these secrets set on Fly.io:" -ForegroundColor White
Write-Host "  fly secrets set DISCORD_TOKEN=<new-token>" -ForegroundColor Gray
Write-Host "  fly secrets set DISCORD_CLIENT_SECRET=<your-secret>" -ForegroundColor Gray
Write-Host "  fly secrets set DATABASE_URL=<your-db-url>" -ForegroundColor Gray
Write-Host "  fly secrets set DISCORD_REDIRECT_URI=https://pulsekeep.fly.dev/auth/discord/callback" -ForegroundColor Gray
Write-Host ""

# 2. Deploy
Write-Host "Step 2: Deploy to Fly.io..." -ForegroundColor Yellow
Write-Host "  fly deploy" -ForegroundColor Gray
Write-Host ""

# 3. After deploy, update bot profile
Write-Host "Step 3 (optional): Update bot Discord profile" -ForegroundColor Yellow
Write-Host "  node scripts/update-bot-profile.mjs <token> [name] [description] [avatar]" -ForegroundColor Gray
Write-Host ""

# 4. Verify
Write-Host "Step 4: Verify deployment" -ForegroundColor Yellow
Write-Host "  fly logs    # watch live logs" -ForegroundColor Gray
Write-Host "  fly status  # check app status" -ForegroundColor Gray
