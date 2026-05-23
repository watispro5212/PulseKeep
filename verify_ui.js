const { chromium } = require('playwright');
const path = require('path');

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  // Serve files locally for verification
  const fs = require('fs');
  const http = require('http');
  const server = http.createServer((req, res) => {
    let filePath = path.join(__dirname, 'web', req.url === '/' ? 'index.html' : req.url);
    if (fs.existsSync(filePath)) {
      res.writeHead(200);
      res.end(fs.readFileSync(filePath));
    } else {
      res.writeHead(404);
      res.end();
    }
  });
  server.listen(3000);

  await page.goto('http://localhost:3000');
  await page.screenshot({ path: 'screenshot_index.png' });
  console.log('Index page screenshot taken.');

  await page.goto('http://localhost:3000/commands.html');
  await page.screenshot({ path: 'screenshot_commands.png' });
  console.log('Commands page screenshot taken.');

  await page.goto('http://localhost:3000/status.html');
  await page.screenshot({ path: 'screenshot_status.png' });
  console.log('Status page screenshot taken.');

  await page.goto('http://localhost:3000/team.html');
  await page.screenshot({ path: 'screenshot_team.png' });
  console.log('Team page screenshot taken.');

  await browser.close();
  server.close();
})();
