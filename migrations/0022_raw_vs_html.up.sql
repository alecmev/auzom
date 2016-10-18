ALTER TABLE comment RENAME COLUMN body TO body_html;
ALTER TABLE comment RENAME COLUMN body_raw TO body;

ALTER TABLE game RENAME COLUMN summary TO summary_html;
ALTER TABLE game RENAME COLUMN summary_raw TO summary;

DROP TRIGGER game_audit_summary_raw ON game;
DROP TRIGGER game_audit_summary ON game;

CREATE TRIGGER game_audit_summary
  AFTER UPDATE OF summary
  ON game
  FOR EACH ROW
  WHEN (new.summary IS DISTINCT FROM old.summary)
  EXECUTE PROCEDURE audit('summary');

ALTER TABLE news_item RENAME COLUMN preview TO preview_html;
ALTER TABLE news_item RENAME COLUMN preview_raw TO preview;
ALTER TABLE news_item RENAME COLUMN body TO body_html;
ALTER TABLE news_item RENAME COLUMN body_raw TO body;

ALTER TABLE season RENAME COLUMN rules TO rules_html;
ALTER TABLE season RENAME COLUMN rules_raw TO rules;

DROP TRIGGER season_audit_rules_raw ON season;

CREATE TRIGGER season_audit_rules
  AFTER UPDATE OF rules
  ON season
  FOR EACH ROW
  WHEN (new.rules IS DISTINCT FROM old.rules)
  EXECUTE PROCEDURE audit('rules');

ALTER TABLE tournament RENAME COLUMN description TO description_html;
ALTER TABLE tournament RENAME COLUMN description_raw TO description;

DROP TRIGGER tournament_audit_description_raw ON tournament;

CREATE TRIGGER tournament_audit_description
  AFTER UPDATE OF description
  ON tournament
  FOR EACH ROW
  WHEN (new.description IS DISTINCT FROM old.description)
  EXECUTE PROCEDURE audit('description');
