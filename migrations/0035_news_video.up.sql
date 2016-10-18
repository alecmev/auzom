ALTER TABLE news_item ADD COLUMN video text NOT NULL DEFAULT '';
ALTER TABLE news_item ALTER COLUMN video DROP DEFAULT;

CREATE TRIGGER news_item_audit_video
  AFTER UPDATE OF video
  ON news_item
  FOR EACH ROW
  WHEN (new.video IS DISTINCT FROM old.video)
  EXECUTE PROCEDURE audit('video');
