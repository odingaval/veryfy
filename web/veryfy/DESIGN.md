---
name: Veryfy
colors:
  surface: '#f8f9ff'
  surface-dim: '#cbdbf5'
  surface-bright: '#f8f9ff'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#eff4ff'
  surface-container: '#e5eeff'
  surface-container-high: '#dce9ff'
  surface-container-highest: '#d3e4fe'
  on-surface: '#0b1c30'
  on-surface-variant: '#44474d'
  inverse-surface: '#213145'
  inverse-on-surface: '#eaf1ff'
  outline: '#75777e'
  outline-variant: '#c5c6cd'
  surface-tint: '#515f78'
  primary: '#000000'
  on-primary: '#ffffff'
  primary-container: '#0d1c32'
  on-primary-container: '#76849f'
  inverse-primary: '#b9c7e4'
  secondary: '#006c49'
  on-secondary: '#ffffff'
  secondary-container: '#6cf8bb'
  on-secondary-container: '#00714d'
  tertiary: '#000000'
  on-tertiary: '#ffffff'
  tertiary-container: '#2a1700'
  on-tertiary-container: '#b87500'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#d6e3ff'
  primary-fixed-dim: '#b9c7e4'
  on-primary-fixed: '#0d1c32'
  on-primary-fixed-variant: '#39475f'
  secondary-fixed: '#6ffbbe'
  secondary-fixed-dim: '#4edea3'
  on-secondary-fixed: '#002113'
  on-secondary-fixed-variant: '#005236'
  tertiary-fixed: '#ffddb8'
  tertiary-fixed-dim: '#ffb95f'
  on-tertiary-fixed: '#2a1700'
  on-tertiary-fixed-variant: '#653e00'
  background: '#f8f9ff'
  on-background: '#0b1c30'
  surface-variant: '#d3e4fe'
typography:
  display-lg:
    fontFamily: Inter
    fontSize: 48px
    fontWeight: '700'
    lineHeight: 56px
    letterSpacing: -0.02em
  headline-lg:
    fontFamily: Inter
    fontSize: 32px
    fontWeight: '600'
    lineHeight: 40px
    letterSpacing: -0.01em
  headline-lg-mobile:
    fontFamily: Inter
    fontSize: 24px
    fontWeight: '600'
    lineHeight: 32px
    letterSpacing: -0.01em
  headline-md:
    fontFamily: Inter
    fontSize: 24px
    fontWeight: '600'
    lineHeight: 32px
  headline-sm:
    fontFamily: Inter
    fontSize: 20px
    fontWeight: '600'
    lineHeight: 28px
  body-lg:
    fontFamily: Inter
    fontSize: 18px
    fontWeight: '400'
    lineHeight: 28px
  body-md:
    fontFamily: Inter
    fontSize: 16px
    fontWeight: '400'
    lineHeight: 24px
  body-sm:
    fontFamily: Inter
    fontSize: 14px
    fontWeight: '400'
    lineHeight: 20px
  label-md:
    fontFamily: Inter
    fontSize: 12px
    fontWeight: '600'
    lineHeight: 16px
    letterSpacing: 0.05em
  code-md:
    fontFamily: Inter
    fontSize: 14px
    fontWeight: '500'
    lineHeight: 20px
rounded:
  sm: 0.125rem
  DEFAULT: 0.25rem
  md: 0.375rem
  lg: 0.5rem
  xl: 0.75rem
  full: 9999px
spacing:
  base: 8px
  xs: 4px
  sm: 12px
  md: 16px
  lg: 24px
  xl: 32px
  xxl: 48px
  container-max: 1280px
  gutter: 24px
---

## Brand & Style
The design system is engineered to project absolute authority, security, and institutional trust. It bridges the gap between the rigorous standards of government portals and the streamlined efficiency of modern fintech platforms. 

The aesthetic is **Corporate / Modern**, characterized by a disciplined use of whitespace, a high-contrast primary palette, and a focus on functional clarity. The emotional response should be one of "verified confidence"—users must feel that their data is handled with precision and that the information presented is immutable. Visual flourishes are minimized in favor of structural integrity and data density, ensuring that complex verification workflows remain intuitive and professional.

