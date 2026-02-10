# Research & Technical Decisions: Cruise Booking System

**Feature**: `001-cruise-booking-system`
**Date**: 2026-02-10

## 1. Map Provider Integration

*   **Decision**: **Gaode Map (AMap)**
*   **Rationale**: The system mandates a Chinese-default experience. Gaode provides superior data accuracy, loading speed, and API stability within China compared to Google Maps or Mapbox. It also offers excellent support for WeChat Mini-programs.
*   **Alternatives Considered**:
    *   *Mapbox*: Good for global style customization but slower in China; styling overkill for MVP.
    *   *Google Maps*: Blocked in China, unsuitable for domestic users.

## 2. OCR Service (Phase 2)

*   **Decision**: **Aliyun OCR (Passport/ID Card)**
*   **Rationale**: Industry standard in China with high accuracy for ID Cards and Passports. Provides a unified SDK for Go.
*   **Alternatives Considered**:
    *   *Tencent Cloud OCR*: Viable alternative, can be a fallback.
    *   *Self-hosted Tesseract*: Too much maintenance overhead and lower accuracy for Chinese IDs.

## 3. Inventory Locking Mechanism

*   **Decision**: **Redis `SETNX` (Key: `lock:cabin:{id}`) with TTL + Lua Scripts**
*   **Rationale**: 
    *   Simple and high performance for the "15-minute lock" requirement.
    *   Lua scripts ensure atomicity when extending locks or releasing them upon payment.
    *   Avoids the complexity of Redlock unless a distributed Redis cluster is explicitly required (MVP uses standard Redis).
*   **Alternatives Considered**:
    *   *Database Row Locking (`SELECT FOR UPDATE`)*: Too heavy on the primary DB for high-concurrency browsing/booking attempts.
    *   *In-memory (Go Sync)*: Doesn't work across multiple backend replicas (Kubernetes).

## 4. Order State Machine

*   **Decision**: **Custom Go Struct-based Implementation**
*   **Rationale**:
    *   Allows strict type checking of states.
    *   Easier to integrate with GORM `BeforeSave`/`AfterSave` hooks for audit logging.
    *   External libraries (like `looplab/fsm`) often add unnecessary reflection or "stringly-typed" logic.
*   **Alternatives Considered**:
    *   *looplab/fsm*: Popular but adds dependency overhead for relatively linear logic.

## 5. Rich Text Content Storage

*   **Decision**: **TipTap (Frontend) + HTML/JSON (Backend) + MinIO (Images)**
*   **Rationale**:
    *   TipTap produces clean JSON/HTML output.
    *   Images are stripped from base64 and uploaded to MinIO, storing only URLs in the content.
    *   Ensures lightweight database records and CDN-ready assets.

## 6. Real-time Notifications

*   **Decision**: **NATS JetStream + Gorilla WebSocket**
*   **Rationale**:
    *   NATS JetStream handles the event decoupling (e.g., "Order Paid" event published by Payment Service).
    *   WebSocket service subscribes to NATS subjects and pushes to connected clients.
    *   Scales horizontally (unlike keeping state in a single monolithic app).
