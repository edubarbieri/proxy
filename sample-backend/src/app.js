
const express = require('express')
const app = express()

const port = process.env.PORT || 3000;
const host = process.env.HOST || '0.0.0.0';

const message = process.env.MESSAGE || 'Hello World!'

const path = process.env.CONTEXT_PATH || '/'

app.get(path, (req, res) => {
  let x = 0
  for (; x < 100000; x++) {
  }
  res.send(message + "-" + x)
})

app.listen(port, host, () => {
  console.log(`Backend listening at http://${host}:${port}`)
})

