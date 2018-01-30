const Cloudant = require('cloudant')
const fs = require('fs')

const cloudant = Cloudant(getCloudantCredentials(process.env.VCAP_SERVICES))

const DBNAME = 'documents'

function getCloudantCredentials(vcap_services) {
  const creds = JSON.parse(vcap_services).cloudantNoSQLDB[0].credentials

  return {
    account: creds.username,
    password: creds.password
  }
}

function setupCloudant(callback){
  let itemsProcessed = 0
  cloudant.db.create(DBNAME, ()=>{
    const db = cloudant.db.use(DBNAME)
    const files = fs.readdirSync('cloudant')
    files.forEach((item) => {
      fs.readFile('cloudant/'+item, (err, data) => {
        db.insert(data, item, (err, body, header)=>{
          console.log('err : ' + err + ' header: ' + header + ' body: ' + body)
          itemsProcessed++
          if (itemsProcessed == files.length) {
            callback()
          }
        })
      })
    })
  })
}

module.exports = {
  setup : setupCloudant,
  getCloudantCredentials: getCloudantCredentials,
}