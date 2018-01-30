const express = require('express');
const app = express();

const Cloudant = require('./setupCloudant')

const Mysql = require('./setupMysql')

const SLEEPTIME = 30000

function sleep(time, callback) {
  console.log("Sleeping for "+time+"ms")
  var stop = new Date().getTime();
  while(new Date().getTime() < stop + time) {
      ;
  }
  callback();
}

const PORT = process.env.PORT || 3000
app.listen(PORT, function () {
  console.log('Example app listening on port '+ PORT);
  
  sleep(SLEEPTIME, () => {
    Cloudant.setup(() => {
      console.log("Completed cloudant setup")

      if (Mysql.credentials(process.env.VCAP_SERVICES) !== null) {
        Mysql.setup(Mysql.credentials(process.env.VCAP_SERVICES), () => {
          console.log("Completed mysql setup")
          // Done
          process.exit(0)
        })
      } else {
        console.log("No mysql setup necessary")
        process.exit(0)
      }
    });
  })
});