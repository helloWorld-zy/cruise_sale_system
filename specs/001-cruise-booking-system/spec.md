# Feature Specification: Cruise Booking System (MVP + Full Experience)

**Feature Branch**: `001-cruise-booking-system`
**Created**: 2026-02-10
**Status**: Draft
**Input**: User description: "MVP + One-Stop Full Experience (Section II)"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Guest Cruise Browsing & Discovery (Priority: P1)

As a potential customer, I want to browse and filter cruises by destination, date, and price, and view detailed information about the ship and cabins so that I can decide which cruise to book.

**Why this priority**: Core entry point for sales; without discovery, no bookings can occur.

**Independent Test**: Can be tested by accessing the public homepage, using filters, and verifying the details page loads with correct content (images, amenities, cabin list).

**Acceptance Scenarios**:

1. **Given** a visitor on the cruise list page, **When** they filter by "Destination: Japan" and "Date: Next Month", **Then** only relevant cruises are displayed.
2. **Given** a visitor on a specific cruise detail page, **When** they click "View Cabins", **Then** the list of available cabin types (Inside, Oceanview, Balcony) is shown with starting prices.
3. **Given** a visitor viewing cruise facilities, **When** they toggle categories (e.g., "Dining"), **Then** only dining facilities are shown with images and opening hours.

---

### User Story 2 - Booking & Checkout Flow (Priority: P1)

As a customer, I want to select a specific cabin, enter passenger details, and pay securely so that I can secure my vacation spot.

**Why this priority**: The primary revenue generation mechanism.

**Independent Test**: Can be tested by completing a full booking flow with a mock payment provider and verifying the order status changes to "Paid".

**Acceptance Scenarios**:

1. **Given** a selected cabin type, **When** the user chooses a specific available cabin number, **Then** the system locks the inventory for 15 minutes.
2. **Given** the passenger information form, **When** the user enters valid ID and contact details for all guests, **Then** the system proceeds to the payment confirmation screen.
3. **Given** a pending order, **When** the payment is successfully processed, **Then** the user receives a confirmation notification and the order status updates to "Confirmed".
4. **Given** a booking in progress, **When** the 15-minute lock expires without payment, **Then** the cabin is released back to inventory.

---

### User Story 3 - User Order Management (Priority: P1)

As a registered user, I want to view my booking history, check order status, and request cancellations if necessary so that I can manage my travel plans.

**Why this priority**: Essential for post-purchase experience and customer support reduction.

**Independent Test**: Can be tested by logging in as a user with existing orders and performing status checks or cancellation requests.

**Acceptance Scenarios**:

1. **Given** a logged-in user, **When** they access "My Orders", **Then** a list of all past and upcoming trips is displayed with current statuses.
2. **Given** a "Confirmed" order eligible for cancellation, **When** the user submits a cancellation request, **Then** the system calculates the refund amount based on the policy and submits it for admin audit.

---

### User Story 4 - Admin Content & Inventory Management (Priority: P1)

As an administrator, I want to manage cruise information, cabin inventory, and pricing rules so that the platform displays accurate products to users.

**Why this priority**: operational requirement to populate the system with sellable inventory.

**Independent Test**: Can be tested by logging into the admin panel, creating a new cruise/cabin entity, and verifying it appears on the public frontend.

**Acceptance Scenarios**:

1. **Given** an admin in the Cruise Management section, **When** they add a new ship with rich text description and images, **Then** the ship becomes visible in the backend list.
2. **Given** a specific sailing, **When** the admin sets a "Holiday" price modifier for a date range, **Then** the frontend prices for that period update accordingly.
3. **Given** a low-inventory warning threshold (e.g., <5 cabins), **When** sales reduce stock below this level, **Then** the system triggers a notification to operations staff.

---

### User Story 5 - Admin Order Processing & Finance (Priority: P2)

As a finance/ops staff member, I want to view all orders, audit refunds, and see daily financial reports so that I can ensure business health.

**Why this priority**: Critical for financial reconciliation and order fulfillment, though manual workarounds exist initially.

