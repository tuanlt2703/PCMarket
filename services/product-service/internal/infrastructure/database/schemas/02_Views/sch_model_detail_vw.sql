CREATE VIEW model_detail AS
SELECT md.id,
       ct.id AS component_type_id,
       ct.name AS component_type_name,
       mf.id AS manufacturer_id,
       mf.name AS manufacturer_name,
       md.model_name,
       md.description,
       md.release_date,
       md.is_active,
       md.created_at,
       md.updated_at
  FROM models md
  JOIN component_types ct ON ct.id = md.component_type_id
  JOIN manufacturers mf ON mf.id = md.manufacturer_id;
