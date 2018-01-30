
SELECT count(table_name)
FROM information_schema.tables
WHERE table_schema = 'dxcdip'
AND table_name in ('agreement', 'element', 'rule', 'parties', 'agreement_has_parties')
;
