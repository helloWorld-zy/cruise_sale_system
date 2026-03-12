CREATE TABLE IF NOT EXISTS content_templates (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(120) NOT NULL,
  kind VARCHAR(40) NOT NULL,
  status SMALLINT NOT NULL DEFAULT 1,
  content_json TEXT NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_content_templates_kind_status ON content_templates(kind, status);

ALTER TABLE voyages
  ADD COLUMN IF NOT EXISTS fee_note_template_id BIGINT NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS fee_note_mode VARCHAR(20) NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS fee_note_content_json TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS booking_notice_template_id BIGINT NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS booking_notice_mode VARCHAR(20) NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS booking_notice_content_json TEXT NOT NULL DEFAULT '';