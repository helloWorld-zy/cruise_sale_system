# 2026-03-05 Admin Voyage Unification Design

## 1. Background and Goal

Current admin has both `routes` and `voyages` management. The target is to fully merge them into a single voyage-centric model and UI.

Goals:
- Remove route management completely.
- Keep voyage CRUD as the only admin entry.
- Add rich voyage itinerary editing in admin create/edit.
- Persist itinerary as structured data in database.

Confirmed constraints:
- Route module is fully removed (not hidden, not compatibility mode).
- Data migration strategy is clean rebuild (no historical route/voyage migration).
- Accommodation data uses `bool + text`.
- One day can contain multiple stop cities.
- Add a manual top-level field `brief_info` (short voyage info), not derived from itinerary.

## 2. High-Level Architecture

### 2.1 Domain Model Direction
- `voyages` is the root aggregate for scheduling and operations.
- `voyage_itineraries` stores day/stop-level planning details.
- `routes` table and `/admin/routes` API are removed.

### 2.2 Existing Core Relation Compatibility
- Keep dependencies that already use `voyage_id`:
  - `cabin_skus.voyage_id`
  - `bookings.voyage_id`
- No change is required in cabin/booking foreign key direction.

## 3. Database Design

### 3.1 `voyages` table changes
Keep fields:
- `id`
- `cruise_id`
- `code`
- `depart_date`
- `return_date`
- `status`
- `created_at`
- `updated_at`

Add fields:
- `brief_info VARCHAR(300) NOT NULL DEFAULT ''`

Remove fields:
- `route_id`

Constraints:
- `code` unique and not null.
- `cruise_id` references `cruises(id)`.
- `return_date >= depart_date`.
- `brief_info` trimmed non-empty in application validation.

### 3.2 New `voyage_itineraries` table
Columns:
- `id BIGSERIAL PRIMARY KEY`
- `voyage_id BIGINT NOT NULL REFERENCES voyages(id) ON DELETE CASCADE`
- `day_no INT NOT NULL`
- `stop_index INT NOT NULL`
- `city VARCHAR(120) NOT NULL`
- `summary TEXT NOT NULL DEFAULT ''`
- `eta_time TIME NULL`
- `etd_time TIME NULL`
- `has_breakfast BOOLEAN NOT NULL DEFAULT FALSE`
- `has_lunch BOOLEAN NOT NULL DEFAULT FALSE`
- `has_dinner BOOLEAN NOT NULL DEFAULT FALSE`
- `has_accommodation BOOLEAN NOT NULL DEFAULT FALSE`
- `accommodation_text VARCHAR(300) NOT NULL DEFAULT ''`
- `created_at TIMESTAMP NOT NULL DEFAULT NOW()`
- `updated_at TIMESTAMP NOT NULL DEFAULT NOW()`

Indexes and constraints:
- `UNIQUE(voyage_id, day_no, stop_index)`
- `CHECK(day_no >= 1)`
- `CHECK(stop_index >= 1)`
- Index on `(voyage_id, day_no, stop_index)` for ordered fetch.

### 3.3 Route removal scope
- Drop `routes` table.
- Remove route-related indexes/constraints linked to old `voyages.route_id`.
- Remove route-related migration tests and add new migration tests for voyage itinerary schema.

## 4. Backend API and Validation Design

## 4.1 Admin endpoints
Keep and refactor:
- `GET /api/v1/admin/voyages`
- `GET /api/v1/admin/voyages/:id`
- `POST /api/v1/admin/voyages`
- `PUT /api/v1/admin/voyages/:id`
- `DELETE /api/v1/admin/voyages/:id`

Remove:
- `/api/v1/admin/routes/*`

## 4.2 DTOs
### Create/Update request
- `cruise_id: number`
- `code: string`
- `brief_info: string`
- `depart_date: string(datetime)`
- `return_date: string(datetime)`
- `itineraries: ItineraryInput[]`

### `ItineraryInput`
- `day_no: number`
- `stop_index: number`
- `city: string`
- `summary: string`
- `eta_time?: string(time)`
- `etd_time?: string(time)`
- `has_breakfast: boolean`
- `has_lunch: boolean`
- `has_dinner: boolean`
- `has_accommodation: boolean`
- `accommodation_text: string`

### List response
- Voyage base fields + `itinerary_days` + `first_stop_city`.

### Detail response
- Voyage base fields + ordered `itineraries` by `(day_no, stop_index)`.

