WITH
    counts
    AS
    (
        SELECT
            project,
            filename,
            timestamp,
            id,
            COUNT(*) AS occurrence_count
        FROM [dbo].[cost_data]
        GROUP BY project, filename, timestamp, id
    )
SELECT
    project,
    filename,
    timestamp,
    id,
    occurrence_count
FROM counts
WHERE occurrence_count > 1
ORDER BY project, filename, timestamp, occurrence_count DESC