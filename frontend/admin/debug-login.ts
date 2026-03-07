import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: false });
  const page = await browser.newPage();
  
  page.on('console', msg => {
    console.log('CONSOLE:', msg.type(), msg.text());
  });
  
  page.on('pageerror', err => {
    console.log('PAGE ERROR:', err.message);
  });
  
  await page.goto('http://127.0.0.1:3013/login');
  await page.waitForTimeout(3000);
  
  await page.fill('#username', 'admin');
  await page.fill('#password', 'admin123');
  await page.waitForTimeout(1000);
  
  await page.click('button[type="submit"]');
  await page.waitForTimeout(5000);
  
  console.log('Current URL:', page.url());
  const errorText = await page.locator('.login-form__error').textContent().catch(() => 'no error element');
  console.log('Error text:', errorText);
  
  await browser.close();
})();
