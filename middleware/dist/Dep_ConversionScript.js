const readFileSync = require('fs').readFileSync
const writeFile = require('fs').writeFileSync;
// C:\Ottawa_DT_Dev\Learning_Computers\layer8\interceptor\bin\interceptor__local.wasm
const wasmCode = readFileSync("./bin/interceptor__local.wasm");
const encoded = Buffer.from(wasmCode, 'binary').toString('base64');

function exportToJson(encoded){
    json = "\""+encoded+"\""
    writeFile("interceptor__local", json, err => {
    if (err) {
      console.error(err);
    }
    // file written successfully
  });
}

exportToJson(encoded)


wasmBin = require("./b64wasm.json")

const crypto = require("crypto").webcrypto;
globalThis.crypto = crypto;
require('./public/wasm/wasm_exec.js');

function decode(encoded) {
    var binaryString =  Buffer.from(encoded, 'base64').toString('binary');
    var bytes = new Uint8Array(binaryString.length);
    for (var i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes.buffer;
}