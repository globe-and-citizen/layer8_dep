const express = require('express')
const cors = require('cors')
const Layer8 = require("../../middleware/dist/loadWASM.js")
require('dotenv').config()

const jwt = require('jsonwebtoken');
const bcrypt = require('bcrypt');

// INIT
const port = 8000
const app = express()

const users = []; // Store users in memory

const SECRET_KEY = 'my_very_secret_key'

// MIDDLEWARE
app.use(express.json())
app.use(Layer8)
app.use(cors())


app.get("/", (req, res)=>{
    console.log("req.body: ", req.body)
    console.log("res.custom_test_prop: ", res.custom_test_prop)
    res.send("Bro, ur poems coming soon. Relax a little.")
})

app.post("/", (req, res)=>{
    console.log("req.body: ", JSON.parse(req.body))
    //console.log("headers:: ", req.headers)
    res.send("Server has registered a POST.")
})


app.post('/api/register', async (req, res) => {
    console.log("req.body: ", req.body)
    const { password, email } = JSON.parse(req.body);
    console.log(password, email)
    const hashedPassword = await bcrypt.hash(password, 10);
    users.push({ email, password: hashedPassword });
    console.log("users: ", users)
    res.status(200).send('User registered successfully!');
});

app.post('/api/login', async (req, res) => {
    console.log("res.custom_test_prop: ", res.custom_test_prop)
    console.log("req.body: ", req.body)
    console.log("users: ", users)
    const { email, password } = JSON.parse(req.body);
    const user = users.find(u => u.email === email);
    console.log("user: ", user)
    if (user && await bcrypt.compare(password, user.password)) {
        const token = jwt.sign({ email }, SECRET_KEY);
        console.log("token", token)
        res.status(200).json({ token });
    } else {
        res.status(401).send('Invalid credentials!');
    }
})

app.listen(port, () => {
    console.log(`\nA mock Service Provider backend is now listening on port ${port}.`)
})

// const express = require('express')
// const cors = require('cors')
// const Layer8 = require("../../middleware/dist/loadWASM.js")
// require('dotenv').config()

// // INIT
// const port = 8000
// const app = express()
// const userDatabase = [{username:"chester", password:"tester"}]


// // MIDDLEWARE
// app.use(Layer8)
// app.use(cors())
// app.use(express.json())

// app.get("/", (req, res)=>{
//     console.log("req.body: ", req.body)
//     console.log("res.custom_test_prop: ", res.custom_test_prop)
//     res.send("Bro, ur poems coming soon. Relax a little.")
// })

// app.get("/success", (req, res)=>{
//     console.log("req.headers: ", req.headers)
//     res.send("Unfortunately, success has no final form.")
// })

// app.post("/success", (req, res)=>{
//     console.log("Dude you're even closer...")
//     console.log("req.header: ", req.headers)
//     console.log("req.body", req.body)
//     res.send("Well done. Never stop hustling.")
// })

// app.get("/login", (req, res)=>{
//     console.log("Arrival at '/login'", req.query)
//     const username = req.query.username
//     const password = req.query.password

//     const index = userDatabase.findIndex((user)=>{
//         return user.username === username
//     })

//     if (userDatabase[index].password == password) {
//         console.log("User now logged in.")
//         res.send("You are logged in")
//     } else {
//         console.log("Error: Usertried to use incorrect password.")
//         res.send("Denied! Get the fuck out you bum.")
//     }
// })

// app.post("/register", (req, res)=>{
//     console.log("Arrival at '/register'")

//     const {username, password} = req.body

//     const newUser = {
//         username,
//         password
//     }

//     if (userDatabase.findIndex((user)=>{
//         return user.username === username
//     }) != -1){
//         console.log("Error: Username already taken.")
//         res.send("Err: Username already exists.")
//     } else {
//         userDatabase.push(newUser)
//         console.log(userDatabase)
//         res.send("new user registered successfully")
//     }
// })

// app.listen(port, () => {
//     console.log(`\nA mock Service Provider backend is now listening on port ${port}.`)
// })