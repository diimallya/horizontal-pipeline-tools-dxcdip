const mysql = require('mysql')
const ConnectionConfig = require('mysql/lib/ConnectionConfig')
const fs = require('fs')
const parties = require('./parties.json')

const CHECKTABLE = "check_table_dxcdip.sql"
const CHECKFUNCNAME = "count(table_name)"
const CREATETABLE = "create_table_dxcdip.sql"
const DROPTABLE = "drop_table_dxcdip.sql"

function getMysqlCredentials(vcap_services) {
  if (JSON.parse(vcap_services)["compose-for-mysql"] !== undefined) {
    return JSON.parse(vcap_services)["compose-for-mysql"][0].credentials.uri
  }
  return null
}

function setupMysql(uri, callback){
  // Run setup script
  const connectionConfig = new ConnectionConfig(uri);
  connectionConfig.multipleStatements = true
  const connection = mysql.createConnection(connectionConfig)

  // connection
  connection.connect((err) => {
    if (err) {
      console.log(err)
      callback()
    }

    const checkTable = fs.readFileSync("mysql/"+CHECKTABLE)

    connection.query(checkTable.toString(),(err, results) => {
      if (err) {
        console.log(err)
        callback()
      }
      if (results[0][CHECKFUNCNAME] === 0) {
        // Create table

        const createTable = fs.readFileSync("mysql/"+CREATETABLE)

        connection.query(createTable.toString(),(err, res) => {
          if (err) {
            console.log(err)
          }
          console.log("Created mysql tables")

          // Add parties to table
          addPartiesToSql(connection, () => {
            connection.end();
            callback()
          })
        })

      } else {
        console.log("Check found an existing table", results)
        // Add parties to table
        addPartiesToSql(connection, () => {
          connection.end();
          callback()
        })
      }
    })
  });
}

function addPartiesToSql(connection, callback) {
  let columnsAdded = 0

  // Add parties to table
  Object.keys(parties).forEach((partyName)=>{
    
    connection.query("INSERT INTO dxcdip.parties (partyid, name, role) VALUES (?,?,?)", [parties[partyName].uuid, partyName, parties[partyName].userType], (err) => {
      if (err) {
        if (err.code === "ER_DUP_ENTRY") {
          // Skip, this is fine
        } else {
          console.log(err)
          callback()
        }
      }

      if (columnsAdded === Object.keys(parties).length-1) {
        callback()
      }
      columnsAdded++
    })
  });
}

module.exports = {
  setup : setupMysql,
  credentials: getMysqlCredentials,
}