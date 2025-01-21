-- Insert test data
INSERT INTO [dbo].[project_data]
    (
    [Id], [project], [filename], [timestamp],
    [category], [cost], [cost_unit]
    )
VALUES
    (
        'test1',
        'Demo Project',
        'test.json',
        GETDATE(),
        'Test Category',
        100.50,
        10
);

-- Query the data
SELECT *
FROM [dbo].[project_data];