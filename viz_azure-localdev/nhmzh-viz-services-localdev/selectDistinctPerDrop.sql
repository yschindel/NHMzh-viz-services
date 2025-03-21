SELECT
    project,
    filename,
    timestamp,
    COUNT(distinct id) AS distinct_id_count
FROM [dbo].[cost_data]
GROUP BY project, filename, timestamp
ORDER BY project, filename, timestamp