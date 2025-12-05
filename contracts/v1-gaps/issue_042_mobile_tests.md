# Issue #042: Add Mobile Responsiveness Tests for Architect View

**Priority:** 🟢 LOW  
**Estimated Tokens:** ~600 (Low complexity)  
**Agent Role:** Test

---

## 1. 🎫 Related Issue Context

**Gap Reference:** GAP-008 from v1-analysis.md

Per Project Charter PR Template: "Confirmed responsiveness on a mobile viewport (360px width check)."

The `layout.spec.ts` tests mobile sidebar, but the Architect view (primary input area) is not tested on mobile.

---

## 2. 📋 Acceptance Criteria

### E2E Tests
- [ ] Test Architect view at 360x640 viewport (mobile portrait)
- [ ] Test Architect view at 414x896 viewport (iPhone 11 Pro Max)
- [ ] Verify textarea is visible and fills available width
- [ ] Verify TokenMeter is visible below textarea
- [ ] Verify typing works correctly on mobile viewport
- [ ] Verify paste works correctly (large text)
- [ ] Verify color change on token limit still works

### Visual Verification
- [ ] No horizontal scroll on mobile
- [ ] Text is readable (not too small)
- [ ] Touch targets are adequate size (min 44px)

---

## 3. 📊 Token Efficiency Strategy

- Add to existing architect.spec.ts file
- Reuse existing test patterns
- ~30 new lines of test code

---

## 4. 🏗️ Technical Specification

### Mobile Viewport Tests
```typescript
// frontend/tests/e2e/architect.spec.ts

test.describe('Architect View - Mobile', () => {
    test.beforeEach(async ({ page }) => {
        await page.setViewportSize({ width: 360, height: 640 });
        await page.goto('/');
    });

    test('should display architect input on mobile', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        await expect(input).toBeVisible();
        
        // Verify it fills the width (accounting for padding)
        const box = await input.boundingBox();
        expect(box?.width).toBeGreaterThan(300); // ~360 - padding
    });

    test('should display token meter on mobile', async ({ page }) => {
        const meter = page.getByTestId('token-meter');
        await expect(meter).toBeVisible();
        
        // Should be below the textarea, not side-by-side
        const input = page.getByTestId('architect-input');
        const inputBox = await input.boundingBox();
        const meterBox = await meter.boundingBox();
        
        expect(meterBox?.y).toBeGreaterThan(inputBox!.y + inputBox!.height - 10);
    });

    test('should update meter when typing on mobile', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        const meterBar = page.locator('[data-testid="token-meter"] .bg-gray-800 > div');
        
        // Simulate mobile typing
        await input.tap();
        await input.fill('a'.repeat(1600)); // 400 tokens = 5%
        
        await expect(meterBar).toHaveAttribute('style', /width: 5%/);
    });

    test('should not have horizontal scroll', async ({ page }) => {
        // Check that body width equals viewport width
        const bodyWidth = await page.evaluate(() => document.body.scrollWidth);
        expect(bodyWidth).toBeLessThanOrEqual(360);
    });
});

test.describe('Architect View - Large Mobile', () => {
    test.beforeEach(async ({ page }) => {
        await page.setViewportSize({ width: 414, height: 896 });
        await page.goto('/');
    });

    test('should work on iPhone 11 Pro Max size', async ({ page }) => {
        const input = page.getByTestId('architect-input');
        await expect(input).toBeVisible();
        
        // Fill with large text and verify meter changes color
        await input.fill('a'.repeat(30000));
        
        const meterBar = page.locator('[data-testid="token-meter"] .bg-gray-800 > div');
        await expect(meterBar).toHaveClass(/bg-red-500/);
    });
});
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| MODIFY | `frontend/tests/e2e/architect.spec.ts` |

---

## 6. ✅ Definition of Done

1. Mobile viewport tests pass at 360x640
2. Mobile viewport tests pass at 414x896
3. No horizontal scrollbar on mobile
4. TokenMeter visible and functional on mobile
5. Typing and paste work correctly on mobile viewport
6. Screenshots saved for visual verification (optional)
