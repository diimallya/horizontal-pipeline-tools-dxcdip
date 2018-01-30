--<ScriptOptions statementTerminator = ;/>

CREATE SCHEMA dxcdip ;

CREATE TABLE dxcdip.agreement (
  agreementid VARCHAR(36) NOT NULL,
  name VARCHAR(250) NOT NULL,
  status VARCHAR(45) NOT NULL,
  parentElement VARCHAR(36) NULL DEFAULT NULL,
  PRIMARY KEY (agreementid))
;

CREATE TABLE dxcdip.element (
  elementid VARCHAR(36) NOT NULL,
  name VARCHAR(250) NOT NULL,
  type VARCHAR(45) NULL DEFAULT NULL,
  status VARCHAR(45) NULL DEFAULT NULL,
  agreement_agreementid VARCHAR(36) NOT NULL,
  PRIMARY KEY (elementid),
  CONSTRAINT fk_element_agreement
    FOREIGN KEY (agreement_agreementid)
    REFERENCES dxcdip.agreement (agreementid)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
;

CREATE INDEX fk_element_agreement_idx ON dxcdip.element (agreement_agreementid);

CREATE TABLE dxcdip.rule (
  ruleid VARCHAR(36) NOT NULL,
  ruletext VARCHAR(5000) NOT NULL,
  status VARCHAR(45) NULL DEFAULT NULL,
  element_elementid VARCHAR(36) NULL DEFAULT NULL,
  PRIMARY KEY (ruleid),
  CONSTRAINT fk_rule_element1
    FOREIGN KEY (element_elementid)
    REFERENCES dxcdip.element (elementid)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
;
CREATE INDEX fk_rule_element1_idx ON dxcdip.rule (element_elementid ASC);

CREATE TABLE dxcdip.parties (
  partyid VARCHAR(36) NOT NULL,
  name VARCHAR(256) NOT NULL,
  role VARCHAR(45) NOT NULL,
  idtype VARCHAR(45) NULL DEFAULT NULL,
  ipaddress VARCHAR(45) NULL DEFAULT NULL,
  PRIMARY KEY (partyid))
;

CREATE TABLE dxcdip.agreement_has_parties (
  agreement_agreementid VARCHAR(36) NOT NULL,
  parties_partyid VARCHAR(36) NOT NULL,
  PRIMARY KEY (agreement_agreementid, parties_partyid),
  CONSTRAINT fk_agreement_has_parties_agreement1
    FOREIGN KEY (agreement_agreementid)
    REFERENCES dxcdip.agreement (agreementid)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT fk_agreement_has_parties_parties1
    FOREIGN KEY (parties_partyid)
    REFERENCES dxcdip.parties (partyid)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
;

CREATE INDEX fk_agreement_has_parties_parties1_idx ON dxcdip.agreement_has_parties (parties_partyid);
CREATE INDEX fk_agreement_has_parties_agreement1_idx ON dxcdip.agreement_has_parties (agreement_agreementid);
i
