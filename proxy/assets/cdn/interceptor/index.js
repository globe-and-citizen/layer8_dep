import * as interceptor from "./interceptor.json" assert { type: "json" };
await import("./wasm_exec.js");

const decode = (encoded) => {
    var str = window.atob(encoded);
    var bytes = new Uint8Array(str.length);
    for (var i = 0; i < str.length; i++) {
        bytes[i] = str.charCodeAt(i);
    }
    return bytes.buffer;
}

(async () => {
    const go = new Go();
    const importObject = go.importObject;
    return await WebAssembly.instantiate(decode(interceptor.default), importObject).then((result) => {
        go.run(result.instance);
        globalThis.layer8 = layer8;
    });
})();
