const chai = require('chai');

const setupCloudant = require('../setupCloudant')

var fakeVCAP = `{
  "cloudantNoSQLDB": [
   {
    "credentials": {
     "host": "a-bluemix.cloudant.com",
     "password": "b",
     "port": 443,
     "url": "https://a-bluemix:b@c-bluemix.cloudant.com",
     "username": "d-bluemix"
    }
   }
  ]
 }`

describe('Cloudant setup', () => {

  describe('Get cloudant credentials', () => {
    it('expects to get from vcap_services', () => {
      chai.expect(setupCloudant.getCloudantURL(fakeVCAP)).to.deep.equal({account: "d-bluemix", password: "b"});
    });
  });

});