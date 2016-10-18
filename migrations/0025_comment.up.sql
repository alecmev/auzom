ALTER TABLE comment RENAME COLUMN user_id TO created_by;
ALTER TABLE comment RENAME CONSTRAINT comment_user_id_fkey TO comment_created_by_fkey;

ALTER TABLE comment DROP COLUMN match_id;
ALTER TABLE comment DROP COLUMN news_id;
ALTER TABLE comment ADD COLUMN target text NOT NULL;
ALTER TABLE comment ADD COLUMN target_id int4 NOT NULL;
ALTER TABLE comment ADD CONSTRAINT comment_target_check CHECK (target IN ('match'));
CREATE INDEX comment_target_target_id_idx ON public.comment (target, target_id);

ALTER TABLE comment RENAME COLUMN last_edited_at TO updated_at;
ALTER TABLE comment ADD COLUMN updated_by integer;
ALTER TABLE comment ADD CONSTRAINT comment_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES "user"(id) ON UPDATE CASCADE;

ALTER TABLE comment ADD COLUMN is_deleted boolean NOT NULL DEFAULT FALSE;

CREATE TRIGGER comment_updated_at
  BEFORE UPDATE
  ON comment
  FOR EACH ROW
  EXECUTE PROCEDURE updated_at();

CREATE TRIGGER comment_audit_1
  BEFORE UPDATE
  ON comment
  FOR EACH ROW
  EXECUTE PROCEDURE audit_1();

CREATE TRIGGER comment_audit_2
  BEFORE UPDATE OF updated_by
  ON comment
  FOR EACH ROW
  EXECUTE PROCEDURE audit_2();

CREATE TRIGGER comment_audit_3
  BEFORE UPDATE
  ON comment
  FOR EACH ROW
  EXECUTE PROCEDURE audit_3();

CREATE TRIGGER comment_audit_body
  AFTER UPDATE OF body
  ON comment
  FOR EACH ROW
  WHEN (new.body IS DISTINCT FROM old.body)
  EXECUTE PROCEDURE audit('body');

CREATE TRIGGER comment_audit_is_deleted
  AFTER UPDATE OF is_deleted
  ON comment
  FOR EACH ROW
  WHEN (new.is_deleted IS DISTINCT FROM old.is_deleted)
  EXECUTE PROCEDURE audit('is_deleted');
