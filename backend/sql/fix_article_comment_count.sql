-- Rebuild article.comment_count from active article comments.
-- Run once after deploying the comment_count synchronization fix.

UPDATE article AS a
SET comment_count = (
    SELECT COUNT(*)
    FROM comment AS c
    WHERE c.target_type = 1
      AND c.target_id = a.id
      AND c.is_deleted = 0
);
