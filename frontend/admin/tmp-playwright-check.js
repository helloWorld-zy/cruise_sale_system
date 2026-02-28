const { chromium } = require('@playwright/test');

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  page.on('request', (r) => {
    if (r.url().includes('cruises') || r.url().includes('undefined')) {
      console.log('REQ', r.method(), r.url());
    }
  });
  page.on('response', (r) => {
    if (r.url().includes('cruises') || r.url().includes('undefined')) {
      console.log('RES', r.status(), r.url());
    }
  });

  await page.goto('http://localhost:3013/cruises/create');
  await page.fill('input[placeholder="邮轮名称"]', 'PW-CRUISE-1');
  await page.fill('input[placeholder="公司 ID"]', '1');
  await page.fill('input[placeholder*="状态"]', 'draft');
  await page.click('button[type="submit"]');
  await page.waitForTimeout(1500);

  const err = await page.locator('.error').allTextContents();
  console.log('ERR_TEXT', JSON.stringify(err));
  console.log('URL_AFTER', page.url());

  await browser.close();
})();
