ALTER TABLE voyages
  DROP COLUMN IF EXISTS booking_notice_content_json,
  DROP COLUMN IF EXISTS booking_notice_mode,
  DROP COLUMN IF EXISTS booking_notice_template_id,
  DROP COLUMN IF EXISTS fee_note_content_json,
  DROP COLUMN IF EXISTS fee_note_mode,
  DROP COLUMN IF EXISTS fee_note_template_id;

DROP TABLE IF EXISTS content_templates;