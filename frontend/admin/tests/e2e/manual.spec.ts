import { test } from '@playwright/test';

test('Manual validation of Web, Admin, and Miniapp', async ({ page, context }) => {
    test.setTimeout(0); // Disable timeout so user can manually verify as long as they want

    // Admin module
    await page.goto('http://localhost:3000');

    // Web module
    const webPage = await context.newPage();
    await webPage.goto('http://localhost:3001');

    // Miniapp module
    const miniappPage = await context.newPage();
    await miniappPage.goto('http://localhost:5173');

    // Pause to allow manual inspection and interaction using Playwright Inspector
    await page.pause();
});
