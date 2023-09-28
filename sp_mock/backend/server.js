const express = require('express')
const cors = require('cors')
const Layer8 = require("../../middleware/dist/loadWASM.js")
require('dotenv').config()

// INIT
const port = 8000
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

app.get("/success", (req, res)=>{
    console.log("dude you're close...")
    console.log("req.body: ", req.headers)
    res.send("unfortunately, success has no final form.")
})

app.post("/success", (req, res)=>{
    console.log("dude you're even closer...")
    console.log("req.header: ", req.headers)
    console.log("req.body: ", req.body)
    res.send("keep hustling...")
})

app.listen(port, () => {
    console.log(`\nA mock Service Provider backend is now listening on port ${port}.`)
})