## Colors
The palette is anchored by a deep navy primary, providing a "heavy" visual weight that signals stability and heritage. 

- **Primary (#0A192F):** Used for navigation, primary actions, and key headers to ground the interface.
- **Success/Emerald (#10B981):** Reserved strictly for "Verified" or "Valid" states, providing a clear, high-contrast signal of truth.
- **Warning/Amber (#F59E0B):** Used for "Pending" or "Expiring Soon" notifications, designed to draw attention without inducing panic.
- **Neutrals:** A range of cool grays (Slate) is used for secondary text and UI borders to maintain a crisp, clean environment.
- **Backgrounds:** Use a very light off-white (#F8FAFC) to separate the main canvas from pure white surface cards, enhancing the sense of layered information.

## Typography
This design system utilizes **Inter** exclusively to ensure maximum legibility across high-density data views. The type scale is optimized for hierarchical clarity.

Headlines use a tighter letter-spacing and heavier weights to command attention, while body text maintains a generous line height to support long-form document reading and multi-step forms. Label styles are set in uppercase with slight letter spacing to differentiate metadata from user-generated content. For license numbers or hash strings, use a medium weight to mimic a "monospaced" feel while staying within the font family.

## Layout & Spacing
The layout employs a **12-column fixed grid** for desktop applications to ensure a consistent, structured dashboard experience. On smaller viewports, the grid transitions to a fluid model with 16px margins.

Spacing follows a strict 8px baseline grid to maintain rhythmic consistency. High-density dashboards should utilize 'sm' and 'md' spacing units to maximize visible data, while marketing or login screens should use 'xl' and 'xxl' units to promote focus and breathing room. Gutters are kept wide (24px) to ensure that even in data-heavy views, columns remain distinct and readable.

## Elevation & Depth
Depth is conveyed through a combination of **tonal layers** and **ambient shadows**. 

- **Level 0 (Background):** The base canvas uses the background neutral (#F8FAFC).
- **Level 1 (Cards/Surfaces):** Main content areas are white with a subtle 1px border (#E2E8F0) and no shadow, creating a "flat but distinct" appearance.
- **Level 2 (Interactive/Floating):** Modals, dropdowns, and active cards use a highly diffused, low-opacity shadow (0px 4px 12px rgba(10, 25, 47, 0.08)) to appear physically lifted without feeling "heavy."
- **Focus States:** Elements receive a 2px offset ring in the primary navy color to ensure clear keyboard navigation and accessibility.

## Shapes
The shape language is **Soft**, utilizing a 0.25rem (4px) base radius. This provides a modern polish that feels contemporary but maintains the structured, "straight-edged" serious nature of an institutional tool. Large containers like dashboard cards may scale to 0.5rem (8px), but buttons and input fields should strictly adhere to the base 4px radius to preserve a crisp, technical look.

## Components

### Buttons
- **Primary:** Deep navy background, white text. No gradients. High-contrast hover state (slight lighten).
- **Secondary:** Transparent background with a 1px navy border. 
- **Destructive:** Used for revoking licenses; deep crimson text and border to signal risk.

### Inputs & Forms
- Inputs use a white background with a 1px gray border. On focus, the border transitions to Primary Navy. 
- Labels are always positioned above the input field in the `label-md` style for maximum clarity during data entry.
- Error states must include both a red border and a trailing icon to support accessibility.

### Cards & Dashboards
- Verification cards should feature a thick 4px left-border accent colored by status (Emerald for valid, Amber for warning).
- Use subtle horizontal dividers between list items to maintain vertical rhythm in high-density tables.

### Status Badges (Chips)
- Statuses use a "Soft" pill shape with a low-opacity background of the status color (e.g., 10% Emerald) and high-contrast text of the same hue (100% Emerald). This ensures the status is the first thing a user's eye gravitates toward.

### Data Tables
- Use a clean, header-less or subtly-headered table style. Alternate row striping is discouraged; use subtle 1px borders instead to keep the UI light.