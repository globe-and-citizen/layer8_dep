const fs = require('fs');
const crypto = require("crypto").webcrypto;
globalThis.crypto = crypto;
require('./wasm_exec.js');

const wasmModule = fs.readFileSync('../../middleware/dist/middleware.wasm');
const go = new Go();
const importObject = go.importObject;
WebAssembly.instantiate(wasmModule, importObject).then(async (results) => {
    const instance = results.instance
    go.run(instance);
}).catch((err)=>{
    console.log("Error running loadWASM script: ", err)
});

module.exports = function Layer8(req, res, next) {
    WASMMiddleware(req, res, next);
};


