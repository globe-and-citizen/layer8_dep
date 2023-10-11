# Installation

```bash
npm i layer8-middleware-wasm
```

In `server.js`:

```js
const express = require('express')
const layer8 = require('layer8-middleware')

const app = express()

app.use(express.json()) // use right after express.json()
app.use(layer8)
```
