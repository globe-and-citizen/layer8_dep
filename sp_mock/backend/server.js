/**
 * Sets up an Express server with middleware and a single route handler.
 * Listens on port 3000 and logs a message when the server starts.
 */
const express = require('express')
const cors = require('cors')
const Layer8 = require("../../middleware/dist/loadWASM.js")
require('dotenv').config()

// INIT
const port = 3000
const app = express()

// MIDDLEWARE
app.use(Layer8)
app.use(cors())
app.use(express.json())

app.get("/", (req, res)=>{
    console.log("req.body: ", req.body)
    console.log("res.custom_test_prop: ", res.custom_test_prop)
    res.send("Bro, ur poems coming soon. Relax a little.")
})

app.listen(port, () => {
    console.log(`\nA mock Service Provider backend is now listening on port ${port}.`)
})