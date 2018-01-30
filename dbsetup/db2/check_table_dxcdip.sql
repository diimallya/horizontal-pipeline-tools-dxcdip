SELECT count(TABNAME)
FROM SYSCAT.TABLES
WHERE LOWER(TABSCHEMA) = 'dxcdip'
AND LOWER(TABNAME) in ('agreement', 'element', 'rule', 'parties', 'agreement_has_parties')
;
