//
// init.js
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

if (!WebAssembly.instantiateStreaming) { // polyfill
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
	const source = await (await resp).arrayBuffer();
	return await WebAssembly.instantiate(source, importObject);
    };
}

let keyboardHandler;
let mouseHandler;

document.addEventListener('keydown', function(ev) {
    console.log("keydown:", ev);
    if (keyboardHandler) {
        keyboardHandler(ev);
    }
})

function init(keyboard, mouse) {
    keyboardHandler = keyboard;
    mouseHandler = mouse
}
function uninit() {
    keyboardHandler = undefined;
    mouseHandler = undefined;
}

const go = new Go();
let mod, inst;
WebAssembly.instantiateStreaming(fetch("hello.wasm"), go.importObject)
    .then((result) => {
        mod = result.module;
        inst = result.instance;
        document.getElementById("runButton").disabled = false;
    });

async function run() {
    console.clear();
    await go.run(inst);
    uninit()
    // reset instance
    inst = await WebAssembly.instantiate(mod, go.importObject);
}
