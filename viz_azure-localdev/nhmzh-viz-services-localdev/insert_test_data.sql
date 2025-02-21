INSERT INTO [dbo].[project_data]
    (project, filename, fileid, timestamp, id, lca, ebkph, ebkph_1, ebkph_2, ebkph_3,
    cost, cost_unit, mat_kbob, gwp_absolute, gwp_relative, penr_absolute, penr_relative,
    ubp_absolute, ubp_relative)
VALUES
    -- Concrete elements
    ('Project_A', 'foundation.ifc', 'Project_A/foundation.ifc', '2024-03-20T10:30:00Z', 'ELEM001', 'A1-A3', 'E2.1', 'E2.1.1', 'E2.1.1.1', 'E2.1.1.1.1',
        25000.00, 250.00, 'KB_CONCRETE_01', 2450.50, 0.85, 12500.75, 0.92, 185000.00, 0.78),
    ('Project_A', 'walls.ifc', 'Project_A/walls.ifc', '2024-03-20T10:35:00Z', 'ELEM002', 'A1-A3', 'E2.2', 'E2.2.1', 'E2.2.1.1', 'E2.2.1.1.1',
        18500.00, 185.00, 'KB_CONCRETE_02', 1850.25, 0.75, 9800.50, 0.88, 145000.00, 0.82),

    -- Steel structures
    ('Project_B', 'steel_beams.ifc', 'Project_B/steel_beams.ifc', '2024-03-21T09:15:00Z', 'ELEM003', 'A1-A3', 'E3.1', 'E3.1.1', 'E3.1.1.1', 'E3.1.1.1.1',
        42000.00, 420.00, 'KB_STEEL_01', 3200.75, 0.95, 18500.25, 0.97, 220000.00, 0.89),
    ('Project_B', 'columns.ifc', 'Project_B/columns.ifc', '2024-03-21T09:20:00Z', 'ELEM004', 'A1-A3', 'E3.2', 'E3.2.1', 'E3.2.1.1', 'E3.2.1.1.1',
        35000.00, 350.00, 'KB_STEEL_02', 2800.50, 0.92, 16500.75, 0.94, 195000.00, 0.86),

    -- Wood elements
    ('Project_C', 'timber_roof.ifc', 'Project_C/timber_roof.ifc', '2024-03-22T14:45:00Z', 'ELEM005', 'A1-A3', 'E4.1', 'E4.1.1', 'E4.1.1.1', 'E4.1.1.1.1',
        15000.00, 150.00, 'KB_WOOD_01', 850.25, 0.45, 5500.50, 0.52, 75000.00, 0.48),
    ('Project_C', 'facade.ifc', 'Project_C/facade.ifc', '2024-03-22T14:50:00Z', 'ELEM006', 'A1-A3', 'E4.2', 'E4.2.1', 'E4.2.1.1', 'E4.2.1.1.1',
        22000.00, 220.00, 'KB_WOOD_02', 1200.75, 0.55, 7800.25, 0.62, 95000.00, 0.58),

    -- Mixed materials
    ('Project_D', 'composite_floor.ifc', 'Project_D/composite_floor.ifc', '2024-03-23T11:20:00Z', 'ELEM007', 'A1-A3', 'E5.1', 'E5.1.1', 'E5.1.1.1', 'E5.1.1.1.1',
        28500.00, 285.00, 'KB_MIXED_01', 2100.50, 0.78, 11500.75, 0.82, 165000.00, 0.75),
    ('Project_D', 'interior_walls.ifc', 'Project_D/interior_walls.ifc', '2024-03-23T11:25:00Z', 'ELEM008', 'A1-A3', 'E5.2', 'E5.2.1', 'E5.2.1.1', 'E5.2.1.1.1',
        19500.00, 195.00, 'KB_MIXED_02', 1650.25, 0.72, 8900.50, 0.76, 125000.00, 0.68);