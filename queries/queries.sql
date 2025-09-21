-- QUERY 1: A quantidade de horas que cada professor tem comprometido em aulas.
SELECT
    p.id AS professor_id,
    p.first_name,
    p.last_name,
    SUM(EXTRACT(EPOCH FROM (cs.end_time - cs.start_time))) / 3600.0 AS committed_hours
FROM
    PROFESSOR p
JOIN
    CLASS c ON p.id = c.professor_id
JOIN
    CLASS_SCHEDULE cs ON c.id = cs.class_id
GROUP BY
    p.id, p.first_name, p.last_name
ORDER BY
    committed_hours DESC;

-- QUERY BREAK --

-- QUERY 2: Lista de salas com seus horários ocupados.
SELECT
    r.id AS room_id,
    r.room_number,
    b.name AS building_name,
    cs.day_of_week,
    TO_CHAR(cs.start_time, 'HH24:MI') AS start_time,
    TO_CHAR(cs.end_time, 'HH24:MI') AS end_time,
    s.subject_code,
    s.name AS subject_name
FROM
    ROOM r
JOIN
    BUILDING b ON r.building_id = b.id
LEFT JOIN
    CLASS_SCHEDULE cs ON r.id = cs.room_id
LEFT JOIN
    CLASS c ON cs.class_id = c.id
LEFT JOIN
    SUBJECT s ON c.subject_id = s.id
WHERE
    cs.id IS NOT NULL
ORDER BY
    b.name, r.room_number, cs.day_of_week, cs.start_time;