## 4.3 Validation Rules
Voyage level:
- `code` required and unique.
- `brief_info` required, trimmed non-empty, max length 300.
- `cruise_id` must exist.
- `return_date >= depart_date`.

Itinerary level:
- At least one itinerary entry.
- `(day_no, stop_index)` must be unique.
- `day_no >= 1`, `stop_index >= 1`.
- For each `day_no`, `stop_index` should be continuous from `1..N`.
- `city` required.
- Optional: if both `eta_time` and `etd_time` exist, `eta_time <= etd_time`.
- If `has_accommodation == false`, normalize `accommodation_text = ''`.

## 4.4 Transaction boundaries
Create:
1. Begin transaction.
2. Insert voyage.
3. Batch insert itinerary rows.
4. Commit.

Update:
1. Begin transaction.
2. Update voyage base fields.
3. Delete old itinerary rows by `voyage_id`.
4. Batch insert new itinerary rows.
5. Commit.

Delete:
- Return dependency conflict when cabin/booking references block deletion.
- Otherwise delete voyage (itinerary rows cascade delete).

## 4.5 Error code additions
Add business-level codes for better frontend UX:
- `ErrVoyageCodeDuplicate`
- `ErrVoyageDateInvalid`
- `ErrVoyageItineraryInvalid`
- `ErrVoyageHasDependents`

## 5. Frontend Admin Design

## 5.1 Navigation and page scope
- Remove route menu entry.
- Keep only voyage management menu.
- Remove route pages under `frontend/admin/app/pages/routes/*`.
- Refactor voyage pages under `frontend/admin/app/pages/voyages/*`.

## 5.2 Voyage list page
Columns:
- `ID`
- `航次代码`
- `简短信息`
- `所属邮轮`
- `出发日期`
- `结束日期`
- `行程天数`
- `首日停靠`
- `操作`

Keep existing confirm dialog pattern for delete.

## 5.3 Voyage create/edit top fields
- Cruise selector.
- Voyage code.
- Brief info (manual input).
- Depart date.
- Return date.

No auto-generation of brief info from itinerary.

## 5.4 Itinerary editor interaction
Editor model supports day and stop nesting:
- Day groups: `第X天`.
- Each day has one or more stop cards.

Per-stop fields:
- City
- Summary
- ETA time
- ETD time
- Breakfast/Lunch/Dinner checkboxes
- Accommodation checkbox
- Accommodation text (only shown when checked)

Actions:
- Add day
- Remove day
- Add stop within day
- Remove stop

Normalization before submit:
- Re-number days to continuous `1..N`.
- Re-number stop index in each day to `1..N`.
- Flatten to backend DTO format.

## 5.5 Component split
To avoid mega page complexity:
- `VoyageItineraryEditor.vue`
- `VoyageDayCard.vue`
- `VoyageStopCard.vue`

Create and edit pages share the same editor and payload mapping helpers.

## 6. Testing Strategy

## 6.1 Backend
- Migration tests for route removal and new itinerary schema.
- Repository/service/handler tests for voyage CRUD + itinerary validation.
- Conflict tests for delete with dependent cabin/booking records.

## 6.2 Frontend (Vitest)
- `voyages-new.spec.ts`:
  - day/stop add/remove and reindex behavior.
  - payload includes `brief_info` and flattened itineraries.
- `voyages-id.spec.ts`:
  - detail load, edit save, delete flow.
- `voyages-index.spec.ts`:
  - list rendering includes brief info.
  - delete error mapping.

## 6.3 Regression focus
- Cabin and booking chain still works with voyage-centric data.
- Admin voyage create -> cabin association -> booking path remains valid.

## 7. Rollout and Operations

Given clean rebuild strategy:
1. Stop services.
2. Rebuild database from migration baseline.
3. Apply new migrations.
4. Start backend and frontend.
5. Run target test suite and smoke checks.

Operational note:
- Historical route/voyage data is intentionally discarded.
- Backup is still recommended before reset.

## 8. Risks and Mitigations

- Risk: Large cross-layer change may break admin forms.
  - Mitigation: Component split + focused unit tests + payload contract tests.
- Risk: Delete behavior regressions due to new FK flow.
  - Mitigation: Explicit conflict code mapping and dedicated tests.
- Risk: DTO mismatch between frontend nested UI and backend flat itinerary rows.
  - Mitigation: Single normalize helper and contract tests on mapped payload.

## 9. Non-goals
- No compatibility mode for routes.
- No migration of historical route/voyage data.
- No automatic brief info generation from itinerary.
