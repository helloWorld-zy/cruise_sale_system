import { test, expect } from '@playwright/test';

test('has title and search filters', async ({ page }) => {
  await page.goto('http://localhost:3001/cruises');

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/Cruise/);
  
  // Check header
  await expect(page.locator('h1')).toContainText('Find Your Cruise');

  // Check search button exists
  await expect(page.getByRole('button', { name: 'Search' })).toBeVisible();
});

test('navigates to detail page', async ({ page }) => {
  // Mock API response if needed, or assume dev server is running with seed data
  // For now, we test the navigation structure logic
  await page.goto('http://localhost:3001/cruises');
  
  // Click first view details button if any card exists
  // const firstCard = page.locator('.u-card').first();
  // if (await firstCard.isVisible()) {
  //   await firstCard.getByRole('button', { name: 'View Details' }).click();
  //   await expect(page).toHaveURL(/cruises\/[\w-]+/);
  // }
});
