const { Console } = require('console');
const fs = require('fs');

const crypto = require("crypto").webcrypto;
globalThis.crypto = crypto;
require('./wasm_exec.js');

const wasmCode = fs.readFileSync("../../middleware/bin/middleware.wasm");
const encoded = Buffer.from(wasmCode, 'binary').toString('base64');

// only necessary if you want to export
function exportToJson(encoded){
    json = "\""+encoded+"\""
    fs.writeFile("b64wasm.json", json, err => {
    if (err) {
      console.error(err);
    }
    // file written successfully
  });
}
// exportToJson(encoded) 

function decode(encoded) {
    var binaryString =  Buffer.from(encoded, 'base64').toString('binary');
    var bytes = new Uint8Array(binaryString.length);
    for (var i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes.buffer;
}

const go = new Go();
const importObject = go.importObject;
WebAssembly.instantiate(decode(encoded), importObject).then(async (results) => {
    const instance = results.instance
    go.run(instance);
    console.log("WASM is Loaded")
}).catch((err)=>{
    console.log("Error running loadWASM script: ", err)
});

module.exports = function Layer8(req, res, next) { 
    WASMMiddleware(req, res, next);
};

// module.exports = function Layer8(req, res, next) { 
//     console.log(TestWASM())
// };