**Independent Test**: Can be tested by generating orders and verifying they appear in admin reports and dashboards.

**Acceptance Scenarios**:

1. **Given** a list of refund requests, **When** an admin approves a request, **Then** the refund is processed via the payment gateway and inventory is released.
2. **Given** the specific date, **When** the admin views the "Daily Reconciliation" report, **Then** the total transaction volume matches the payment provider's records.

---

### User Story 6 - Intelligent Discovery & Enhanced Booking (Priority: P2)

As a user, I want smart recommendations, price trend analysis, and flexible payment options (installments/multi-currency) to make informed and easier booking decisions.

**Why this priority**: Enhances conversion rates and user experience for the "Complete Experience" phase.

**Independent Test**: Verify recommendation API returns context-aware results; Verify price trend chart renders; Verify installment payment option availability.

**Acceptance Scenarios**:

1. **Given** a user with history of "Luxury" bookings, **When** they visit the homepage, **Then** "Suite" cabins are prioritized in recommendations.
2. **Given** a flight/cruise price history, **When** the user views a sailing, **Then** a "Price Trend" chart is displayed with a "Best Time to Buy" tip.
3. **Given** a booking total eligible for installments, **When** the user selects "Pay", **Then** a "Deposit + Balance" option is available.

---

### User Story 7 - Post-Booking & Onboard Experience (Priority: P2)

As a traveler, I want to book shore excursions, receive my Group Departure Notification, and receive real-time notifications about my trip to ensure a smooth journey.

**Why this priority**: Extends value beyond the initial sale.

**Independent Test**: Book an excursion add-on; Upload a departure notice as admin and verify user receipt; Simulate a WebSocket notification.

**Acceptance Scenarios**:

1. **Given** a confirmed booking, **When** the admin uploads a "Departure Notification" PDF, **Then** the user receives a pop-up alert and can download the file from their order details.
2. **Given** an upcoming port stop, **When** the user browses "Shore Excursions", **Then** available tours for that specific port/time are shown.
3. **Given** a change in boarding time, **When** the system updates the schedule, **Then** a real-time notification is pushed to the user's device.

---

### User Story 8 - Social & Community Engagement (Priority: P3)

As a user, I want to share my itinerary, read reviews, and potentially find group-buy deals to enhance the social aspect of travel.

**Why this priority**: Drives organic growth and retention (Ecosystem phase).

**Independent Test**: Submit a review; Share an itinerary poster; Join a group buy.

**Acceptance Scenarios**:

1. **Given** a completed trip, **When** the user submits a photo review, **Then** it enters a moderation queue before appearing on the cruise page.
2. **Given** a booked itinerary, **When** the user clicks "Share", **Then** a branded poster image is generated for social media sharing.

---

### User Story 9 - Intelligent Operations & Marketing (Priority: P2)

As an operations manager, I want dynamic pricing, CRM insights, and automated marketing tools to maximize revenue and efficiency.

**Why this priority**: Optimizes long-term profitability and operational scale.

**Independent Test**: Verify dynamic pricing rules trigger updates; Check CRM segment generation; Verify automated marketing email triggers.

**Acceptance Scenarios**:

1. **Given** inventory drops below 20%, **When** the Dynamic Pricing Engine runs, **Then** the price automatically increases by the configured percentage.
2. **Given** a user who browsed but didn't buy, **When** the "Cart Abandonment" rule triggers, **Then** a marketing message is sent with a reminder.

## Requirements *(mandatory)*

### Functional Requirements

**Core Product & Sales (MVP)**
- **FR-001**: System MUST allow administrators to create, update, and delete Cruise entities with rich media.
- **FR-002**: System MUST support hierarchical Cabin Type management (e.g., Balcony, Suite).
- **FR-003**: System MUST support Facility management classified by type.
- **FR-004**: System MUST allow defining Cabins (SKUs) associated with specific sailings.
- **FR-005**: System MUST maintain real-time inventory preventing double-booking.
- **FR-006**: System MUST support temporary "Inventory Lock" (15 mins).
- **FR-007**: System MUST support dynamic pricing matrices (Date/Occupancy/Pax Type).
- **FR-008**: System MUST trigger low-inventory alerts.

