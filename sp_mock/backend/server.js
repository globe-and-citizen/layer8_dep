const express = require('express')
const cors = require('cors')
const Layer8 = require("../../middleware/dist/loadWASM.js")
require('dotenv').config()

const jwt = require('jsonwebtoken');
const bcrypt = require('bcrypt');

// INIT
const port = 3000
const app = express()

const users = []; // Store users in memory

const SECRET_KEY = 'my_very_secret_key'

// MIDDLEWARE
app.use(Layer8)
app.use(cors())
app.use(express.json())

app.get("/", (req, res)=>{
    console.log("req.body: ", req.body)
    console.log("res.custom_test_prop: ", res.custom_test_prop)
    res.send("Bro, ur poems coming soon. Relax a little.")
})

app.post('/api/register', async (req, res) => {
    console.log("req.body: ", req.body)
    const { email, password } = req.body;
    const hashedPassword = await bcrypt.hash(password, 10);
    users.push({ email, password: hashedPassword });
    res.status(200).send('User registered successfully!');
});

app.post('/api/login', async (req, res) => {
    console.log("req.body: ", req.body)
    const { email, password } = req.body;
    const user = users.find(u => u.email === email);
    if (user && await bcrypt.compare(password, user.password)) {
        const token = jwt.sign({ email }, SECRET_KEY);
        res.status(200).json({ token });
    } else {
        res.status(401).send('Invalid credentials!');
    }
})

app.listen(port, () => {
    console.log(`\nA mock Service Provider backend is now listening on port ${port}.`)
})