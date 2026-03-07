import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: false });
  const page = await browser.newPage();
  await page.goto('http://127.0.0.1:3013/login');
  console.log('Page URL:', page.url());
  console.log('Page title:', await page.title());
  const content = await page.content();
  console.log('Page has content:', content.length > 500);
  await page.waitForTimeout(30000);
  await browser.close();
})();