**Core Booking & Order (MVP)**
- **FR-009**: System MUST provide public search filters (Destination, Date, Port, Price).
- **FR-010**: System MUST collect detailed passenger information during booking.
- **FR-011**: System MUST track Order lifecycle (Pending -> Paid -> Confirmed -> etc.).
- **FR-012**: System MUST support automated cancellation policies.
- **FR-013**: System MUST support registration/login via Mobile/SMS and WeChat.
- **FR-014**: System MUST send automated notifications for key order events.
- **FR-015**: System MUST support "Frequent Passengers" management.

**Localization (Global)**
- **FR-016**: System MUST default to **Chinese (Simplified)** for all user and admin interfaces.
- **FR-017**: System MUST support manual language switching for other supported languages (e.g., English, Japanese).

**Intelligent Discovery & Decision (Phase 2+)**
- **FR-018**: System MUST provide AI-based cabin recommendations based on user history and preferences.
- **FR-019**: System MUST display visual Price Trend charts and "Best Time to Buy" indicators.
- **FR-020**: System MUST allow side-by-side comparison of up to 3 cabin types.

**Enhanced Booking & Payment (Phase 2+)**
- **FR-021**: System MUST support Installment payments (Deposit + Balance).
- **FR-022**: System MUST support Multi-currency display and settlement (CNY default, plus USD/HKD/JPY).
- **FR-023**: System MUST support OCR scanning for ID/Passport data entry.
- **FR-024**: System MUST support Bulk/Group booking workflows.

**Post-Booking & Interaction (Phase 2+)**
- **FR-025**: System MUST allow administrators to upload "Departure Notification" files (PDF) for orders and notify users.
- **FR-026**: System MUST provide specific "Countdown" and "Pre-trip Checklist" features.
- **FR-027**: System MUST support Real-time WebSocket updates for inventory/status and Price Alerts.
- **FR-028**: System MUST integrate AI Customer Service (FAQ) + Human handover.

**Social, Community & Ecosystem (Phase 2+)**
- **FR-029**: System MUST generate shareable "Itinerary Posters".
- **FR-030**: System MUST support "Refer-a-Friend" incentives.
- **FR-031**: System MUST support User Reviews (Score + Media) with moderation.
- **FR-032**: System MUST support UGC Travel Logs and Group Buying (Team Booking) features.

**Intelligent Ops (Phase 2+)**
- **FR-033**: System MUST support Dynamic Pricing Rules (Inventory-based, Competitor-based).
- **FR-034**: System MUST support Channel Inventory Management (OTA vs Direct).
- **FR-035**: System MUST provide comprehensive Revenue Management dashboards (RevPAR, etc.).
- **FR-036**: System MUST support CRM (Lifecycle Management) and Automated Marketing triggers.

### Key Entities *(include if feature involves data)*

- **Cruise**: Physical ship.
- **CabinType/Cabin**: Room definitions and physical instances.
- **Voyage**: Specific sailing.
- **Order/OrderItem**: Booking records.
- **Passenger**: Traveler details.
- **Review**: User feedback.
- **ShoreExcursion**: Add-on product.
- **MarketingRule**: Automation configuration.
- **PriceRule**: Dynamic pricing configuration.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can complete a booking flow in under 5 minutes.
- **SC-002**: System handles inventory locks with 0% double-booking.
- **SC-003**: Search results load in under 2 seconds.
- **SC-004**: 100% of orders reconciled automatically.
- **SC-005**: 100% Test Coverage (Strict Constitution Mandate).
- **SC-006**: AI Recommendations account for >10% of clicks on details page (Phase 2).
- **SC-007**: OCR recognition accuracy >95% for standard IDs (Phase 2).

### Edge Cases

- **Inventory Race Condition**: Simultaneous booking of last cabin.
- **Payment Timeout**: Lock expiry handling.
- **Currency Fluctuation**: Rate changes between booking and payment (Fixed at booking time?).
- **Dynamic Pricing Conflict**: Manual override vs Auto-rule.