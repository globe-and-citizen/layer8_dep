const express = require("express");
const cors = require("cors");
const POEMS = require("./poems.json");
const jwt = require("jsonwebtoken");
const bcrypt = require("bcrypt");
const app = express();
const users = []; // Store users in memory
const SECRET_KEY = "my_very_secret_key";
// TODO: For future, use a layer8 npm published package for initialising the client and variables
const popsicle = require("popsicle");
const ClientOAuth2 = require("client-oauth2");
const fs = require("fs");

require("dotenv").config();
const port = process.env.PORT;
const FRONTEND_URL = process.env.FRONTEND_URL;
const LAYER8_URL = process.env.LAYER8_URL;
// const port = 8000;
// const FRONTEND_URL = "http://localhost:5173"
// const LAYER8_URL = "http://localhost:5001"
const LAYER8_CALLBACK_URL = `${FRONTEND_URL}/oauth2/callback`;
const LAYER8_RESOURCE_URL = `${LAYER8_URL}/api/user`;

const layer8Auth = new ClientOAuth2({
  clientId: "notanid",
  clientSecret: "absolutelynotasecret!",
  accessTokenUri: `${LAYER8_URL}/api/oauth`,
  authorizationUri: `${LAYER8_URL}/authorize`,
  redirectUri: LAYER8_CALLBACK_URL,
  scopes: ["read:user"],
});

app.get("/healthcheck", (req, res) => {
  console.log("Enpoint for testing");
  console.log("req.body: ", req.body);
  res.send("Bro, ur poems coming soon. Relax a little.");
});

const Layer8 = require("./dist/loadWASM.js");
app.use(Layer8);
app.use(cors());

app.post("/", (req, res) => {
  console.log("Enpoint for testing");
  console.log("headers:: ", req.headers);
  console.log("req.body: ", req.body);
  res.setHeader("x-header-test", "1234");
  res.send("Server has registered a POST.");
});

let counter = 0;
app.get("/nextpoem", (req, res) => {
  counter++;
  let marker = counter % 3;
  console.log("Served: ", POEMS.data[marker].title);
  res.status(200).json(POEMS.data[marker]);
});

app.post("/api/register", async (req, res) => {
  console.log("req.body: ", req.body);
  const { password, email } = req.body;
  console.log(password, email);
  try {
    hashedPassword = await bcrypt.hash(password, 10);
  } catch (err) {
    console.log("err: ", err);
  }
  users.push({ email, password: hashedPassword });
  console.log("users: ", users);
  res.status(200).send("User registered successfully!");
});

app.post("/api/login", async (req, res) => {
  console.log("res.custom_test_prop: ", res.custom_test_prop);
  console.log("req.body: ", req.body);
  console.log("users: ", users);
  const { email, password } = req.body;
  const user = users.find((u) => u.email === email);
  console.log("user: ", user);
  if (user && (await bcrypt.compare(password, user.password))) {
    const token = jwt.sign({ email }, SECRET_KEY);
    console.log("token", token);
    res.status(200).json({ token });
  } else {
    res.status(401).json({ message: "Invalid credentials!" });
  }
});

// Layer8 Components start here
app.get("/api/login/layer8/auth", async (req, res) => {
  console.log("layer8Auth.code.getUri(): ", layer8Auth.code.getUri());
  res.status(200).json({ authURL: layer8Auth.code.getUri() });
});

app.post("/api/login/layer8/auth", async (req, res) => {
  console.log("Do I even run?");
  const { callback_url } = req.body;
  const user = await layer8Auth.code
    .getToken(callback_url)
    .then(async (user) => {
      return await popsicle
        .request(
          user.sign({
            method: "GET",
            url: LAYER8_RESOURCE_URL,
          })
        )
        .then((res) => {
          console.log("response: ", res);
          return JSON.parse(res.body);
        })
        .catch((err) => {
          console.log("from popsicle: ", err);
        });
    })
    .catch((err) => {
      console.log("err: ", err);
    });

  console.log("user: ", user)
  const isEmailVerified = user.is_email_verified.value;
  let displayName = "";
  let countryName = "";

  if (user.display_name) {
    displayName = user.display_name.value;
  }

  if (user.country_name) {
    countryName = user.country_name.value;
  }

  console.log("Display Name: ", displayName);
  console.log("Country Name: ", countryName);
  const token = jwt.sign(
    { isEmailVerified, displayName, countryName },
    SECRET_KEY
  );
  res.status(200).json({ token });
});

app.post("/api/upload", async (req, res) => {
  const file = req.body.get('file');
  const uint8Array = new Uint8Array(await file.arrayBuffer());
  fs.writeFileSync(`./uploads/${file.name}`, uint8Array);
  res.status(200).json({ message: "File uploaded successfully!" });
});

// Layer8 Components end here

app.listen(port, () => {
  console.log(
    `\nA mock Service Provider backend is now listening on port ${port}.`
  );
});
